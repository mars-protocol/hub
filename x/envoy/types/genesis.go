package types

// DefaultGenesisState returns the module's default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// Validate validates the given instance of the module's genesis state.
func (gs GenesisState) Validate() error {
	return nil
}
