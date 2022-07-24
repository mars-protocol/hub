package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the incentives module's name
	ModuleName = "incentives"

	// StoreKey is the incentives module's store key
	StoreKey = ModuleName

	// RouterKey is the incentives module's message route
	RouterKey = ModuleName

	// QuerierRoute is the incentives module's querier route
	QuerierRoute = ModuleName
)

// Keys for the incentives module substore
// Items are stored with the following key: values
//
// - 0x00: uint64
// - 0x01<uint64_bytes>: Schedule
var (
	KeyNextScheduleID = []byte{0x00} // key for the the next schedule id
	KeySchedule       = []byte{0x01} // key for the incentives schedules
)

// GetScheduleKey creates the key for the incentives schedule of the given id
func GetScheduleKey(id uint64) []byte {
	return append(KeySchedule, sdk.Uint64ToBigEndian(id)...)
}
