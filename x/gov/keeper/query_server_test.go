package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func TestQueryServer(t *testing.T) {
	ctx, app, proposal, valoper, voters := setupTest(t, []VotingPower{
		{Staked: 30_000_000, Vesting: 21_000_000},
		{Staked: 48_000_000, Vesting: 0},
	})

	// validator votes yes
	// voters[0] votes no, overriding the validator's vote
	// this should results in 49 yes vs 51 no
	// in comparison, with the vanilla tallying logic, the result would be 49 yes vs 30 no
	app.GovKeeper.AddVote(ctx, proposal.Id, valoper, govv1.NewNonSplitVoteOption(govv1.OptionYes), "")
	app.GovKeeper.AddVote(ctx, proposal.Id, voters[0], govv1.NewNonSplitVoteOption(govv1.OptionNo), "")

	queryClient := govv1.NewQueryClient(&baseapp.QueryServiceTestHelper{
		Ctx:             ctx,
		GRPCQueryRouter: app.GRPCQueryRouter(),
	})

	// query proposal - for this one we use the vanilla logic
	{
		res, err := queryClient.Proposal(context.Background(), &govv1.QueryProposalRequest{ProposalId: proposal.Id})
		require.NoError(t, err)
		require.Equal(t, proposal.String(), res.Proposal.String())
	}

	// query votes - for this one we use the vanilla logic
	{
		res, err := queryClient.Vote(context.Background(), &govv1.QueryVoteRequest{ProposalId: proposal.Id, Voter: voters[0].String()})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Vote.Options))
		require.Equal(t, &govv1.WeightedVoteOption{Option: govv1.OptionNo, Weight: sdk.NewDec(1).String()}, res.Vote.Options[0])
	}

	// query tally result - this one is replaced with our custom logic
	{
		res, err := queryClient.TallyResult(context.Background(), &govv1.QueryTallyResultRequest{ProposalId: proposal.Id})
		require.NoError(t, err)
		require.Equal(t, "49000000", res.Tally.YesCount)
		require.Equal(t, "51000000", res.Tally.NoCount)
		require.Equal(t, "0", res.Tally.NoWithVetoCount)
		require.Equal(t, "0", res.Tally.AbstainCount)
	}
}
