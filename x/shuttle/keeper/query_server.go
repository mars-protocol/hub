package keeper

import (
	"context"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type queryServer struct{ k Keeper }

// NewQuerySErverImpl creates an implementation of the `QueryServer` interface
// for the given keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

func (qs queryServer) Account(goCtx context.Context, req *types.QueryAccountRequest) (*types.QueryAccountResponse, error) {
	// TODO
	return &types.QueryAccountResponse{Address: "TODO"}, nil
}
