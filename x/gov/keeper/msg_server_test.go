package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/gov/keeper"
)

func TestProposalMetadataTypeCheck(t *testing.T) {
	ctx, app, _, _, _ := setupTest(t, []VotingPower{{Staked: 1_000_000, Vesting: 0}})

	testCases := []struct {
		name        string
		metadataStr string
		expPass     bool
	}{
		{
			"a valid proposal metadata",
			`{
				"title": "Mock Proposal",
				"authors": ["Larry Engineer <gm@larry.engineer>"],
				"summary": "Mock proposal for testing purposes",
				"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module.",
				"proposal_forum_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				"vote_option_context": "Vote yes if you like this proposal, Vote no if you don't like it."
			}`,
			true,
		},
		{
			"a valid metadata with missing optional fields",
			`{
				"title": "Mock Proposal",
				"authors": [],
				"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module."
			}`,
			true,
		},
		{
			"an invalid metadata with mandatory fields missing",
			`{
				"title": "Mock Proposal",
				"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module."
			}`,
			false,
		},
		{
			"an invalid proposal with extra unexpected fields",
			`{
				"title": "Mock Proposal",
				"authors": ["Larry Engineer <gm@larry.engineer>"],
				"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module.",
				"foo": "bar"
			}`,
			false,
		},
		{
			"empty proposal metadata string is not accepted",
			"",
			false,
		},
	}

	msgServer := keeper.NewMsgServerImpl(app.GovKeeper)

	for _, tc := range testCases {
		_, err := msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, tc.metadataStr))

		if tc.expPass {
			require.NoError(t, err, "expect success but failed: name = %s", tc.name)
		} else {
			require.Error(t, err, "expect error but succeeded: name = %s", tc.name)
		}
	}
}

func TestVoteMetadataTypeCheck(t *testing.T) {
	ctx, app, _, _, _ := setupTest(t, []VotingPower{{Staked: 1_000_000, Vesting: 0}})

	testCases := []struct {
		name        string
		metadataStr string
		expPass     bool
	}{
		{
			"a valid vote metadata",
			`{"justification":"I like the proposal"}`,
			true,
		},
		{
			"a valid metadata with missing optional fields",
			"{}",
			true,
		},
		{
			"an invalid metadata with extra unexpected fields",
			`{"foo":"bar"}`,
			false,
		},
		{
			"empty metadata string is accepted",
			"",
			true,
		},
	}

	msgServer := keeper.NewMsgServerImpl(app.GovKeeper)

	for _, tc := range testCases {
		_, err := msgServer.Vote(ctx, newMsgVote(tc.metadataStr))

		if tc.expPass {
			require.NoError(t, err, "expect success but failed: name = %s", tc.name)
		} else {
			require.Error(t, err, "expect error but succeeded: name = %s", tc.name)
		}
	}
}

func newMsgSubmitProposal(t *testing.T, metadataStr string) *govv1.MsgSubmitProposal {
	addrs := marsapptesting.MakeRandomAccounts(1)
	proposer := addrs[0]

	proposal, err := govv1.NewMsgSubmitProposal([]sdk.Msg{}, sdk.NewCoins(), proposer.String(), metadataStr)
	require.NoError(t, err)

	return proposal
}

func newMsgVote(metadataStr string) *govv1.MsgVote {
	addrs := marsapptesting.MakeRandomAccounts(1)
	voter := addrs[0]

	return govv1.NewMsgVote(voter, 1, govv1.VoteOption_VOTE_OPTION_YES, metadataStr)
}
