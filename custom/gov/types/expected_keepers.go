package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// WasmKeeper defines the expected interface needed to query smart contracts
type WasmKeeper interface {
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
}
