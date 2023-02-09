package types

import (
	"cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NOTE: The latest version (v0.47.0) of the vanilla gov module already uses
// error codes 2-16, so we start from 17.
// https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/types/errors.go
var (
	ErrFailedToQueryVesting = errors.Register(govtypes.ModuleName, 17, "failed to query vesting contract")
	ErrInvalidMetadata      = errors.Register(govtypes.ModuleName, 18, "invalid proposal or vote metadata")
)
