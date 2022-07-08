package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// the vanilla gov module already registered 2-9, so we start from 10
var ErrFailedToQueryVesting = sdkerrors.Register(govtypes.ModuleName, 10, "failed to query vesting contract")
