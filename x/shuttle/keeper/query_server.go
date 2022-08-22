package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type queryServer struct{ Keeper }

// NewQueryServerImpl creates an implementation of the `QueryServer` interface for the given keeper
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (k queryServer) Account(goCtx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	maccAddr := k.GetModuleAddress()

	portID, err := icatypes.NewControllerPortID(maccAddr.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create port for account: %s", err)
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no ICA channel found for connection %s and port %s", req.ConnectionId)
	}

	addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no interchain account found for connection %s", req.ConnectionId)
	}

	return &types.QueryAccountResponse{ChannelId: channelID, Address: addr}, nil
}

func (k queryServer) Accounts(goCtx context.Context, req *types.QueryAccountsRequest) (*types.QueryAccountsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO")
}
