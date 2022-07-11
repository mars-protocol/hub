package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// ErrFailedToQueryVesting represents an error where the gov module fails to query the vesting contract
// NOTE: the vanilla gov module already registered 2-15, so we start from 16
var ErrFailedToQueryVesting = sdkerrors.Register(govtypes.ModuleName, 16, "failed to query vesting contract")
