package keeper

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

const (
	// timeoutTime is the timeout time for packets.
	//
	// Currently we set this as a constant of 5 minutes. Should be make this a
	// configurable parameter? Or a part of the proposal?
	timeoutTime = 15 * time.Minute

	// memo is the memo string to be attached to packets.
	memo = "🛰️ INCOMING TRANSMISSION FROM MARS HUB"
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for
// the given keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

// RegisterAccount creates a new interchain account, or if an interchain account
// already exists but its channel is closed, reopen a new channel.
//
// There are two cases this function may be invoked:
//
//   - For a connection where the shuttle module has never registered an ICA,
//     register one and open a channel for it;
//   - For a connection where the shuttle module already owns an ICA, but the
//     active channel associated with it has been closed, then open a new
//     channel and set it as the new open channel.
//
// Per IBC specs, an ordered channel is closed if a packet times out:
// https://github.com/cosmos/ibc-go/blob/v6.1.0/modules/core/04-channel/keeper/timeout.go#L173-L175
//
// We don't need to check for duplicate ICAs here; the controller module does
// this for us:
// https://github.com/cosmos/ibc-go/blob/v6.1.0/modules/apps/27-interchain-accounts/controller/keeper/account.go#L52-L56
//
// ## IMPORTANT NOTE
//
// In order versions of ibc-go there is a bug with the ICA host module that
// prevents a closed ICA channel from being reopened.
//
// Specifically, the issue is this: When opening the new channel, the host chain
// asserts that the old and new version strings match. On my controller chain,
// the ICA controller module is wrapped in an ICS-29 fee middleware. Therefore
// the channel's version is something like this:
//
// {"fee_version":"ics29-1","app_version":"{\"version\":\"ics27-1\", ...}"}
//
// However during channel handshake, the counterpartyVersion that the ICA host
// module receives is:
//
// {"version":"ics27-1", ...}
//
// That is, without the fee middleware wrapper. Comparing these two will of
// course result in a mismatch error!
//
// This has been fixed in v6:
// https://github.com/cosmos/ibc-go/pull/2302
//
// However, as of ibc-go v4.2.0 (I'm using wasmd 0.30.0 for testing, which comes
// with this version of ibc-go), the fix is not included.
//
// This means that if we are to create an ICA on a host chain with ibc-go <v6,
// we must be suuuuper careful not to have packets timed out... in which case we
// won't be able to reopen the channel!!!
func (ms msgServer) RegisterAccount(goCtx context.Context, req *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// the interchain account is to be owned by the shuttle module account
	owner := ms.k.GetModuleAddress()

	// build and execute the register interchain account message
	//
	// use an empty string as version here. the controller module will generate
	// a default version string for us
	msg := icacontrollertypes.NewMsgRegisterInterchainAccount(req.ConnectionId, owner.String(), "")
	res, err := ms.k.executeMsg(ctx, msg)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: emit the events!
	// the IBC relayer listens to these events
	ctx.EventManager().EmitEvents(res.GetEvents())

	// TODO: currently this gets printed to CLI even during tx simulations.
	// how can we make it only appear during actual deliverTx?
	ms.k.Logger(ctx).Info(
		"initiated interchain account channel handshake",
		"connectionID", req.ConnectionId,
	)

	return &types.MsgRegisterAccountResponse{}, nil
}

func (ms msgServer) SendFunds(goCtx context.Context, req *types.MsgSendFunds) (*types.MsgSendFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", ms.k.authority, req.Authority)
	}

	owner, portID, err := ms.k.GetOwnerAndPortID()
	if err != nil {
		return nil, err
	}

	// query details of the transfer channel
	//
	// the objective is to find the connection id associated with the channel
	channel, found := ms.k.channelKeeper.GetChannel(ctx, ibctransfertypes.PortID, req.ChannelId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "channel with port ID %s and channel ID %s does not exist", ibctransfertypes.PortID, req.ChannelId)
	}

	// the transfer channel must only have one hop
	//
	// we do not need to support multihop channels, as Mars Hub will establish
	// direct connections with all its outpost chains.
	if len(channel.ConnectionHops) > 1 {
		return nil, sdkerrors.Wrapf(types.ErrMultihopUnsupported, "%s has more than one connection hops", req.ChannelId)
	}

	// find the interchain account address associated with the connection
	connectionID := channel.ConnectionHops[0]
	address, found := ms.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no interchain account exists on %s", connectionID)
	}

	// find token balances of the shuttle module account
	balance := ms.k.bankKeeper.GetAllBalances(ctx, owner)

	// if the proposal requires sending more coins than what the module acocunt
	// holds, then draw the difference from the community pool
	shortfall := saturateSub(req.Amount, balance)
	if !shortfall.Empty() {
		if err = ms.k.distrKeeper.DistributeFromFeePool(ctx, shortfall, owner); err != nil {
			return nil, err
		}
	}

	// set timeout parameters
	// we use the timestamp and not the height.
	// note that the timeoutTimestamp in MsgTransfer is in nanoseconds
	timeoutHeight := ibcclienttypes.Height{}
	timeoutTimestamp := uint64(ctx.BlockTime().Add(timeoutTime).UnixNano())

	// send the funds via ICS-20
	//
	// because ICS-20 only supports one coin per packet, we need to dispatch a
	// packet for each coin.
	//
	// the ibctransferkeeper has a sendTransfer method but it's not public.
	// therefore we need to send a MsgTransfer to the baseapp msgRouter.
	for _, coin := range req.Amount {
		msg := ibctransfertypes.NewMsgTransfer(
			ibctransfertypes.PortID,
			req.ChannelId,
			coin,
			owner.String(),
			address,
			timeoutHeight,
			timeoutTimestamp,
			memo,
		)

		res, err := ms.k.executeMsg(ctx, msg)
		if err != nil {
			return nil, err
		}

		// IMPORTANT: emit the events!
		// the IBC relayer listens to these events
		ctx.EventManager().EmitEvents(res.GetEvents())
	}

	ms.k.Logger(ctx).Info(
		"initiated ICS-20 transfer(s) to interchain account",
		"connectionID", connectionID,
		"channelID", req.ChannelId,
		"amount", req.Amount.String(),
	)

	return &types.MsgSendFundsResponse{}, nil
}

