package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	ibcchanneltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibccore "github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/mars-protocol/hub/v2/x/envoy/types"
)

type queryServer struct{ k Keeper }

// NewQuerySErverImpl creates an implementation of the `QueryServer` interface
// for the given keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Account(goCtx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, portID, err := qs.k.GetOwnerAndPortID()
	if err != nil {
		return nil, err
	}

	account, err := qs.queryAccount(ctx, req.ConnectionId, portID)
	if err != nil {
		return nil, err
	}

	return &types.QueryAccountResponse{Account: account}, nil
}

func (qs queryServer) Accounts(goCtx context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, portID, err := qs.k.GetOwnerAndPortID()
	if err != nil {
		return nil, err
	}

	// we only need to fetch the channels with PortID equal to envoy module account's
	// we filter the channels by provided method GetAllChannelsWithPortPrefix()
	// note, this will probably return valid ICA channels for envoy module
	// but we will compare PortID explicitly to keep things safe
	allChannels := qs.k.channelKeeper.GetAllChannelsWithPortPrefix(ctx, portID)
	accounts := []*types.AccountInfo{}

	for _, channel := range allChannels {
		// allChannels is filtered by portID prefix, but not equality
		// the following if-condition may seem unnecessary
		// but, must remain for future-proofing
		if channel.PortId == portID {
			account, err := qs.queryAccountFromChannel(ctx, channel.ChannelId, portID)
			if err != nil {
				return nil, err
			}

			accounts = append(accounts, account)
		}
	}

	return &types.QueryAccountsResponse{Accounts: accounts}, nil
}

func (qs queryServer) queryAccount(ctx sdk.Context, connectionID, portID string) (*types.AccountInfo, error) {
	address, found := qs.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("envoy module-owned ICA: connection ID (%s)", connectionID)
	}

	// ordered channels are closed if a packet times out:
	// https://github.com/cosmos/ibc-go/blob/v6.1.0/modules/core/04-channel/keeper/timeout.go#L173-L175
	//
	// in this case, this method call will fail. simply reopen a new channel by
	// sending a MsgRegisterAccount.
	//
	// it may be better to use a more informative error message here (e.g. "an
	// interchain account exists but the channel is closed.")
	channelID, found := qs.k.icaControllerKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("ICA open active channel: connectionID (%s), portID (%s)", connectionID, portID)
	}

	channel, found := qs.k.channelKeeper.GetChannel(ctx, portID, channelID)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("IBC channel: portID (%s) channelID (%s)", portID, channelID)
	}

	connection, err := qs.k.channelKeeper.GetConnection(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	return composeAccountInfo(address, connectionID, portID, channelID, connection, channel), nil
}

func (qs queryServer) queryAccountFromChannel(ctx sdk.Context, channelID, portID string) (*types.AccountInfo, error) {
	connectionID, connection, err := qs.k.channelKeeper.GetChannelConnection(ctx, portID, channelID)
	if err != nil {
		return nil, err
	}

	channel, found := qs.k.channelKeeper.GetChannel(ctx, portID, channelID)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("IBC channel: portID (%s) channelID (%s)", portID, channelID)
	}

	address, found := qs.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("envoy module-owned ICA: connection ID (%s)", connectionID)
	}

	return composeAccountInfo(address, connectionID, portID, channelID, connection, channel), nil
}

func composeAccountInfo(
	address, connectionID, portID, channelID string,
	connection ibccore.ConnectionI, channel ibcchanneltypes.Channel,
) *types.AccountInfo {
	return &types.AccountInfo{
		Controller: &types.ChainInfo{
			ClientId:     connection.GetClientID(),
			ConnectionId: connectionID,
			PortId:       portID,
			ChannelId:    channelID,
		},
		Host: &types.ChainInfo{
			ClientId:     connection.GetCounterparty().GetClientID(),
			ConnectionId: connection.GetCounterparty().GetConnectionID(),
			PortId:       channel.Counterparty.PortId,
			ChannelId:    channel.Counterparty.ChannelId,
		},
		Address: address,
	}
}
