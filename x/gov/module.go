package gov

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/mars-protocol/hub/v2/x/gov/keeper"
)

// AppModule must implement the `module.AppModule` interface
var _ module.AppModule = AppModule{}

// AppModule implements an application module for the custom gov module
//
// NOTE: our custom AppModule wraps the vanilla `gov.AppModule` to inherit most
// of its functions. However, we overwrite the `EndBlock` function to replace it
// with our custom vote tallying logic.
type AppModule struct {
	gov.AppModule

	keeper        keeper.Keeper
	accountKeeper govtypes.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak govtypes.AccountKeeper, bk govtypes.BankKeeper) AppModule {
	return AppModule{
		AppModule:     gov.NewAppModule(cdc, keeper.Keeper, ak, bk),
		keeper:        keeper,
		accountKeeper: ak,
	}
}

// EndBlock returns the end blocker for the gov module. It returns no validator
// updates.
//
// NOTE: this overwrites the vanilla gov module EndBlocker with our custom vote
// tallying logic.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

// RegisterServices registers module services.
//
// NOTE: this overwrites the vanilla gov module RegisterServices function
func (am AppModule) RegisterServices(cfg module.Configurator) {
	macc := am.accountKeeper.GetModuleAddress(govtypes.ModuleName).String()

	// msg server - use the vanilla implementation
	// The changes we've made to execution are in EndBlocker, so the msgServer
	// doesn't need to be changed.
	msgServer := keeper.NewMsgServerImpl(am.keeper)
	govv1beta1.RegisterMsgServer(cfg.MsgServer(), keeper.NewLegacyMsgServerImpl(macc, msgServer))
	govv1.RegisterMsgServer(cfg.MsgServer(), msgServer)

	// query server - use our custom implementation
	queryServer := keeper.NewQueryServerImpl(am.keeper)
	govv1beta1.RegisterQueryServer(cfg.QueryServer(), keeper.NewLegacyQueryServerImpl(queryServer))
	govv1.RegisterQueryServer(cfg.QueryServer(), queryServer)
}
