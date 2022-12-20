package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Schedule corresponding to the Rust type `mars_vesting::msg::Schedule`
type Schedule struct {
	StartTime uint64 `json:"start_time"`
	Cliff     uint64 `json:"cliff"`
	Duration  uint64 `json:"duration"`
}

// InstantiateMsg corresponding to the Rust type `mars_vesting::msg::InstantiateMsg`
type InstantiateMsg struct {
	Owner          string    `json:"owner"`
	UnlockSchedule *Schedule `json:"unlock_schedule"`
}

// ExecuteMsg corresponding to the Rust enum `mars_vesting::msg::ExecuteMsg`
//
// NOTE: For covenience, we don't include other enum variants, as they are not
// needed here.
type ExecuteMsg struct {
	CreatePosition *CreatePosition `json:"create_position,omitempty"`
	Withdraw       *Withdraw       `json:"withdraw,omitempty"`
}

// CreatePosition corresponding to the Rust enum variant `mars_vesting::msg::ExecuteMsg::CreatePosition`
type CreatePosition struct {
	User         string    `json:"user"`
	VestSchedule *Schedule `json:"vest_schedule"`
}

// Withdraw corresponding to the Rust enum variant `mars_vesting::msg::ExecuteMsg::Withdraw`
type Withdraw struct{}

// QueryMsg corresponding to the Rust enum `mars_vesting::msg::QueryMsg`
//
// NOTE: For covenience, we don't include other enum variants, as they are not
// needed here.
type QueryMsg struct {
	VotingPower  *VotingPowerQuery  `json:"voting_power,omitempty"`
	VotingPowers *VotingPowersQuery `json:"voting_powers,omitempty"`
}

// VotingPowerQuery corresponding to the Rust enum variant `mars_vesting::msg::QueryMsg::VotingPower`
type VotingPowerQuery struct {
	User string `json:"user,omitempty"`
}

// VotingPowersQuery corresponding to the Rust enum variant `mars_vesting::msg::QueryMsg::VotingPowers`
type VotingPowersQuery struct {
	StartAfter string `json:"start_after,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`
}

// VotingPowerResponseItem corresponding to the `voting_powers` query's respons
// type's repeating element
type VotingPowerResponse struct {
	User        string   `json:"user"`
	VotingPower sdk.Uint `json:"voting_power"`
}

// VotingPowerResponse corresponding to the response type of the `voting_powers`
// query
type VotingPowersResponse []VotingPowerResponse
