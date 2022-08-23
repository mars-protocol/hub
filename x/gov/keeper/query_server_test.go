package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestQueryServer(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 30, Vesting: 21},
		{Staked: 49, Vesting: 0},
	})

	// validator votes yes
	// voters[0] votes no, overriding the validator's vote
	// this should results in 49 yes vs 51 no
	// in comparison, with the vanilla tallying logic, the result would be 49 yes vs 30 no
	app.GovKeeper.AddVote(ctx, proposal.ProposalId, valoper, govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
	app.GovKeeper.AddVote(ctx, proposal.ProposalId, voters[0], govtypes.NewNonSplitVoteOption(govtypes.OptionNo))

	queryClient := govtypes.NewQueryClient(&baseapp.QueryServiceTestHelper{
		Ctx:             ctx,
		GRPCQueryRouter: app.GRPCQueryRouter(),
	})

	// query proposal - for this one we use the vanilla logic
	{
		res, err := queryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: proposal.ProposalId})
		require.NoError(t, err)
		require.Equal(t, proposal.Content.String(), res.Proposal.Content.String())
	}

	// query votes - for this one we use the vanilla logic
	{
		res, err := queryClient.Vote(context.Background(), &govtypes.QueryVoteRequest{ProposalId: proposal.ProposalId, Voter: voters[0].String()})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Vote.Options))
		require.Equal(t, govtypes.WeightedVoteOption{Option: govtypes.OptionNo, Weight: sdk.NewDec(1)}, res.Vote.Options[0])
	}

	// query tally result - this one is replaced with our custom logic
	{
		res, err := queryClient.TallyResult(context.Background(), &govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId})
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(49), res.Tally.Yes)
		require.Equal(t, sdk.NewInt(51), res.Tally.No)
		require.Equal(t, sdk.NewInt(0), res.Tally.NoWithVeto)
		require.Equal(t, sdk.NewInt(0), res.Tally.Abstain)
	}
}
