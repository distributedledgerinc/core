package types

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceBallot is a convinience wrapper arounda a PriceVote slice
type PriceBallot []PriceVote

// Returns the total amount of voting power in the ballot
func (pb PriceBallot) Power(ctx sdk.Context, sk StakingKeeper) int64 {
	totalPower := int64(0)
	for _, vote := range pb {
		totalPower += vote.getPower(ctx, sk)
	}

	return totalPower
}

// Returns the median weighted by the power of the PriceVote.
func (pb PriceBallot) WeightedMedian(ctx sdk.Context, sk StakingKeeper) sdk.Dec {
	totalPower := pb.Power(ctx, sk)
	if pb.Len() > 0 {
		if !sort.IsSorted(pb) {
			sort.Sort(pb)
		}

		pivot := int64(0)
		for _, v := range pb {
			votePower := v.getPower(ctx, sk)

			pivot += votePower
			if pivot >= (totalPower / 2) {
				return v.Price
			}
		}
	}
	return sdk.ZeroDec()
}

// Returns the standard deviation by the power of the PriceVote.
func (pb PriceBallot) StandardDeviation(ctx sdk.Context, sk StakingKeeper) (standardDeviation sdk.Dec) {
	if len(pb) == 0 {
		return sdk.ZeroDec()
	}

	median := pb.WeightedMedian(ctx, sk)

	sum := sdk.ZeroDec()
	for _, v := range pb {
		deviation := v.Price.Sub(median)
		sum = sum.Add(deviation.Mul(deviation))
	}

	variance := sum.Quo(sdk.NewDec(int64(len(pb))))

	floatNum, _ := strconv.ParseFloat(variance.String(), 64)
	floatNum = math.Sqrt(floatNum)
	standardDeviation, _ = sdk.NewDecFromStr(fmt.Sprintf("%f", floatNum))

	return
}

// Len implements sort.Interface
func (pb PriceBallot) Len() int {
	return len(pb)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pb PriceBallot) Less(i, j int) bool {
	return pb[i].Price.LTE(pb[j].Price)
}

// Swap implements sort.Interface.
func (pb PriceBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

// String implements fmt.Stringer interface
func (pb PriceBallot) String() (out string) {
	out = fmt.Sprintf("PriceBallot of %d votes\n", pb.Len())
	for _, pv := range pb {
		out += fmt.Sprintf("\n  %s", pv.String())
	}
	return
}
