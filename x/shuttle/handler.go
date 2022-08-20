package shuttle

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/x/shuttle/keeper"
	"github.com/mars-protocol/hub/x/shuttle/types"
)

// NewMsgHandler creates a new handler for messages
func NewMsgHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized shuttle message type: %T", msg)
	}
}

// NewProposalHandler creates a new handler for governance proposals
func NewProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.ExecuteRemoteContractProposal:
			return handleExecuteRemoteContractProposal(ctx, k, c)
		case *types.MigrateRemoteContractProposal:
			return handleMigrateRemoteContractProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized shuttle proposal content type: %T", c)
		}
	}
}

func handleExecuteRemoteContractProposal(ctx sdk.Context, k keeper.Keeper, p *types.ExecuteRemoteContractProposal) error {
	return errors.New("unimplemented")
}

func handleMigrateRemoteContractProposal(ctx sdk.Context, k keeper.Keeper, p *types.MigrateRemoteContractProposal) error {
	return errors.New("unimplemented")
}
