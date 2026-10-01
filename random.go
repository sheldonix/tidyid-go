package tidyid

import "crypto/rand"

var (
	letterLookup    = makeAlphabetLookup(Letters)
	uppercaseLookup = makeAlphabetLookup(LettersWithUppercase)
	digitLookup     = makeAlphabetLookup(Digits)
)

func generateID(length int, allowUppercase bool) string {
	var output [MaxLength]byte
	var randomBuffer [MaxLength * 2]byte
	random := randomBuffer[:length*2]
	outputOffset := 0
	for outputOffset < length {
		rand.Read(random)
		outputOffset = fillStructuredID(
			output[:length],
			outputOffset,
			random,
			allowUppercase,
		)
	}
	return string(output[:length])
}

func fillStructuredID(
	output []byte,
	outputOffset int,
	random []byte,
	allowUppercase bool,
) int {
	letters := Letters
	if allowUppercase {
		letters = LettersWithUppercase
	}
	letterLimit := 256 - 256%len(letters)

	for _, value := range random {
		if outputOffset == len(output) {
			break
		}
		isDigit := outputOffset%3 == 2
		if !isDigit && int(value) >= letterLimit {
			continue
		}

		if isDigit {
			output[outputOffset] = Digits[value&7]
		} else {
			output[outputOffset] = letters[int(value)%len(letters)]
		}
		outputOffset++
	}
	return outputOffset
}

func makeAlphabetLookup(alphabet string) [256]bool {
	var lookup [256]bool
	for index := range alphabet {
		lookup[alphabet[index]] = true
	}
	return lookup
}

func isValidID(value string, allowUppercase bool) bool {
	return isValidLength(len(value)) && hasValidIDFormat(value, allowUppercase)
}

func hasValidIDFormat(value string, allowUppercase bool) bool {
	letters := &letterLookup
	if allowUppercase {
		letters = &uppercaseLookup
	}
	index := 0
	for index+2 < len(value) {
		if !letters[value[index]] ||
			!letters[value[index+1]] ||
			!digitLookup[value[index+2]] {
			return false
		}
		index += 3
	}
	for index < len(value) {
		if !letters[value[index]] {
			return false
		}
		index++
	}
	return true
}
