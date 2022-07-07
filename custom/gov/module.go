package gov

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/mars-protocol/hub/custom/gov/keeper"
)

// AppModule must implement the `module.AppModule` interface
var _ module.AppModule = AppModule{}

// AppModule implements an application module for the custom gov module
//
// NOTE: our custom AppModule wraps the vanilla `gov.AppModule` to inherit most of its functions.
// However, we overwrite the `EndBlock` function to replace it with our custom vote tallying logic
type AppModule struct {
	gov.AppModule

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper) AppModule {
	return AppModule{
		AppModule: gov.NewAppModule(cdc, keeper.Keeper, ak, bk),
		keeper:    keeper,
	}
}

// EndBlock returns the end blocker for the gov module. It returns no validator updates.
//
// NOTE: this overwrites the vanilla gov module EndBlocker with our custom vote tallying logic
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
