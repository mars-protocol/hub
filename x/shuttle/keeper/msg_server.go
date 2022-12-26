package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for
// the given keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (ms msgServer) RegisterAccount(goCtx context.Context, req *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner := ms.k.GetModuleAddress().String()

	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	// there must not already be an interchain account associated with this
	// connection id
	if _, found := ms.k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID); found {
		return nil, sdkerrors.Wrapf(types.ErrAccountExists, "an interchain account already exists on %s", req.ConnectionId)
	}

	// build the register interchain account message
	// TODO: what to use as the version string??
	version := ""
	msg := icacontrollertypes.NewMsgRegisterInterchainAccount(req.ConnectionId, owner, version)

	// handle the message
	handler := ms.k.router.Handler(msg)
	if _, err = handler(ctx, msg); err != nil {
		return nil, err
	}

	ms.k.Logger(ctx).Info(
		"initiated handshake process for interchain account registration",
		"connectionID", req.ConnectionId,
	)

	return &types.MsgRegisterAccountResponse{}, nil
}
