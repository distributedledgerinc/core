package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/oracle/internal/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace

	distrKeeper   types.DistributionKeeper
	StakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper

	distrName string

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	paramspace params.Subspace, distrKeeper types.DistributionKeeper,
	stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper,
	distrName string, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramspace.WithKeyTable(ParamKeyTable()),
		distrKeeper:   distrKeeper,
		StakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
		distrName:     distrName,
		codespace:     codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns a codespace of keeper
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

//-----------------------------------
// Prevote logic

// IteratePrevotes iterates rate over prevotes in the store
func (k Keeper) IteratePrevotes(ctx sdk.Context, handler func(prevote types.PricePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PrevoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.PricePrevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

// iteratePrevotesWithPrefix iterates over prevotes in the store with given prefix
func (k Keeper) iteratePrevotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.PricePrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var prevote types.PricePrevote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &prevote)
		if handler(prevote) {
			break
		}
	}
}

//-----------------------------------
// Votes logic

// CollectVotes collects all oracle votes for the period, categorized by the votes' denom parameter
func (k Keeper) CollectVotes(ctx sdk.Context) (votes map[string]types.PriceBallot) {
	votes = map[string]types.PriceBallot{}
	handler := func(vote types.PriceVote) (stop bool) {
		votes[vote.Denom] = append(votes[vote.Denom], vote)
		return false
	}
	k.IterateVotes(ctx, handler)

	return
}

// IterateVotes iterates over votes in the store
func (k Keeper) IterateVotes(ctx sdk.Context, handler func(vote types.PriceVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.VoteKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// Iterate over votes in the store
func (k Keeper) iterateVotesWithPrefix(ctx sdk.Context, prefix []byte, handler func(vote types.PriceVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.PriceVote
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &vote)
		if handler(vote) {
			break
		}
	}
}

// GetPrevote retrieves a prevote from the store
func (k Keeper) GetPrevote(ctx sdk.Context, denom string, voter sdk.ValAddress) (prevote types.PricePrevote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetPrevoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoPrevote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &prevote)
	return
}

// AddPrevote adds a prevote to the store
func (k Keeper) AddPrevote(ctx sdk.Context, prevote types.PricePrevote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(prevote)
	store.Set(types.GetPrevoteKey(prevote.Denom, prevote.Voter), bz)
}

// DeletePrevote deletes a prevote from the store
func (k Keeper) DeletePrevote(ctx sdk.Context, prevote types.PricePrevote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPrevoteKey(prevote.Denom, prevote.Voter))
}

// Retrieves a vote from the store
func (k Keeper) getVote(ctx sdk.Context, denom string, voter sdk.ValAddress) (vote types.PriceVote, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetVoteKey(denom, voter))
	if b == nil {
		err = types.ErrNoVote(k.codespace, voter, denom)
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vote)
	return
}

// AddVote adds a vote to the store
func (k Keeper) AddVote(ctx sdk.Context, vote types.PriceVote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.GetVoteKey(vote.Denom, vote.Voter), bz)
}

// DeleteVote deletes a vote from the store
func (k Keeper) DeleteVote(ctx sdk.Context, vote types.PriceVote) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVoteKey(vote.Denom, vote.Voter))
}

//-----------------------------------
// Price logic

// GetLunaPrice gets the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) GetLunaPrice(ctx sdk.Context, denom string) (price sdk.Dec, err sdk.Error) {
	if denom == core.MicroLunaDenom {
		return sdk.OneDec(), nil
	}

	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetPriceKey(denom))
	if b == nil {
		return sdk.ZeroDec(), types.ErrUnknownDenomination(k.codespace, denom)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &price)
	return
}

// SetLunaPrice sets the consensus exchange rate of Luna denominated in the denom asset to the store.
func (k Keeper) SetLunaPrice(ctx sdk.Context, denom string, price sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(price)
	store.Set(types.GetPriceKey(denom), bz)
}

// IterateLunaPrices iterates over luna prices in the store
func (k Keeper) IterateLunaPrices(ctx sdk.Context, handler func(denom string, price sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PriceKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.PriceKey):])
		var price sdk.Dec
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &price)
		if handler(denom, price) {
			break
		}
	}
}

