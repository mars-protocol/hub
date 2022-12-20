package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/x/safety/types"
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the `MsgServer` interface for
// the given keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

func (ms msgServer) SafetyFundSpend(goCtx context.Context, req *types.MsgSafetyFundSpend) (*types.MsgSafetyFundSpendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Authority != ms.k.authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", ms.k.authority, req.Authority)
	}

	recipientAddr, err := sdk.AccAddressFromBech32(req.Recipient)
	if err != nil {
		return nil, err
	}

	if err := ms.k.ReleaseFund(ctx, recipientAddr, req.Amount); err != nil {
		return nil, err
	}

	ms.k.Logger(ctx).Info(
		"released coins from safety fund",
		"recipient", req.Recipient,
		"amount", req.Amount.String(),
	)

	return &types.MsgSafetyFundSpendResponse{}, nil
}
