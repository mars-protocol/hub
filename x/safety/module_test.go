package safety_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	marsapptesting "github.com/mars-protocol/hub/v2/app/testing"
	"github.com/mars-protocol/hub/v2/x/safety/types"
)

// TestCreatesModuleAccountAtGenesis asserts that the safety module account is
// properly registered with the auth module at genesis.
func TestCreatesModuleAccountAtGenesis(t *testing.T) {
	app := marsapptesting.MakeSimpleMockApp()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}
