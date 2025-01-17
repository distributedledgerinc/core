package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParamsEqual(t *testing.T) {
	p1 := DefaultParams()
	err := p1.Validate()
	require.NoError(t, err)

	// minus vote period
	p1.VotePeriod = -1
	err = p1.Validate()
	require.Error(t, err)

	// small vote threshold
	p2 := DefaultParams()
	p2.VoteThreshold = sdk.ZeroDec()
	err = p2.Validate()
	require.Error(t, err)

	// negative reward band
	p3 := DefaultParams()
	p3.RewardBand = sdk.NewDecWithPrec(-1, 2)
	err = p3.Validate()
	require.Error(t, err)

	// negative reward fraction
	p4 := DefaultParams()
	p4.RewardFraction = sdk.NewDecWithPrec(-1, 2)
	err = p4.Validate()
	require.Error(t, err)

	// zero slash window
	p5 := DefaultParams()
	p5.VotesWindow = 0
	err = p5.Validate()
	require.Error(t, err)

	// negative slash fraction
	p6 := DefaultParams()
	p6.SlashFraction = sdk.NewDecWithPrec(-1, 2)
	err = p6.Validate()
	require.Error(t, err)
}
