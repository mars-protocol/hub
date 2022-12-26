package types

// DefaultGenesisState returns the shuttle module's default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// Validate validates the given instance of the shuttle module's genesis state.
func (gs GenesisState) Validate() error {
	return nil
}
