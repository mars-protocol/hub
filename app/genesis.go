package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// The genesis state of the blockchain is represented here as a map of raw JSON
// messages key'd by an identifier string.
//
// The identifier is used to determine which module genesis information belongs
// to, so it may be appropriately routed during init chain.
//
// Within this application, default genesis information is retieved from the
// `ModuleBasicManager` which populates JSON from each `BasicModule` object
// provided to it during init.
type GenesisState map[string]json.RawMessage

// DefaultGenesisState generates the default state for the application.
func DefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}
