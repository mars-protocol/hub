package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcchannelkeeper "github.com/cosmos/ibc-go/v6/modules/core/04-channel/keeper"

	"github.com/mars-protocol/hub/x/envoy/types"
)

// Keeper is the envoy module's keeper.
type Keeper struct {
	cdc codec.Codec

	accountKeeper       authkeeper.AccountKeeper
	bankKeeper          bankkeeper.Keeper
	distrKeeper         distrkeeper.Keeper
	channelKeeper       ibcchannelkeeper.Keeper
	icaControllerKeeper icacontrollerkeeper.Keeper

	// The baseapp's message service router.
	// We use this to dispatch messages upon successful governance proposals.
	router *baseapp.MsgServiceRouter

	// The account who can execute envoy module messages.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new envoy module keeper.
func NewKeeper(
	cdc codec.Codec, accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper, distrKeeper distrkeeper.Keeper,
	channelKeeper ibcchannelkeeper.Keeper, icaControllerKeeper icacontrollerkeeper.Keeper,
	router *baseapp.MsgServiceRouter, authority string,
) Keeper {
	// ensure envoy module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// make sure the codec is ProtoCodec
	// ICA controller only accepts ProtoCodec for encoding messages:
	// https://github.com/cosmos/ibc-go/blob/v6.1.0/modules/apps/27-interchain-accounts/types/codec.go#L32
	if _, ok := cdc.(*codec.ProtoCodec); !ok {
		panic(fmt.Sprintf("%s module keeper only accepts ProtoCodec; found %T", types.ModuleName, cdc))
	}

	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		distrKeeper:         distrKeeper,
		channelKeeper:       channelKeeper,
		icaControllerKeeper: icaControllerKeeper,
		router:              router,
		authority:           authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetModuleAddress returns the envoy module account's address.
func (k Keeper) GetModuleAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetOwnerAndPortID is a convenience method that returns the envoy module
// account, which acts as the owner of interchain accounts, as well as the ICA
// controller port ID associated with it.
func (k Keeper) GetOwnerAndPortID() (sdk.AccAddress, string, error) {
	owner := k.GetModuleAddress()
	portID, err := icatypes.NewControllerPortID(owner.String())
	return owner, portID, err
}

// executeMsg executes message using the baseapp's message router.
func (k Keeper) executeMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	handler := k.router.Handler(msg)
	return handler(ctx, msg)
}
