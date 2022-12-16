package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govv1beta1.Content)(nil),
		&SafetyFundSpendProposal{},
	)
}
