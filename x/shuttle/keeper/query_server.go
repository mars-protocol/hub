package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type queryServer struct{ k Keeper }

// NewQueryServerImpl creates an implementation of the `QueryServer` interface for the given keeper
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	params := qs.k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (qs queryServer) Account(goCtx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	maccAddr := qs.k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(maccAddr.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create port for account: %s", err)
	}

	channelID, found := qs.k.icaControllerKeeper.GetActiveChannelID(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no ICA channel found for connection %s and port %s", req.ConnectionId)
	}

	addr, found := qs.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no interchain account found for connection %s", req.ConnectionId)
	}

	return &types.QueryAccountResponse{ChannelId: channelID, Address: addr}, nil
}

func (qs queryServer) Accounts(goCtx context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	maccAddr := qs.k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(maccAddr.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create port for account: %s", err)
	}

	channels := qs.k.icaControllerKeeper.GetAllActiveChannels(ctx)
	items := []types.QueryAccountsResponseItem{}
	for _, channel := range channels {
		if channel.PortId == portID {
			addr, _ := qs.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, channel.ConnectionId, portID)

			items = append(items, types.QueryAccountsResponseItem{
				ConnectionId: channel.ConnectionId,
				ChannelId:    channel.ChannelId,
				Address:      addr,
			})
		}
	}

	return &types.QueryAccountsResponse{Accounts: items}, nil
}
