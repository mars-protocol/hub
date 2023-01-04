package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/mars-protocol/hub/x/gov/types"
)

type msgServer struct{ k Keeper }

// NewMsgServerImpl creates an implementation of the gov v1 MsgServer interface
// for the given keeper.
func NewMsgServerImpl(k Keeper) govv1.MsgServer {
	return &msgServer{k}
}

func (ms msgServer) SubmitProposal(goCtx context.Context, msg *govv1.MsgSubmitProposal) (*govv1.MsgSubmitProposalResponse, error) {
	// the metadata string must not be empty. attempt to deserialize it using
	// the given schema return error if fails.
	if _, err := types.UnmarshalProposalMetadata(msg.Metadata); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, err.Error())
	}

	// if metadata is good, we just hand over the rest to the vanilla msgServer
	return govkeeper.NewMsgServerImpl(ms.k.Keeper).SubmitProposal(goCtx, msg)
}

func (ms msgServer) Vote(goCtx context.Context, msg *govv1.MsgVote) (*govv1.MsgVoteResponse, error) {
	// if the metadata string is not empty, attempt to deserialize it using the
	// given schema. return error if fails
	if len(msg.Metadata) > 0 {
		if _, err := types.UnmarshalVoteMetadata(msg.Metadata); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, err.Error())
		}
	}

	// if metadata is good, we just hand over the rest to the vanilla msgServer
	return govkeeper.NewMsgServerImpl(ms.k.Keeper).Vote(goCtx, msg)
}

func (ms msgServer) ExecLegacyContent(goCtx context.Context, msg *govv1.MsgExecLegacyContent) (*govv1.MsgExecLegacyContentResponse, error) {
	return govkeeper.NewMsgServerImpl(ms.k.Keeper).ExecLegacyContent(goCtx, msg)
}

func (ms msgServer) VoteWeighted(goCtx context.Context, msg *govv1.MsgVoteWeighted) (*govv1.MsgVoteWeightedResponse, error) {
	return govkeeper.NewMsgServerImpl(ms.k.Keeper).VoteWeighted(goCtx, msg)
}

func (ms msgServer) Deposit(goCtx context.Context, msg *govv1.MsgDeposit) (*govv1.MsgDepositResponse, error) {
	return govkeeper.NewMsgServerImpl(ms.k.Keeper).Deposit(goCtx, msg)
}
