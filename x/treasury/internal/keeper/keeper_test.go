package keeper

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/treasury/internal/types"
)

func TestRewardWeight(t *testing.T) {
	input := CreateTestInput(t)

	// See that we can get and set reward weights
	blocksPerEpoch := core.BlocksPerEpoch
	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		input.TreasuryKeeper.SetRewardWeight(input.Ctx, sdk.NewDecWithPrec(i, 2))
	}

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		require.Equal(t, sdk.NewDecWithPrec(i, 2), input.TreasuryKeeper.GetRewardWeight(input.Ctx, i))
	}
}

func TestTaxRate(t *testing.T) {
	input := CreateTestInput(t)

	// See that we can get and set tax rate
	blocksPerEpoch := core.BlocksPerEpoch
	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		input.TreasuryKeeper.SetTaxRate(input.Ctx, sdk.NewDecWithPrec(i, 2))
	}

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		require.Equal(t, sdk.NewDecWithPrec(i, 2), input.TreasuryKeeper.GetTaxRate(input.Ctx, i))
	}
}

func TestTaxCap(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroCNYDenom, sdk.NewInt(i))
		require.Equal(t, sdk.NewInt(i), input.TreasuryKeeper.GetTaxCap(input.Ctx, core.MicroCNYDenom))
	}
}

func TestTaxProceeds(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * core.BlocksPerEpoch)

		proceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(i)))
		input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, proceeds)
		input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, proceeds)
		input.TreasuryKeeper.RecordTaxProceeds(input.Ctx, proceeds)
	}

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * core.BlocksPerEpoch)
		proceeds := sdk.NewCoins(sdk.NewCoin(core.MicroSDRDenom, sdk.NewInt(i*3)))

		require.Equal(t, proceeds, input.TreasuryKeeper.PeekTaxProceeds(input.Ctx, i))
	}
}

func TestMicroLunaIssuance(t *testing.T) {
	input := CreateTestInput(t)

	supply := input.SupplyKeeper.GetSupply(input.Ctx)
	supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.ZeroInt())))
	input.SupplyKeeper.SetSupply(input.Ctx, supply)

	// See that we can get and set luna issuance
	blocksPerEpoch := core.BlocksPerEpoch
	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, sdk.NewInt(i))))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)
		input.TreasuryKeeper.UpdateIssuance(input.Ctx)
	}

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * blocksPerEpoch)

		require.Equal(t, sdk.NewInt(i), input.TreasuryKeeper.GetHistoricalIssuance(input.Ctx, i).AmountOf(core.MicroLunaDenom))
	}
}

func TestPeekEpochSeigniorage(t *testing.T) {
	input := CreateTestInput(t)

	for i := int64(0); i < 10; i++ {
		input.Ctx = input.Ctx.WithBlockHeight(i * core.BlocksPerEpoch)
		supply := input.SupplyKeeper.GetSupply(input.Ctx)

		preIssuance := sdk.NewInt(rand.Int63() + 1)
		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, preIssuance)))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)
		input.TreasuryKeeper.UpdateIssuance(input.Ctx)

		nowIssuance := sdk.NewInt(rand.Int63() + 1)
		supply = supply.SetTotal(sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, nowIssuance)))
		input.SupplyKeeper.SetSupply(input.Ctx, supply)

		targetSeigniorage := preIssuance.Sub(nowIssuance)
		if targetSeigniorage.IsNegative() {
			targetSeigniorage = sdk.ZeroInt()
		}

		require.Equal(t, targetSeigniorage, input.TreasuryKeeper.PeekEpochSeigniorage(input.Ctx, i+1))
	}
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	defaultParams := types.DefaultParams()
	input.TreasuryKeeper.SetParams(input.Ctx, defaultParams)

	retrievedParams := input.TreasuryKeeper.GetParams(input.Ctx)
	require.Equal(t, defaultParams, retrievedParams)
}
