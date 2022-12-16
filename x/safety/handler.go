package safety

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/mars-protocol/hub/x/safety/keeper"
	"github.com/mars-protocol/hub/x/safety/types"
)

// NewMsgHandler creates a new handler for messages
func NewMsgHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized safety message type: %T", msg)
	}
}

// NewProposalHandler creates a new handler for governance proposals
func NewProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.SafetyFundSpendProposal:
			return handleSafetyFundSpendProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized safety proposal content type: %T", c)
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
