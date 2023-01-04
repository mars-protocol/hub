package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NOTE: latest version (v0.46.0) of the vanilla gov module already registered
// 2-15, so we start from 16.
var (
	// ErrFailedToQueryVesting represents an error where the gov module fails to
	// query the vesting contract.
	ErrFailedToQueryVesting = sdkerrors.Register(govtypes.ModuleName, 16, "failed to query vesting contract")

	// ErrInvalidMetadata represents an error where the metadata (can be one for
	// a proposal or for a vote) doesn't conform to the required schema.
	ErrInvalidMetadata = sdkerrors.Register(govtypes.ModuleName, 17, "invalid proposal or vote metadata")
)