// DeletePrice deletes the consensus exchange rate of Luna denominated in the denom asset from the store.
func (k Keeper) DeletePrice(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPriceKey(denom))
}

// GetActiveDenoms returns all active oracle asset denoms from the store
func (k Keeper) GetActiveDenoms(ctx sdk.Context) (denoms types.DenomList) {
	denoms = types.DenomList{}

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PriceKey)
	for ; iter.Valid(); iter.Next() {
		n := len(types.PriceKey)
		denom := string(iter.Key()[n:])
		denoms = append(denoms, denom)
	}
	iter.Close()

	return
}

//-----------------------------------
// Feeder delegation logic

// GetFeedDelegate gets the account address that the feeder right was delegated to by the validator operator.
func (k Keeper) GetFeedDelegate(ctx sdk.Context, operator sdk.ValAddress) (delegate sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetFeederDelegationKey(operator))
	if b == nil {
		// By default the right is delegated to the validator itself
		return sdk.AccAddress(operator)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &delegate)
	return
}

// SetFeedDelegate sets the account address that the feeder right was delegated to by the validator operator.
func (k Keeper) SetFeedDelegate(ctx sdk.Context, operator sdk.ValAddress, delegatedFeeder sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(delegatedFeeder)
	store.Set(types.GetFeederDelegationKey(operator), bz)
}

// IterateFeederDelegations iterates over the feeder delegations
// and performs a callback function
func (k Keeper) IterateFeederDelegations(ctx sdk.Context,
	handler func(delegator sdk.ValAddress, delegatee sdk.AccAddress) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.FeederDelegationKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		delegator := sdk.ValAddress(iter.Key()[len(types.FeederDelegationKey):])

		var delegatee sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &delegatee)

		if handler(delegator, delegatee) {
			break
		}
	}
}

//-----------------------------------
// Reward pool logic

// getRewardPool retrieves the balance of the oracle module account
func (k Keeper) getRewardPool(ctx sdk.Context) sdk.Coins {
	acc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	return acc.GetCoins()
}

//-----------------------------------
// VotingInfo logic

// SetVotingInfo sets voting info for a validator
func (k Keeper) SetVotingInfo(ctx sdk.Context, address sdk.ValAddress, votingInfo types.VotingInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(votingInfo)
	store.Set(types.GetVotingInfoKey(address), bz)
	return
}

// getVotingInfo gets voting info for a validator
func (k Keeper) getVotingInfo(ctx sdk.Context, address sdk.ValAddress) (votingInfo types.VotingInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetVotingInfoKey(address))

	if bz == nil {
		found = false
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &votingInfo)
	found = true
	return
}

// IterateVotingInfos iterates over the stored VotingInfo
func (k Keeper) IterateVotingInfos(ctx sdk.Context,
	handler func(info types.VotingInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.VotingInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {

		var info types.VotingInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &info)

		if handler(info) {
			break
		}
	}
}

//-----------------------------------
// MissedVoteBitArray logic

// GetMissedVoteBitArray gets the bit for the missed votes array
// only wrong(price) votes will be marked as missed
func (k Keeper) GetMissedVoteBitArray(ctx sdk.Context, address sdk.ValAddress, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMissedVoteBitArrayKey(address, index))
	if bz == nil {
		// lazy: treat empty key as not missed
		missed = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

// SetMissedVoteBitArray sets the bit that checks if the validator has
// missed a block in the current window
func (k Keeper) SetMissedVoteBitArray(ctx sdk.Context, address sdk.ValAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(types.GetMissedVoteBitArrayKey(address, index), bz)
}

// clearMissedVoteBitArray deletes every instance of MissedVoteBitArray in the store
func (k Keeper) clearMissedVoteBitArray(ctx sdk.Context, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetMissedVoteBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// IterateMissedVoteBitArray iterates over the signed votes window
// and performs a callback function
func (k Keeper) IterateMissedVoteBitArray(ctx sdk.Context,
	address sdk.ValAddress, handler func(index int64, missed bool) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	// Array may be sparse
	for ; index < k.VotesWindow(ctx); index++ {
		var missed bool
		bz := store.Get(types.GetMissedVoteBitArrayKey(address, index))
		if bz == nil {
			continue
		}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
		if handler(index, missed) {
			break
		}
	}
}
