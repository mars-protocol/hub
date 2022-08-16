package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/relay/types"
)

type msgServer struct{ Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for the given keeper
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	maccAddr := k.GetModuleAddress()

	portID, err := icatypes.NewControllerPortID(maccAddr.String())
	if err != nil {
		return nil, err
	}

	_, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, msg.ConnectionId, portID)
	if found {
		return nil, sdkerrors.Wrapf(types.ErrAccountExists, "interchain account found for connection %s and port %s", msg.ConnectionId, portID)
	}

	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionId, maccAddr.String()); err != nil {
		return nil, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}
