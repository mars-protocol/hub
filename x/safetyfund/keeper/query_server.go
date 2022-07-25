package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

type queryServer struct{ k Keeper }

// NewQueryServerImpl creates an implementation of the `QueryServer` interface for the given keeper
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Balances(goCtx context.Context, req *types.QueryBalancesRequest) (*types.QueryBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	balances := qs.k.GetBalances(ctx)

	return &types.QueryBalancesResponse{Balances: balances}, nil
}