func (ms msgServer) SendMessages(goCtx context.Context, req *types.MsgSendMessages) (*types.MsgSendMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", ms.k.authority, req.Authority)
	}

	owner := ms.k.GetModuleAddress()

	protoMsgs, err := convertToProtoMessages(req.Messages)
	if err != nil {
		return nil, err
	}

	data, err := icatypes.SerializeCosmosTx(ms.k.cdc, protoMsgs)
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: memo,
	}
	msg := icacontrollertypes.NewMsgSendTx(
		owner.String(),
		req.ConnectionId,
		uint64(timeoutTime.Nanoseconds()), // NOTE: should be nanoseconds not seconds
		packetData,
	)

	res, err := ms.k.executeMsg(ctx, msg)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: emit the events!
	// the IBC relayer listens to these events
	ctx.EventManager().EmitEvents(res.GetEvents())

	ms.k.Logger(ctx).Info(
		"initiated ICS-20 transfer(s) to interchain account",
		"connectionID", req.ConnectionId,
		"numMsgs", len(req.Messages),
	)

	return &types.MsgSendMessagesResponse{}, nil
}

// saturateSub subtracts a set of coins from another. If the amount goes below
// zero, it's set to zero.
//
// Example:
// {2A, 3B, 4C} - {1A, 5B, 3D} = {1A, 4C}
func saturateSub(coinsA sdk.Coins, coinsB sdk.Coins) sdk.Coins {
	return coinsA.Sub(coinsA.Min(coinsB)...)
}

// getProtoMessages converts []*codectypes.Any to []proto.Message; returns error
// if any of the Any's does not implemenent the proto.Message interface.
//
// In ../types/tx.go we use sdktx.GetMsgs to convert []*Any to []sdk.Msg.
// Why can't we do the same here? Because the function we want to call next
// wants []proto.Message instead of []sdk.Msg.
//
// Despite sdk.Msg is defined by extending proto.Message (and thus inherits all
// the functions required by the proto.Message interface), a function that takes
// a []proto.Message as argument will not accept []sdk.Msg!
//
// Instead, we need to explicitly build a []proto.Message.
func convertToProtoMessages(anys []*codectypes.Any) ([]proto.Message, error) {
	protoMsgs := []proto.Message{}

	for _, any := range anys {
		protoMsg, ok := any.GetCachedValue().(proto.Message)
		if !ok {
			return nil, sdkerrors.Wrapf(types.ErrInvalidProposalMsg, "%s does not implement the proto.Message interface", any.TypeUrl)
		}

		protoMsgs = append(protoMsgs, protoMsg)
	}

	return protoMsgs, nil
}
