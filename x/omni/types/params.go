package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		// number of block epoch in refreshing balance
		NumEpochs: 5,
		// minimum number of validators voted
		MinConsensus: MIN_CONSENSUS,
		// currnet consensus round
		CurRound: 0,
		// Is whitelisting enabled,
		IsWhitelistEnabled: false,
		// ContractAddress,
		ContractAddress: "0xcc7F90c440ddBd4B082EE7eAA4e7E82E56869C4B",
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateEpochLength(p.NumEpochs); err != nil {
		return err
	}

	if err := validateMinConsensus(p.MinConsensus); err != nil {
		return err
	}

	if err := validateCurRound(p.CurRound); err != nil {
		return err
	}

	if err := validateIsWhitelistEnabled(p.IsWhitelistEnabled); err != nil {
		return err
	}

	if err := validateContractAddress(p.ContractAddress); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateEpochLength(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if epoch interval is 0 block count,
	if v < 1 {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateMinConsensus(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if min consensus < MIN_CONSENSUS
	if v < MIN_CONSENSUS {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateCurRound(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// if cur round is below 0
	if v < 0 {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIsWhitelistEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateContractAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
