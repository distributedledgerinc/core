package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidation(t *testing.T) {
	genState := DefaultGenesisState()
	require.NoError(t, ValidateGenesis(genState))

	genState.Params.DailyLunaDeltaCap = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	genState.Params.DailyLunaDeltaCap = sdk.OneDec()
	genState.Params.MaxSwapSpread = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))

	genState.Params.MaxSwapSpread = sdk.ZeroDec()
	genState.Params.MinSwapSpread = sdk.NewDec(-1)
	require.Error(t, ValidateGenesis(genState))
}

func TestGenesisEqual(t *testing.T) {
	genState1 := DefaultGenesisState()
	genState2 := DefaultGenesisState()

	require.True(t, genState1.Equal(genState2))
}

func TestGenesisEmpty(t *testing.T) {
	genState := GenesisState{}
	require.True(t, genState.IsEmpty())
}
