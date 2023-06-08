package testing

import (
	"encoding/json"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/CosmWasm/wasmd/x/wasm"

	marsapp "github.com/mars-protocol/hub/app"
)

// MakeMockApp create a mockup `MarsApp` instance for testing purposes.
//
// Parameters:
//   - accounts: Addresses of base accounts to be registered in the auth module.
//     Do not include module accounts. Do include validator operators.
//   - balances: Coin balances. DO NOT include the community pool and staking
//     pools; their balances will be populated automatically for your
//     convenience. DO include balances of other module accounts if necessary,
//     such as the incentives and safety fund modules.
//   - operators: Addresses of validator operators. Consensus pubkeys will be
//     generated randomly. Each validator does a self-bond of 1_000_000 umars.
//   - communityPool: Initial balance of the community pool.
func MakeMockApp(accounts []sdk.AccAddress, balances []banktypes.Balance, operators []sdk.AccAddress, communityPool sdk.Coins) *marsapp.MarsApp {
	encodingConfig := marsapp.MakeEncodingConfig()

	// create app
	app := marsapp.NewMarsApp(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		marsapp.DefaultNodeHome,
		5,
		encodingConfig,
		simapp.EmptyAppOptions{},
		[]wasm.Option{},
	)

	// create genesis state
	genState := genesisStateWithValSet(encodingConfig.Codec, accounts, balances, operators, communityPool)
	stateBytes, err := json.MarshalIndent(genState, "", " ")
	if err != nil {
		panic(err)
	}

	// run init chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	return app
}

func genesisStateWithValSet(cdc codec.JSONCodec, accounts []sdk.AccAddress, balances []banktypes.Balance, operators []sdk.AccAddress, communityPool sdk.Coins) marsapp.GenesisState {
	// start with the defaut genesis state
	genState := marsapp.DefaultGenesisState(cdc)

	// change default bond denom to umars
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = marsapp.BondDenom

	// each validator self-bonds 1000000 umars
	bondAmt := sdk.DefaultPowerReduction

	// create a random consensus pubkey for each operator
	// create validators and delegations
	pks := MakeRandomConsensusPubkeys(len(operators))
	validators := []stakingtypes.Validator{}
	delegations := []stakingtypes.Delegation{}
	for i, operator := range operators {
		valOper := sdk.ValAddress(operator).String()
		validators = append(validators, stakingtypes.Validator{
			OperatorAddress:   valOper,
			ConsensusPubkey:   pks[i],
			Jailed:            false,
			Status:            stakingtypes.Bonded, // important
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.OneInt(),
		})
		delegations = append(delegations, stakingtypes.Delegation{
			DelegatorAddress: operator.String(),
			ValidatorAddress: valOper,
			Shares:           sdk.OneDec(),
		})
	}

	// genesis accounts
	// all accounts holding unlocked tokens and all validators
	genAccts := []authtypes.GenesisAccount{}
	for _, operator := range accounts {
		genAccts = append(genAccts, authtypes.NewBaseAccountWithAddress(operator))
	}

	// add module account balances:
	// - bonded amount to bonded pool module account
	// - community pool funds to distr module account
	//
	// we do this step after the genesis account step, so that these module
	// accounts don't get dded to genAccts as base accounts
	balances = append(
		balances,
		banktypes.Balance{
			Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(marsapp.BondDenom, bondAmt.Mul(sdk.NewInt(int64(len(operators)))))),
		},
		banktypes.Balance{
			Address: authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
			Coins:   communityPool,
		},
	)

	// total supply of tokens
	supply := sdk.NewCoins()
	for _, balance := range balances {
		supply = supply.Add(balance.Coins...)
	}

	genState[authtypes.ModuleName] = cdc.MustMarshalJSON(authtypes.NewGenesisState(
		authtypes.DefaultParams(),
		genAccts,
	))
	genState[banktypes.ModuleName] = cdc.MustMarshalJSON(banktypes.NewGenesisState(
		banktypes.Params{
			DefaultSendEnabled: true,
		},
		balances,
		supply,
		[]banktypes.Metadata{},
	))
	genState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingtypes.NewGenesisState(
		stakingParams,
		validators,
		delegations,
	))

	distrGenState := distrtypes.DefaultGenesisState()
	distrGenState.FeePool = distrtypes.FeePool{CommunityPool: sdk.NewDecCoinsFromCoins(communityPool...)}
	genState[distrtypes.ModuleName] = cdc.MustMarshalJSON(distrGenState)

	return genState
}

// MakeSimpleMockApp is a shorthand for MakeMockApp that takes no argument.
// This creates a mock app with only one validator and no other accounts.
// Useful for very simple tests.
func MakeSimpleMockApp() *marsapp.MarsApp {
	accts := MakeRandomAccounts(1)
	return MakeMockApp([]sdk.AccAddress{}, []banktypes.Balance{}, []sdk.AccAddress{accts[0]}, sdk.NewCoins())
}

// MakeRandomAccounts returns a list of randomly generated AccAddresses.
// forked from https://github.com/osmosis-labs/osmosis/blob/v9.0.0-rc0/app/apptesting/test_suite.go#L276
func MakeRandomAccounts(numAccts int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, numAccts)
	for i := 0; i < numAccts; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// MakeRandomAccounts returns a list of randonly generated consensus pubkeys.
func MakeRandomConsensusPubkeys(numPks int) []*codectypes.Any {
	pks := simapp.CreateTestPubKeys(numPks)
	anys := []*codectypes.Any{}

	for _, pk := range pks {
		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}

		anys = append(anys, pkAny)
	}

	return anys
}
