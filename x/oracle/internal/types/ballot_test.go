package types

import (
	"fmt"
	"math"
	"strconv"

	"testing"

	"github.com/stretchr/testify/require"

	core "github.com/terra-project/core/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSqrt(t *testing.T) {
	num := sdk.NewDecWithPrec(144, 4)
	floatNum, err := strconv.ParseFloat(num.String(), 64)
	require.NoError(t, err)

	floatNum = math.Sqrt(floatNum)
	num, err = sdk.NewDecFromStr(fmt.Sprintf("%f", floatNum))
	require.NoError(t, err)

	require.Equal(t, sdk.NewDecWithPrec(12, 2), num)
}

func checkFloatEquality(a sdk.Dec, b float64, precision int) bool {
	base := math.Pow10(precision)

	a2 := a.MulInt64(int64(base)).TruncateInt64()
	b2 := int64(b * base)

	return a2 == b2
}

func TestPBPower(t *testing.T) {

	ctx := sdk.NewContext(nil, abci.Header{}, false, nil)
	_, valAccAddrs, sk := GenerateRandomTestCase()
	pb := PriceBallot{}
	ballotPower := int64(0)

	for i := 0; i < len(sk.Validators()); i++ {
		vote := NewPriceVote(sdk.ZeroDec(), core.MicroSDRDenom, sdk.ValAddress(valAccAddrs[i]))
		pb = append(pb, vote)

		valPower := vote.getPower(ctx, sk)
		require.NotEqual(t, int64(0), valPower)

		ballotPower += valPower
	}

	require.Equal(t, ballotPower, pb.Power(ctx, sk))

	// Mix in a fake validator, the total power should not have changed.
	pubKey := secp256k1.GenPrivKey().PubKey()
	faceValAddr := sdk.ValAddress(pubKey.Address())
	fakeVote := NewPriceVote(sdk.OneDec(), core.MicroSDRDenom, faceValAddr)
	pb = append(pb, fakeVote)
	require.Equal(t, ballotPower, pb.Power(ctx, sk))
}

func TestPBWeightedMedian(t *testing.T) {
	tests := []struct {
		inputs      []float64
		weights     []int64
		isValidator []bool
		median      sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Adding fake validator doesn't change outcome
			[]float64{1.0, 2.0, 10.0, 100000.0, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDecWithPrec(10, 0),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(2, 0),
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	var mockValset []MockValidator
	base := math.Pow10(oracleDecPrecision)
	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			mockVal := NewMockValidator(valAddr, power)

			if tc.isValidator[i] {
				mockValset = append(mockValset, mockVal)
			}
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*base), int64(oracleDecPrecision)), core.MicroSDRDenom, valAddr)
			pb = append(pb, vote)
		}

		sk := NewDummyStakingKeeper(mockValset)

		ctx := sdk.NewContext(nil, abci.Header{}, false, nil)
		require.Equal(t, tc.median, pb.WeightedMedian(ctx, sk))
	}
}

func TestPBStandardDeviation(t *testing.T) {
	tests := []struct {
		inputs            []float64
		weights           []int64
		isValidator       []bool
		standardDeviation sdk.Dec
	}{
		{
			// Supermajority one number
			[]float64{1.0, 2.0, 10.0, 100000.0},
			[]int64{1, 1, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(49995000363, oracleDecPrecision),
		},
		{
			// Adding fake validator doesn't change outcome
			[]float64{1.0, 2.0, 10.0, 100000.0, 10000000000},
			[]int64{1, 1, 100, 1, 10000},
			[]bool{true, true, true, true, false},
			sdk.NewDecWithPrec(4472135950751006, oracleDecPrecision),
		},
		{
			// Tie votes
			[]float64{1.0, 2.0, 3.0, 4.0},
			[]int64{1, 100, 100, 1},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(1224745, oracleDecPrecision),
		},
		{
			// No votes
			[]float64{},
			[]int64{},
			[]bool{true, true, true, true},
			sdk.NewDecWithPrec(0, 0),
		},
	}

	var mockValset []MockValidator
	base := math.Pow10(oracleDecPrecision)
	for _, tc := range tests {
		pb := PriceBallot{}
		for i, input := range tc.inputs {
			valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

			power := tc.weights[i]
			mockVal := NewMockValidator(valAddr, power)

			if tc.isValidator[i] {
				mockValset = append(mockValset, mockVal)
			}
			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*base), int64(oracleDecPrecision)), core.MicroSDRDenom, valAddr)
			pb = append(pb, vote)
		}

		sk := NewDummyStakingKeeper(mockValset)

		ctx := sdk.NewContext(nil, abci.Header{}, false, nil)
		require.Equal(t, tc.standardDeviation, pb.StandardDeviation(ctx, sk))
	}
}

func TestString(t *testing.T) {
	pb := PriceBallot{}
	require.Equal(t, "PriceBallot of 0 votes\n", pb.String())

	vote := NewPriceVote(sdk.NewDecWithPrec(int64(1123400), int64(oracleDecPrecision)), core.MicroSDRDenom, sdk.ValAddress{})
	pb = append(pb, vote)
	require.Equal(t, "PriceBallot of 1 votes\n\n  PriceVote\n\tDenom:    usdr, \n\tVoter:    , \n\tPrice:    1.123400000000000000", pb.String())
}

// func TestPBTally(t *testing.T) {
// 	_, addrs :=mock.GeneratePrivKeyAddressPairs(3)
// 	tests := []struct {
// 		inputs    []float64
// 		weights   []int64
// 		rewardees []sdk.AccAddress
// 	}{
// 		{
// 			// Supermajority one number
// 			[]float64{1.0, 2.0, 10.0, 100000.0},
// 			[]int64{1, 1, 100, 1},
// 			[]sdk.AccAddress{addrs[2]},
// 		},
// 		{
// 			// Tie votes
// 			[]float64{1.0, 2.0, 3.0, 4.0},
// 			[]int64{1, 100, 100, 1},
// 			[]sdk.AccAddress{addrs[1]},
// 		},
// 		{
// 			// No votes
// 			[]float64{},
// 			[]int64{},
// 			[]sdk.AccAddress{},
// 		},

// 		{
// 			// Lots of random votes
// 			[]float64{1.0, 78.48, 78.11, 79.0},
// 			[]int64{1, 51, 79, 33},
// 			[]sdk.AccAddress{addrs[1], addrs[2], addrs[3]},
// 		},
// 	}

// 	for _, tc := range tests {
// 		pb := PriceBallot{}
// 		for i, input := range tc.inputs {
// 			vote := NewPriceVote(sdk.NewDecWithPrec(int64(input*100), 2), "",
// 				sdk.NewInt(tc.weights[i]), addrs[i])
// 			pb = append(pb, vote)
// 		}

// 		_, rewardees := pb.Tally()
// 		require.Equal(t, len(tc.rewardees), len(rewardees))
// 	}
// }
