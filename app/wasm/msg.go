package wasm

import sdk "github.com/cosmos/cosmos-sdk/types"

type MarsMsg struct {
	FundCommunityPool *FundCommunityPool `json:"fund_community_pool,omitempty"`
}

type FundCommunityPool struct {
	Amount sdk.Coins `json:"amount"`
}
