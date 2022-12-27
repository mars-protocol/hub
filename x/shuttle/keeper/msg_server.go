package keeper

import (
	"context"
	"time"

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
	// timeoutTime is the timeout time for ICS-20 packets.
	//
	// Currently we set this as a constant of 5 minutes. Should be make this a
	// configurable parameter? Or a part of the proposal?
	timeoutTime = 5 * time.Minute

	// memo is the memo string to be attached to ICS-20 packets.
	//
	// TODO: think of something cool to use as memo
	memo = ""
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for
// the given keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (ms msgServer) RegisterAccount(goCtx context.Context, req *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", ms.k.authority, req.Authority)
	}

	owner := ms.k.GetModuleAddress().String()
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	// there must not already be an interchain account associated with this
	// connection id
	if address, found := ms.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID); found {
		return nil, sdkerrors.Wrapf(types.ErrAccountExists, "an interchain account with address %s already exists on %s", address, req.ConnectionId)
	}

	// build the register interchain account message
	//
	// we use an empty string as version here. in this case, the ICA controller
	// middleware will create the default metadata:
	// https://github.com/cosmos/ibc-go/blob/v6.1.0/modules/apps/27-interchain-accounts/controller/keeper/handshake.go#L45-L51
	msg := icacontrollertypes.NewMsgRegisterInterchainAccount(req.ConnectionId, owner, "")

	// handle the message
	handler := ms.k.router.Handler(msg)
	if _, err = handler(ctx, msg); err != nil {
		return nil, err
	}

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

	owner := ms.k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(owner.String())
	if err != nil {
		return nil, err
	}

	// query details of the transfer channel
	//
	// the objective is to find the connection id associated with the channel
	channel, found := ms.k.channelKeeper.GetChannel(ctx, ibctransfertypes.PortID, req.ChannelId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrChannelNotFound, "channel with port ID %s and channel ID %s does not exist", ibctransfertypes.PortID, req.ChannelId)
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
		return nil, sdkerrors.Wrapf(types.ErrAccountNotFound, "no interchain account exists on %s", connectionID)
	}

	// find token balances of the shuttle module account
	balance := ms.k.bankKeeper.GetAllBalances(ctx, owner)

	// if the proposal requires sending more coins than what the module acocunt
	// holds, then draw the difference from the community pool
	shortfall := saturateSub(req.Amount, balance)
	if !shortfall.Empty() {
		ms.k.distrKeeper.DistributeFromFeePool(ctx, shortfall, owner)
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
	handler := ms.k.router.HandlerByTypeURL(sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}))
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
		if _, err := handler(ctx, msg); err != nil {
			return nil, err
		}
	}

	ms.k.Logger(ctx).Info(
		"initiated ICS-20 transfer(s) to interchain account",
		"connectionID", connectionID,
		"channelID", req.ChannelId,
		"amount", req.Amount.String(),
	)

	return &types.MsgSendFundsResponse{}, nil
}

// saturateSub subtracts a set of coins from another. If the amount goes below
// zero, it's set to zero.
//
// Example:
// {2A, 3B, 4C} - {1A, 5B, 3D} = {1A, 4C}
func saturateSub(coinsA sdk.Coins, coinsB sdk.Coins) sdk.Coins {
	return coinsA.Sub(coinsA.Min(coinsB)...)
}
