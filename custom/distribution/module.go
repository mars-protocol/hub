package distribution

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/mars-protocol/hub/custom/distribution/keeper"
)

// AppModule must implement the `module.AppModule` interface
var _ module.AppModule = AppModule{}

// AppModule implements an application module for the custom distribution module
type AppModule struct {
	distr.AppModule

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, accountKeeper distrtypes.AccountKeeper,
	bankKeeper distrtypes.BankKeeper, stakingKeeper distrtypes.StakingKeeper,
) AppModule {
	return AppModule{
		AppModule: distr.NewAppModule(cdc, keeper.Keeper, accountKeeper, bankKeeper, stakingKeeper),
		keeper:    keeper,
	}
}

// BeginBlocker returns the begin blocker for the custom distr module
//
// NOTE: this overwrites the vanilla distr module BeginBlocker with our custom fee distribution logic
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}
