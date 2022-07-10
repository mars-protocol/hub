package types

// RewardDenom is the denom of the coin to be distributed as staking rewards. All coins that are NOT
// of this denom goes directly to the community pool.
//
// In a later release, this portion of the fees will be directed to a separate "safety fund".
//
// TODO: make this a configurable parameter
const RewardDenom = "umars"
