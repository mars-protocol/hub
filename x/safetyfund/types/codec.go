package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&SafetyFundSpendProposal{},
	)
}
