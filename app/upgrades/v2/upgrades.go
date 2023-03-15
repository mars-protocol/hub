package v2

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"

	"github.com/mars-protocol/hub/x/envoy"
	envoytypes "github.com/mars-protocol/hub/x/envoy/types"
)

// CreateUpgradeHandler creates the upgrade handler for the v2 upgrade.
//
// In this upgrade, we add two new modules, ICA and envoy, without making any
// change to the existing modules.
func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("initializing interchain account module")
		initICAModule(ctx, mm, vm)

		ctx.Logger().Info("initializing envoy module")
		initEnvoyModule(ctx, mm, vm)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func initICAModule(ctx sdk.Context, mm *module.Manager, vm module.VersionMap) {
	vm[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

	controllerParams := icacontrollertypes.Params{
		ControllerEnabled: true,
	}

	hostParams := icahosttypes.Params{
		HostEnabled: true,
		AllowMessages: []string{
			"/cosmos.authz.v1beta1.MsgExec",
			"/cosmos.authz.v1beta1.MsgGrant",
			"/cosmos.authz.v1beta1.MsgRevoke",
			"/cosmos.bank.v1beta1.MsgSend",
			"/cosmos.bank.v1beta1.MsgMultiSend",
			"/cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse",
			"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
			"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
			"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
			"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
			"/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
			"/cosmos.gov.v1.MsgSubmitProposal",
			"/cosmos.gov.v1.MsgDeposit",
			"/cosmos.gov.v1.MsgVote",
			"/cosmos.gov.v1.MsgVoteWeighted",
			"/cosmos.slashing.v1beta1.MsgUnjail",
			"/cosmos.staking.v1beta1.MsgCreateValidator",
			"/cosmos.staking.v1beta1.MsgEditValidator",
			"/cosmos.staking.v1beta1.MsgDelegate",
			"/cosmos.staking.v1beta1.MsgBeginRedelegate",
			"/cosmos.staking.v1beta1.MsgUndelegate",
			"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
			"/ibc.applications.transfer.v1.MsgTransfer",
			"/cosmwasm.wasm.v1.MsgStoreCode",
			"/cosmwasm.wasm.v1.MsgInstantiateContract",
			"/cosmwasm.wasm.v1.MsgInstantiateContract2",
			"/cosmwasm.wasm.v1.MsgExecuteContract",
			"/cosmwasm.wasm.v1.MsgMigrateContract",
			"/cosmwasm.wasm.v1.MsgUpdateAdmin",
			"/cosmwasm.wasm.v1.MsgClearAdmin",
			"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfig",
		},
	}

	icaModule := getAppModule[ica.AppModule](mm, icatypes.ModuleName, "ica.AppModule")
	icaModule.InitModule(ctx, controllerParams, hostParams)
}

func initEnvoyModule(ctx sdk.Context, mm *module.Manager, vm module.VersionMap) {
	vm[envoytypes.ModuleName] = mm.Modules[envoytypes.ModuleName].ConsensusVersion()

	envoyModule := getAppModule[envoy.AppModule](mm, envoytypes.ModuleName, "envoy.AppModule")
	envoyModule.InitModule(ctx, &envoytypes.GenesisState{})
}

func getAppModule[T module.AppModule](mm *module.Manager, moduleName, typeName string) T {
	module, correctTypeCast := mm.Modules[moduleName].(T)
	if !correctTypeCast {
		panic(fmt.Sprintf("mm.Modules[\"%s\"] is not of %s", moduleName, typeName))
	}

	return module
}
