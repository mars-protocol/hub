package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/mars-protocol/hub/x/gov/types"
)

// queryVotingPowers queries the vesting contract of user voting powers based on the given query msg
func queryVotingPowers(ctx sdk.Context, k wasmtypes.ViewKeeper, contractAddr sdk.AccAddress, query *types.VotingPowersQuery) (types.VotingPowersResponse, error) {
	var votingPowersResponse types.VotingPowersResponse

	req, err := json.Marshal(&types.QueryMsg{VotingPowers: query})
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrFailedToQueryVesting, "failed to marshal query request: %s", err)
	}

	res, err := k.QuerySmart(ctx, contractAddr, req)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrFailedToQueryVesting, "query returned error: %s", err)
	}

	err = json.Unmarshal(res, &votingPowersResponse)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrFailedToQueryVesting, "failed to unmarshal query response: %s", err)
	}

	return votingPowersResponse, nil
}

// incrementVotingPowers increments the voting power counter based on the contract query response
func incrementVotingPowers(votingPowersResponse types.VotingPowersResponse, tokensInVesting map[string]sdk.Int, totalTokensInVesting *sdk.Int) error {
	for _, item := range votingPowersResponse {
		if _, ok := tokensInVesting[item.User]; ok {
			return sdkerrors.Wrapf(types.ErrFailedToQueryVesting, "query response contains duplicate address: %s", item.User)
		}

		tokensInVesting[item.User] = sdk.Int(item.VotingPower)
		*totalTokensInVesting = totalTokensInVesting.Add(sdk.Int(item.VotingPower))
	}

	return nil
}

// GetTokensInVesting queries the vesting contract for an array of users who have tokens locked in the
// contract and their respective amount, as well as computing the total amount of locked tokens.
func GetTokensInVesting(ctx sdk.Context, k wasmtypes.ViewKeeper, contractAddr sdk.AccAddress) (map[string]sdk.Int, sdk.Int, error) {
	tokensInVesting := make(map[string]sdk.Int)
	totalTokensInVesting := sdk.ZeroInt()

	votingPowersResponse, err := queryVotingPowers(ctx, k, contractAddr, &types.VotingPowersQuery{})
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	if err = incrementVotingPowers(votingPowersResponse, tokensInVesting, &totalTokensInVesting); err != nil {
		return nil, sdk.ZeroInt(), err
	}

	for {
		count := len(votingPowersResponse)
		if count == 0 {
			break
		}

		startAfter := votingPowersResponse[count-1].User

		votingPowersResponse, err = queryVotingPowers(ctx, k, contractAddr, &types.VotingPowersQuery{StartAfter: startAfter})
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		if err = incrementVotingPowers(votingPowersResponse, tokensInVesting, &totalTokensInVesting); err != nil {
			return nil, sdk.ZeroInt(), err
		}
	}

	return tokensInVesting, totalTokensInVesting, nil
}

// MustGetTokensInVesting is the same with `GetTokensInVesting`, but panics on error
func MustGetTokensInVesting(ctx sdk.Context, k wasmtypes.ViewKeeper, contractAddr sdk.AccAddress) (map[string]sdk.Int, sdk.Int) {
	tokensInVesting, totalTokensInVesting, err := GetTokensInVesting(ctx, k, contractAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to tally vote: %s", err))
	}

	return tokensInVesting, totalTokensInVesting
}
