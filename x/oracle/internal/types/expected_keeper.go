package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// expected keeper for staking module
type StakingKeeper interface {
	Validator(ctx sdk.Context, address sdk.ValAddress) stakingexported.ValidatorI // get validator by operator address; nil when validator not found
	TotalBondedTokens(sdk.Context) sdk.Int                                        // total bonded tokens within the validator set
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingexported.ValidatorI) (stop bool))
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) // slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	MaxValidators(sdk.Context) uint16                          // MaxValidators returns the maximum amount of bonded validators
}

// expected keeper for distribution module
type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val stakingexported.ValidatorI, tokens sdk.DecCoins)
}

// expected keeper for supply module
type SupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
	SetSupply(ctx sdk.Context, supply supplyexported.SupplyI)
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) sdk.Error
}

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) // Must be called when a validator is bonded
}
