package market

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/market/internal/types"
)

// NewHandler creates a new handler for all market type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			errMsg := "Unrecognized market Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSwap handles the logic of a MsgSwap
func handleMsgSwap(ctx sdk.Context, k Keeper, ms MsgSwap) sdk.Result {

	// Can't swap to the same coin
	if ms.OfferCoin.Denom == ms.AskDenom {
		return ErrRecursiveSwap(DefaultCodespace, ms.AskDenom).Result()
	}

	// Compute exchange rates between the ask and offer
	swapCoin, spread, swapErr := k.GetSwapCoin(ctx, ms.OfferCoin, ms.AskDenom, false)
	if swapErr != nil {
		return swapErr.Result()
	}

	// Send offer coins to module account
	offerCoins := sdk.NewCoins(ms.OfferCoin)
	err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, ms.Trader, ModuleName, offerCoins)
	if err != nil {
		return err.Result()
	}

	// Charge a spread if applicable; distributed to vote winners in the oracle module
	var swapFee sdk.Coin
	if spread.IsPositive() {
		swapFeeAmt := spread.MulInt(swapCoin.Amount).TruncateInt()
		if swapFeeAmt.IsPositive() {
			swapFee = sdk.NewCoin(swapCoin.Denom, swapFeeAmt)
			swapCoin = swapCoin.Sub(swapFee)
		}
	}

	// Burn offered coins and subtract from the trader's account
	burnErr := k.SupplyKeeper.BurnCoins(ctx, ModuleName, offerCoins)
	if burnErr != nil {
		return burnErr.Result()
	}

	// Mint asked coins and credit Trader's account
	swapCoins := sdk.NewCoins(swapCoin)
	mintErr := k.SupplyKeeper.MintCoins(ctx, ModuleName, swapCoins)
	if mintErr != nil {
		return mintErr.Result()
	}

	sendErr := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, ModuleName, ms.Trader, swapCoins)
	if sendErr != nil {
		return sendErr.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventSwap,
			sdk.NewAttribute(types.AttributeKeyOffer, ms.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyTrader, ms.Trader.String()),
			sdk.NewAttribute(types.AttributeKeySwapCoin, swapCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapFee, swapFee.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
