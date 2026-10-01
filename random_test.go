package tidyid

import (
	"slices"
	"testing"
)

func TestStructuredSamplingConsumesOneFreshBlock(t *testing.T) {
	output := make([]byte, 32)
	written := fillStructuredID(output, 0, make([]byte, len(output)*2), false)
	if written != len(output) {
		t.Fatalf("written = %d, want %d", written, len(output))
	}
	if string(output) != "aa2aa2aa2aa2aa2aa2aa2aa2aa2aa2aa" {
		t.Fatalf("output = %q", output)
	}
}

func TestLetterAcceptanceRangesAreExactAndUniform(t *testing.T) {
	for _, test := range []struct {
		alphabet  string
		uppercase bool
	}{
		{Letters, false},
		{LettersWithUppercase, true},
	} {
		alphabet := test.alphabet
		limit := 256 - 256%len(alphabet)
		counts := make([]int, len(alphabet))
		rejected := 0
		for input := 0; input < 256; input++ {
			output := make([]byte, 1)
			written := fillStructuredID(output, 0, []byte{byte(input)}, test.uppercase)
			if input >= limit {
				rejected++
				if written != 0 {
					t.Fatalf("input %d was accepted", input)
				}
				continue
			}
			if written != 1 {
				t.Fatalf("input %d was rejected", input)
			}
			index := slices.Index([]byte(alphabet), output[0])
			if index < 0 {
				t.Fatalf("input %d produced invalid letter %q", input, output[0])
			}
			counts[index]++
		}
		if rejected != 256-limit {
			t.Fatalf("rejected %d inputs, want %d", rejected, 256-limit)
		}
		for index, count := range counts {
			if count != limit/len(alphabet) {
				t.Fatalf("%q accepted %d times, want %d", alphabet[index], count, limit/len(alphabet))
			}
		}
	}

	digitCounts := make([]int, len(Digits))
	for input := 0; input < 256; input++ {
		output := make([]byte, 3)
		written := fillStructuredID(output, 2, []byte{byte(input)}, false)
		if written != len(output) {
			t.Fatalf("input %d wrote through position %d", input, written)
		}
		index := slices.Index([]byte(Digits), output[2])
		if index < 0 {
			t.Fatalf("input %d produced invalid digit %q", input, output[2])
		}
		digitCounts[index]++
	}
	for index, count := range digitCounts {
		if count != 256/len(Digits) {
			t.Fatalf("%q accepted %d times, want %d", Digits[index], count, 256/len(Digits))
		}
	}
}

func TestRejectionSamplingContinuesWithNextBlock(t *testing.T) {
	output := make([]byte, 3)
	written := fillStructuredID(
		output,
		0,
		[]byte{255, 255, 255, 255, 255, 255},
		false,
	)
	if written != 0 {
		t.Fatalf("rejected block wrote %d characters", written)
	}
	written = fillStructuredID(output, written, make([]byte, 6), false)
	if written != len(output) {
		t.Fatalf("written = %d, want %d", written, len(output))
	}
	if string(output) != "aa2" {
		t.Fatalf("output = %q, want aa2", output)
	}
}

func TestUppercaseSamplingAndStructure(t *testing.T) {
	output := make([]byte, 3)
	written := fillStructuredID(output, 0, []byte{0, 22, 0}, true)
	if written != len(output) {
		t.Fatalf("written = %d, want %d", written, len(output))
	}
	if string(output) != "Aa2" {
		t.Fatalf("output = %q, want Aa2", output)
	}
}
