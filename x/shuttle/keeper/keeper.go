package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibcchannelkeeper "github.com/cosmos/ibc-go/v6/modules/core/04-channel/keeper"

	"github.com/mars-protocol/hub/x/shuttle/types"
)

// Keeper is the shuttle module's keeper.
type Keeper struct {
	accountKeeper       authkeeper.AccountKeeper
	bankKeeper          bankkeeper.Keeper
	distrKeeper         distrkeeper.Keeper
	channelKeeper       ibcchannelkeeper.Keeper
	transferKeeper      ibctransferkeeper.Keeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	scopedKeeper        capabilitykeeper.ScopedKeeper

	// The baseapp's message service router.
	// We use this to dispatch messages upon successful governance proposals.
	router *baseapp.MsgServiceRouter

	// The account who can execute shuttle module messages.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new shuttle module keeper.
func NewKeeper(
	accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper,
	distrKeeper distrkeeper.Keeper, channelKeeper ibcchannelkeeper.Keeper,
	transferKeeper ibctransferkeeper.Keeper, icaControllerKeeper icacontrollerkeeper.Keeper,
	scopedKeeper capabilitykeeper.ScopedKeeper, router *baseapp.MsgServiceRouter,
	authority string,
) Keeper {
	// ensure shuttle module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		distrKeeper:         distrKeeper,
		channelKeeper:       channelKeeper,
		transferKeeper:      transferKeeper,
		icaControllerKeeper: icaControllerKeeper,
		scopedKeeper:        scopedKeeper,
		router:              router,
		authority:           authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the shuttle module account's address
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}
