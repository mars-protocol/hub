package types

// DefaultGenesisState returns the default genesis state of the relay module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// ValidateGenesis validates the given instance of the relay module's genesis state
func (GenesisState) Validate() error {
	return nil
}
