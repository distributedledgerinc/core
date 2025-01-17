package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all market state that must be provided at genesis
type GenesisState struct {
	Params       Params  `json:"params" yaml:"params"` // market params
	TaxRate      sdk.Dec `json:"tax_rate" yaml:"tax_rate"`
	RewardWeight sdk.Dec `json:"reward_weight" yaml:"reward_weight"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, taxRate sdk.Dec, rewardWeight sdk.Dec) GenesisState {
	return GenesisState{
		Params:       params,
		TaxRate:      taxRate,
		RewardWeight: rewardWeight,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		TaxRate:      DefaultTaxRate,
		RewardWeight: DefaultRewardWeight,
	}
}

// ValidateGenesis validates the provided oracle genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data GenesisState) error {
	if data.TaxRate.LT(data.Params.TaxPolicy.RateMin) || data.TaxRate.GT(data.Params.TaxPolicy.RateMax) {
		return fmt.Errorf("tax-rate must less than RateMax(%s) and bigger than RateMin(%s)", data.Params.TaxPolicy.RateMax, data.Params.TaxPolicy.RateMin)
	}

	if data.RewardWeight.LT(data.Params.RewardPolicy.RateMin) || data.RewardWeight.GT(data.Params.RewardPolicy.RateMax) {
		return fmt.Errorf("reward-weight must less than WeightMax(%s) and bigger than RateMin(%s)", data.Params.RewardPolicy.RateMax, data.Params.RewardPolicy.RateMin)
	}

	return data.Params.Validate()
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}
