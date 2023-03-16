package v2

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"

	"github.com/mars-protocol/hub/v2/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v2",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			icacontrollertypes.StoreKey,
			icahosttypes.StoreKey,
			// envoy module does not store anything in the chain state, so doesn't
			// need a store upgrade for it
		},
	},
}
