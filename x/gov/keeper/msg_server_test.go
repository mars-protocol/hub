package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	marsapptesting "github.com/mars-protocol/hub/app/testing"

	"github.com/mars-protocol/hub/x/gov/keeper"
	"github.com/mars-protocol/hub/x/gov/types"
)

func TestProposalMetadataTypeCheck(t *testing.T) {
	ctx, app, _, _, _ := setupTest(t, []VotingPower{{Staked: 1_000_000, Vesting: 0}})

	msgServer := keeper.NewMsgServerImpl(app.GovKeeper)

	// a valid proposal
	_, err := msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, `{
		"title": "Mock Proposal",
		"authors": ["Larry Engineer <gm@larry.engineer>"],
		"summary": "Mock proposal for testing purposes",
		"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module.",
		"proposal_forum_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"vote_option_context": "Vote yes if you like this proposal, Vote no if you don't like it."
	}`))
	require.NoError(t, err)

	// a valid proposal with missing optional fields
	// authors can also be an empty array
	_, err = msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, `{
		"title": "Mock Proposal",
		"authors": [],
		"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module."
	}`))
	require.NoError(t, err)

	// an invalid proposal with mandatory fields missing
	_, err = msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, `{
		"title": "Mock Proposal",
		"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module."
	}`))
	require.Error(t, err, types.ErrInvalidMetadata)

	// extra unexpected fields are not allowed
	_, err = msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, `{
		"title": "Mock Proposal",
		"authors": ["Larry Engineer <gm@larry.engineer>"],
		"details": "This is a mock-up proposal for use in the unit tests of Mars Hub's gov module.",
		"foo": "bar"
	}`))
	require.Error(t, err, types.ErrInvalidMetadata)

	// empty metadata string is not allowed
	_, err = msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, ""))
	require.Error(t, err, types.ErrInvalidMetadata)
}

func TestVoteMetadataTypeCheck(t *testing.T) {
	ctx, app, _, _, _ := setupTest(t, []VotingPower{{Staked: 1_000_000, Vesting: 0}})

	msgServer := keeper.NewMsgServerImpl(app.GovKeeper)

	// a valid vote
	_, err := msgServer.Vote(ctx, newMsgVote(`{"justification":"I like the proposal"}`))
	require.NoError(t, err)

	// a valid proposal with missing optional fields
	_, err = msgServer.Vote(ctx, newMsgVote(`{}`))
	require.NoError(t, err)

	// extra unexpected fields are not allowed
	_, err = msgServer.SubmitProposal(ctx, newMsgSubmitProposal(t, `{"foo":"bar"}`))
	require.Error(t, err, types.ErrInvalidMetadata)

	// empty metadata string is accepted
	_, err = msgServer.Vote(ctx, newMsgVote(""))
	require.NoError(t, err)
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
