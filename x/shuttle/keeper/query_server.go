package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/mars-protocol/hub/x/shuttle/types"
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

	// the icaControllerKeeper does not provide a method to enumerate interchain
	// accounts of a particular owner or port id. instead, we simply fetch _all_
	// interchain accounts and filter them by the port id.
	allAccounts := qs.k.icaControllerKeeper.GetAllInterchainAccounts(ctx)
	accounts := []*types.AccountInfo{}
	for _, account := range allAccounts {
		if account.PortId == portID {
			account, err := qs.queryAccount(ctx, account.ConnectionId, portID)
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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "shuttle module-owned ICA: connection ID (%s)", connectionID)
	}

	channelID, found := qs.k.icaControllerKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "ICA open active channel: connectionID (%s), portID (%s)", connectionID, portID)
	}

	channel, found := qs.k.channelKeeper.GetChannel(ctx, portID, channelID)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "IBC channel: portID (%s) channelID (%s)")
	}

	connection, err := qs.k.channelKeeper.GetConnection(ctx, connectionID)
	if err != nil {
		return nil, err
	}

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
	}, nil
}