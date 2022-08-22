package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type msgServer struct{ Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for the given keeper
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.registerAccount(ctx, msg.ConnectionId); err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}
