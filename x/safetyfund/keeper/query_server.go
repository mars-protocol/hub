package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/safetyfund/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Balances(goCtx context.Context, req *types.QueryBalancesRequest) (*types.QueryBalancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	balances := k.GetBalances(ctx)

	return &types.QueryBalancesResponse{Balances: balances}, nil
}
