package shuttle

import (
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
	sequence, err := k.ExecuteRemoteContract(ctx, p.ConnectionId, p.Contract, p.Msg, p.Funds)
	if err != nil {
		return nil
	}

	msgBytes, _ := p.Msg.MarshalJSON()

	// TODO: how should we format the keys in the logging info? i'm using snake_case here for now
	logger := k.Logger(ctx)
	logger.Info(
		"submitted execute remote contract message via interchain account",
		"connection_id", p.ConnectionId,
		"sequence", sequence,
		"contract", p.Contract,
		"msg", string(msgBytes),
		"funds", p.Funds.String(),
	)

	return nil
}

func handleMigrateRemoteContractProposal(ctx sdk.Context, k keeper.Keeper, p *types.MigrateRemoteContractProposal) error {
	sequence, err := k.MigrateRemoteContract(ctx, p.ConnectionId, p.Contract, p.CodeId, p.Msg)
	if err != nil {
		return nil
	}

	msgBytes, _ := p.Msg.MarshalJSON()

	// TODO: how should we format the keys in the logging info? i'm using snake_case here for now
	logger := k.Logger(ctx)
	logger.Info(
		"submitted execute remote contract message via interchain account",
		"connection_id", p.ConnectionId,
		"sequence", sequence,
		"contract", p.Contract,
		"code_id", p.CodeId,
		"execute_msg", string(msgBytes),
	)

	return nil
}
