package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type codeType = sdk.CodeType

// Oracle error codes
const (
	DefaultCodespace sdk.CodespaceType = "oracle"

	CodeUnknownDenom       codeType = 1
	CodeInvalidPrice       codeType = 2
	CodeVoterNotValidator  codeType = 3
	CodeInvalidVote        codeType = 4
	CodeNoVotingPermission codeType = 5
	CodeInvalidHashLength  codeType = 6
	CodeInvalidPrevote     codeType = 7
	CodeVerificationFailed codeType = 8
	CodeNotRevealPeriod    codeType = 9
	CodeInvalidSaltLength  codeType = 10
	CodeInvalidMsgFormat   codeType = 11
	CodeMissingVotingInfo  codeType = 12
)

// ----------------------------------------
// Error constructors

// ErrInvalidHashLength called when the denom is not known
func ErrInvalidHashLength(codespace sdk.CodespaceType, hashLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashLength, fmt.Sprintf("The hash length should equal %d but given %d", tmhash.TruncatedSize, hashLength))
}

// ErrUnknownDenomination called when the denom is not known
func ErrUnknownDenomination(codespace sdk.CodespaceType, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownDenom, fmt.Sprintf("The denom is not known: %s", denom))
}

// ErrInvalidPrice called when the price submitted is not valid
func ErrInvalidPrice(codespace sdk.CodespaceType, price sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrice, fmt.Sprintf("Price is invalid: %s", price.String()))
}

// ErrVoterNotValidator called when the voter is not a validator
func ErrVoterNotValidator(codespace sdk.CodespaceType, voter sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeVoterNotValidator, fmt.Sprintf("Voter is not a validator: %s", voter.String()))
}

// ErrInvalidSignature called when no prevote exists
func ErrVerificationFailed(codespace sdk.CodespaceType, hash []byte, retrivedHash []byte) sdk.Error {
	return sdk.NewError(codespace, CodeVerificationFailed, fmt.Sprintf("Retrieved hash [%s] differs from prevote hash [%s]", retrivedHash, hash))
}

// ErrNoPrevote called when no prevote exists
func ErrNoPrevote(codespace sdk.CodespaceType, voter sdk.ValAddress, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPrevote, fmt.Sprintf("No prevote exists from %s with denom: %s", voter, denom))
}

// ErrNoVote called when no vote exists
func ErrNoVote(codespace sdk.CodespaceType, voter sdk.ValAddress, denom string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("No vote exists from %s with denom: %s", voter, denom))
}

// ErrNoVotingPermission called when the feeder has no permission to submit a vote for the given operator
func ErrNoVotingPermission(codespace sdk.CodespaceType, feeder sdk.AccAddress, operator sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeNoVotingPermission, fmt.Sprintf("Feeder %s not permitted to vote on behalf of: %s", feeder.String(), operator.String()))
}

// ErrNotRevealPeriod called when the feeder submit price reveal vote in wrong period.
func ErrNotRevealPeriod(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotRevealPeriod, fmt.Sprintf("Now is not proper reveal period."))
}

// ErrInvalidSaltLength called when the salt length is not equal 1
func ErrInvalidSaltLength(codespace sdk.CodespaceType, saltLength int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSaltLength, fmt.Sprintf("Salt legnth should be 1~4, but given %d", saltLength))
}

// ErrInvalidMsgFormat called when the msg has invalid format
func ErrInvalidMsgFormat(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMsgFormat, fmt.Sprintf("Invalid Msg Format: %s", msg))
}

// ErrNoVotingInfoFound called when no voting info found
func ErrNoVotingInfoFound(codespace sdk.CodespaceType, valAddr sdk.ValAddress) sdk.Error {
	return sdk.NewError(codespace, CodeMissingVotingInfo, fmt.Sprintf("no signing info found for address: %s", valAddr))
}
