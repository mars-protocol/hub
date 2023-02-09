package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/gov/keeper"
	"github.com/mars-protocol/hub/x/gov/types"
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
				"summary": "Mock proposal for testing purposes"
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
			"extra unexpected fields are accepted",
			`{
				"title": "Mock Proposal",
				"summary": "Mock proposal for testing purposes",
				"foo": "bar"
			}`,
			true,
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
			"extra unexpected fields are accepted",
			`{"foo":"bar"}`,
			true,
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

func TestLegacyProposalMetadata(t *testing.T) {
	ctx, app, _, _, _ := setupTest(t, []VotingPower{{Staked: 1_000_000, Vesting: 0}})

	macc := app.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	addrs := marsapptesting.MakeRandomAccounts(1)
	proposer := addrs[0]

	msgServer := keeper.NewMsgServerImpl(app.GovKeeper)
	legacyMsgServer := keeper.NewLegacyMsgServerImpl(macc.String(), msgServer)

	content := govv1beta1.NewTextProposal(
		"Test community pool spend proposal",
		"This is a mock proposal for testing the conversion of v1beta1 to v1 proposal",
	)

	expectedMetadataStr, err := json.Marshal(&types.ProposalMetadata{
		Title:   content.GetTitle(),
		Summary: content.GetDescription(),
	})
	require.NoError(t, err)

	legacyMsg, err := govv1beta1.NewMsgSubmitProposal(content, sdk.NewCoins(), proposer)
	require.NoError(t, err)

	_, err = legacyMsgServer.SubmitProposal(ctx, legacyMsg)
	require.NoError(t, err)

	proposal, found := app.GovKeeper.GetProposal(ctx, 1)
	require.Equal(t, true, found)
	require.Equal(t, string(expectedMetadataStr), proposal.Metadata)
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
