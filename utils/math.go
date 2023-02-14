package utils

import sdk "github.com/cosmos/cosmos-sdk/types"

// SaturateSub subtracts a set of coins from another. If the amount goes below
// zero, it's set to zero.
//
// Example:
// {2A, 3B, 4C} - {1A, 5B, 3D} = {1A, 4C}
func SaturateSub(coinsA sdk.Coins, coinsB sdk.Coins) sdk.Coins {
	return coinsA.Sub(coinsA.Min(coinsB)...)
}
