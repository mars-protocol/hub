package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NOTE: The latest version (v0.47.0) of the vanilla gov module already uses
// error codes 2-16, so we start from 17.
// https://github.com/cosmos/cosmos-sdk/blob/main/x/gov/types/errors.go
var (
	ErrFailedToQueryVesting = sdkerrors.Register(govtypes.ModuleName, 17, "failed to query vesting contract")
	ErrInvalidMetadata      = sdkerrors.Register(govtypes.ModuleName, 18, "invalid proposal or vote metadata")
)
