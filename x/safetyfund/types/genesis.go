package types

// NewGenesisState returns a new instance of the safetyfund module's genesis state
func NewGenesisState() *GenesisState {
	return &GenesisState{}
}

// DefaultGenesisState returns the default genesis state of the safetyfund module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// ValidateGenesis validates the given instance of the safetyfund module's genesis state
func (GenesisState) Validate() error {
	return nil
}
