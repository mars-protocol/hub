package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v3"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// DefaultTimeoutDuration is the default timeout duration
	DefaultTimeoutDuration = 5 * time.Minute

	// KeyTimeoutDuration is the store's key for the TimeoutDuration parameter
	KeyTimeoutDuration = []byte("TimeoutDuration")
)

// DefaultParams is the module's default parameter configuration
func DefaultParams() Params {
	return Params{
		TimeoutDuration: DefaultTimeoutDuration,
	}
}

// ParamKeyTable for the module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTimeoutDuration, &p.TimeoutDuration, validateTimeoutDuration),
	}
}

func validateTimeoutDuration(i interface{}) error {
	td, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid type for parameter TimeoutDuration: %T", i)
	}

	if td == 0 {
		return fmt.Errorf("invalid parameter TimeoutDuration: must be greater than zero")
	}

	return nil
}

// Validate validates the parameters
func (p Params) Validate() error {
	if err := validateTimeoutDuration(p.TimeoutDuration); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
