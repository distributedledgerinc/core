package types

import (
	"encoding/hex"
	"testing"

	core "github.com/terra-project/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
)

func TestMsgPricePrevote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	bz, err := VoteHash("1", sdk.OneDec(), core.MicroSDRDenom, sdk.ValAddress(addrs[0]))
	require.Nil(t, err)

	tests := []struct {
		hash       string
		denom      string
		voter      sdk.AccAddress
		expectPass bool
	}{
		{hex.EncodeToString(bz), "", addrs[0], false},
		{hex.EncodeToString(bz), core.MicroCNYDenom, addrs[0], true},
		{hex.EncodeToString(bz), core.MicroCNYDenom, addrs[0], true},
		{hex.EncodeToString(bz), core.MicroCNYDenom, sdk.AccAddress{}, false},
		{"", core.MicroCNYDenom, addrs[0], false},
	}

	for i, tc := range tests {
		msg := NewMsgPricePrevote(tc.hash, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgPriceVote(t *testing.T) {
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})

	tests := []struct {
		denom      string
		voter      sdk.AccAddress
		salt       string
		price      sdk.Dec
		expectPass bool
	}{
		{"", addrs[0], "123", sdk.OneDec(), false},
		{core.MicroCNYDenom, addrs[0], "123", sdk.OneDec().MulInt64(core.MicroUnit), true},
		{core.MicroCNYDenom, addrs[0], "123", sdk.ZeroDec(), false},
		{core.MicroCNYDenom, sdk.AccAddress{}, "123", sdk.OneDec().MulInt64(core.MicroUnit), false},
		{core.MicroCNYDenom, addrs[0], "", sdk.OneDec().MulInt64(core.MicroUnit), false},
	}

	for i, tc := range tests {
		msg := NewMsgPriceVote(tc.price, tc.salt, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
