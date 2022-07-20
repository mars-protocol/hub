package safetyfund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/x/safetyfund/keeper"
	"github.com/mars-protocol/hub/x/safetyfund/types"
)

func NewHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized safetyfund message type: %T", msg)
	}
}

func NewSafetyFundSpendProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SafetyFundSpendProposal:
			return handleSafetyFundSpendProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized safetyfund proposal content type: %T", c)
		}
	}
}

func handleSafetyFundSpendProposal(ctx sdk.Context, k keeper.Keeper, p *types.SafetyFundSpendProposal) error {
	recipientAddr, err := sdk.AccAddressFromBech32(p.Recipient)
	if err != nil {
		return err
	}

	if err := k.ReleaseFund(ctx, recipientAddr, p.Amount); err != nil {
		return err
	}

	logger := k.Logger(ctx)
	logger.Info("transferred from safety fund to recipient", "recipient", p.Recipient, "amount", p.Amount.String())

	return nil
}
