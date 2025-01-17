package treasury

import (
	"github.com/terra-project/core/x/treasury/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k Keeper) {

	// Check epoch last block
	if !core.IsPeriodLastBlock(ctx, core.BlocksPerEpoch) {
		return
	}

	// Update luna issuance after finish all works
	defer k.UpdateIssuance(ctx)

	// Check probation period
	if ctx.BlockHeight() < (core.BlocksPerEpoch * k.WindowProbation(ctx)) {
		return
	}

	k.SettleSeigniorage(ctx)

	// Update tax-rate and reward-weight of next epoch
	taxRate := k.UpdateTaxPolicy(ctx)
	rewardWeight := k.UpdateRewardPolicy(ctx)
	taxCap := k.UpdateTaxCap(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypePolichUpdate,
			sdk.NewAttribute(types.AttributeKeyTax, taxRate.String()),
			sdk.NewAttribute(types.AttributeKeyReward, rewardWeight.String()),
			sdk.NewAttribute(types.AttributeKeyTaxCap, taxCap.String()),
		),
	)

}
