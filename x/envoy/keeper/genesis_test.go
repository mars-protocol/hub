package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	marsapp "github.com/mars-protocol/hub/v2/app"
	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"
	"github.com/mars-protocol/hub/v2/x/envoy/types"
)

var mockGenesisState = &types.GenesisState{}

func setupGenesisTest() (ctx sdk.Context, app *marsapp.MarsApp) {
	app = marsapptesting.MakeSimpleMockApp()
	ctx = app.BaseApp.NewContext(false, tmproto.Header{})

	app.EnvoyKeeper.InitGenesis(ctx, mockGenesisState)

	return ctx, app
}

func TestInitGenesis(t *testing.T) {
	ctx, app := setupGenesisTest()

	// make sure that the module account is registered at the auth module
	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}

func TestExportGenesis(t *testing.T) {
	ctx, app := setupGenesisTest()

	exported := app.EnvoyKeeper.ExportGenesis(ctx)
	require.Equal(t, exported, mockGenesisState)
}
