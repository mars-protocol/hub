package testing

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm"

	marsapp "github.com/mars-protocol/hub/app"
)

// MakeMockApp create a mockup `MarsApp` instance for testing purposes.
// forked from https://github.com/CosmosContracts/juno/blob/v6.0.0/x/mint/keeper/integration_test.go
func MakeMockApp() *marsapp.MarsApp {
	encodingConfig := marsapp.MakeEncodingConfig()

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

	// init chain must be called to stop deliverState from being nil
	genesisState := marsapp.DefaultGenesisState(encodingConfig.Marshaler)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	return app
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
