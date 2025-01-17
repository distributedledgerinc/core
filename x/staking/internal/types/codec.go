package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(staking.MsgCreateValidator{}, "staking/MsgCreateValidator", nil)
	cdc.RegisterConcrete(staking.MsgEditValidator{}, "staking/MsgEditValidator", nil)
	cdc.RegisterConcrete(staking.MsgDelegate{}, "staking/MsgDelegate", nil)
	cdc.RegisterConcrete(staking.MsgUndelegate{}, "staking/MsgUndelegate", nil)
	cdc.RegisterConcrete(staking.MsgBeginRedelegate{}, "staking/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
