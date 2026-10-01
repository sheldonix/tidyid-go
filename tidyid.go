package tidyid

import (
	"math"
	"math/big"
)

const (
	// Letters contains the unambiguous lowercase letters used at letter positions.
	Letters = "abcdefghjkmnpqrtuvwxyz"
	// Digits contains the unambiguous digits used at every third position.
	Digits = "23456789"
	// LettersWithUppercase contains the uppercase and lowercase letter alphabet.
	LettersWithUppercase = "ABCDEFGHJKMNPQRTUVWXYZabcdefghjkmnpqrtuvwxyz"

	// DefaultLength is the recommended default ID length.
	DefaultLength = 32
	// MinLength is the smallest supported ID length.
	MinLength = 3
	// MaxLength is the largest supported ID length.
	MaxLength = 256

	// Version is the TidyID Go package version.
	Version = "2.1.0"

	digitAlphabetBits = 3 // len(Digits) is 8, so 8^n is a left shift by 3n.
)

// InvalidIDLengthError reports a length outside 3 through 256.
type InvalidIDLengthError struct{}

func (InvalidIDLengthError) Error() string {
	return "length must be between 3 and 256"
}

// InvalidIDFormatError reports a value that does not follow the TidyID format.
type InvalidIDFormatError struct{}

func (InvalidIDFormatError) Error() string {
	return "value is not a valid TidyID"
}

var (
	lowercaseLetterEntropy = math.Log2(float64(len(Letters)))
	uppercaseLetterEntropy = math.Log2(float64(len(LettersWithUppercase)))
)

// New generates a DefaultLength lowercase TidyID.
func New() (string, error) {
	return Generate(DefaultLength, false)
}

// Generate returns a TidyID of the requested length.
//
// Set allowUppercase to sample letter positions from LettersWithUppercase.
// Every call obtains fresh bytes from the operating system CSPRNG. Generate is
// safe for concurrent use and retains no random state between calls.
func Generate(length int, allowUppercase bool) (string, error) {
	if !isValidLength(length) {
		return "", InvalidIDLengthError{}
	}
	return generateID(length, allowUppercase), nil
}

// IsValidID reports whether value follows the TidyID format.
func IsValidID(value string, allowUppercase bool) bool {
	return isValidID(value, allowUppercase)
}

// IsValidIDOfLength reports whether value follows the TidyID format and has
// exactly the requested valid length.
func IsValidIDOfLength(value string, length int, allowUppercase bool) bool {
	return isValidLength(length) && len(value) == length && hasValidIDFormat(value, allowUppercase)
}

// EnsureValidID validates value and returns InvalidIDFormatError on failure.
func EnsureValidID(value string, allowUppercase bool) error {
	if !isValidID(value, allowUppercase) {
		return InvalidIDFormatError{}
	}
	return nil
}

// EnsureValidIDOfLength validates an exact length and TidyID format. An invalid
// requested length returns InvalidIDLengthError before the value is inspected.
func EnsureValidIDOfLength(value string, length int, allowUppercase bool) error {
	if !isValidLength(length) {
		return InvalidIDLengthError{}
	}
	if len(value) != length || !hasValidIDFormat(value, allowUppercase) {
		return InvalidIDFormatError{}
	}
	return nil
}

// IDCapacity returns the exact number of possible IDs for a length and mode.
// The returned integer is newly allocated and may be modified by the caller.
func IDCapacity(length int, allowUppercase bool) (*big.Int, error) {
	if !isValidLength(length) {
		return nil, InvalidIDLengthError{}
	}

	digitCount := length / 3
	letterCount := length - digitCount
	letterBase := int64(len(Letters))
	if allowUppercase {
		letterBase = int64(len(LettersWithUppercase))
	}
	capacity := new(big.Int).Exp(
		big.NewInt(letterBase),
		big.NewInt(int64(letterCount)),
		nil,
	)
	return capacity.Lsh(capacity, uint(digitAlphabetBits*digitCount)), nil
}

// IDEntropy returns the entropy, in bits, for a length and character mode.
func IDEntropy(length int, allowUppercase bool) (float64, error) {
	if !isValidLength(length) {
		return 0, InvalidIDLengthError{}
	}

	digitCount := length / 3
	letterCount := length - digitCount
	letterEntropy := lowercaseLetterEntropy
	if allowUppercase {
		letterEntropy = uppercaseLetterEntropy
	}
	return float64(letterCount)*letterEntropy +
		float64(digitCount*digitAlphabetBits), nil
}

func isValidLength(length int) bool {
	return length >= MinLength && length <= MaxLength
}
