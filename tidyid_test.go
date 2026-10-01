package tidyid

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"testing"
)

func TestConstantsMatchReferenceFormat(t *testing.T) {
	if Letters != "abcdefghjkmnpqrtuvwxyz" {
		t.Fatalf("Letters = %q", Letters)
	}
	if Digits != "23456789" {
		t.Fatalf("Digits = %q", Digits)
	}
	if LettersWithUppercase != "ABCDEFGHJKMNPQRTUVWXYZabcdefghjkmnpqrtuvwxyz" {
		t.Fatalf("LettersWithUppercase = %q", LettersWithUppercase)
	}
	if DefaultLength != 32 || MinLength != 3 || MaxLength != 256 {
		t.Fatalf("length constants = %d, %d, %d", DefaultLength, MinLength, MaxLength)
	}
	if len(Digits) != 1<<digitAlphabetBits {
		t.Fatalf("digit alphabet length = %d, want %d", len(Digits), 1<<digitAlphabetBits)
	}
}

func TestGenerateBoundariesAndModes(t *testing.T) {
	for _, length := range []int{3, 8, 10, 16, 32, 256} {
		for _, uppercase := range []bool{false, true} {
			for sample := 0; sample < 25; sample++ {
				value, err := Generate(length, uppercase)
				if err != nil {
					t.Fatalf("Generate(%d, %t): %v", length, uppercase, err)
				}
				if !IsValidIDOfLength(value, length, uppercase) {
					t.Fatalf("Generate(%d, %t) = %q", length, uppercase, value)
				}
			}
		}
	}
}

func TestNewUsesDefaultFormat(t *testing.T) {
	value, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if !IsValidIDOfLength(value, DefaultLength, false) {
		t.Fatalf("New() = %q", value)
	}
}

func TestGenerateRejectsInvalidLengths(t *testing.T) {
	for _, length := range []int{-1, 0, 1, 2, 257, int(^uint(0) >> 1)} {
		_, err := Generate(length, false)
		var target InvalidIDLengthError
		if !errors.As(err, &target) {
			t.Fatalf("Generate(%d, false) error = %T %v", length, err, err)
		}
		if err.Error() != "length must be between 3 and 256" {
			t.Fatalf("Generate(%d, false) error = %q", length, err)
		}
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		value     string
		uppercase bool
		valid     bool
	}{
		{"mk7qw2xy", false, true},
		{"mk7qw2x", false, true},
		{"mk7qw2x9", false, false},
		{"m27qw2xy", false, false},
		{"mk7q52xy", false, false},
		{"MK7QW2XY", false, false},
		{"MK7QW2XY", true, true},
		{"aZ2", true, true},
		{"a2Z", true, false},
		{"aZQ", true, false},
		{"aI2", true, false},
		{"aZ0", true, false},
		{"mk", false, false},
		{"\u00e9a2", false, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%q/%t", test.value, test.uppercase), func(t *testing.T) {
			if got := IsValidID(test.value, test.uppercase); got != test.valid {
				t.Fatalf("IsValidID(%q, %t) = %t", test.value, test.uppercase, got)
			}
		})
	}

	if !IsValidIDOfLength("mk7qw2xy", 8, false) {
		t.Fatal("valid exact length rejected")
	}
	if IsValidIDOfLength("mk7qw2xy", 16, false) {
		t.Fatal("wrong exact length accepted")
	}
	if IsValidIDOfLength("mk", 2, false) {
		t.Fatal("invalid requested length accepted")
	}
}

func TestValidationAtEveryLengthAndPosition(t *testing.T) {
	for length := MinLength; length <= MaxLength; length++ {
		lowercase := make([]byte, length)
		uppercase := make([]byte, length)
		for index := range lowercase {
			if (index+1)%3 == 0 {
				lowercase[index] = Digits[0]
				uppercase[index] = Digits[0]
			} else {
				lowercase[index] = Letters[0]
				uppercase[index] = LettersWithUppercase[0]
			}
		}

		if !IsValidID(string(lowercase), false) {
			t.Fatalf("valid lowercase ID rejected at length %d", length)
		}
		if !IsValidID(string(uppercase), true) {
			t.Fatalf("valid uppercase ID rejected at length %d", length)
		}
		if IsValidID(string(uppercase), false) {
			t.Fatalf("uppercase ID accepted in lowercase mode at length %d", length)
		}

		for index := range lowercase {
			lowercaseCharacter := lowercase[index]
			lowercase[index] = '0'
			if IsValidID(string(lowercase), false) {
				t.Fatalf("invalid lowercase byte accepted at length %d, position %d", length, index)
			}
			lowercase[index] = lowercaseCharacter

			uppercaseCharacter := uppercase[index]
			uppercase[index] = '0'
			if IsValidID(string(uppercase), true) {
				t.Fatalf("invalid uppercase byte accepted at length %d, position %d", length, index)
			}
			uppercase[index] = uppercaseCharacter
		}
	}
}

func TestEnsureValidIDReturnsExplicitErrors(t *testing.T) {
	if err := EnsureValidID("mk7qw2xy", false); err != nil {
		t.Fatal(err)
	}
	if err := EnsureValidIDOfLength("aZ2", 3, true); err != nil {
		t.Fatal(err)
	}

	var formatError InvalidIDFormatError
	if err := EnsureValidIDOfLength("invalid", 8, false); !errors.As(err, &formatError) {
		t.Fatalf("format error = %T %v", err, err)
	}
	var lengthError InvalidIDLengthError
	if err := EnsureValidIDOfLength("mk", 2, false); !errors.As(err, &lengthError) {
		t.Fatalf("length error = %T %v", err, err)
	}
}

func TestCapacityAndEntropy(t *testing.T) {
	tests := []struct {
		length    int
		uppercase bool
		capacity  string
		entropy   float64
	}{
		{8, false, "7256313856", 32.756589711823786},
		{16, false, "19146942100646395904", 64.05374780501026},
		{3, true, "15488", 13.918863237274595},
	}
	for _, test := range tests {
		capacity, err := IDCapacity(test.length, test.uppercase)
		if err != nil {
			t.Fatal(err)
		}
		if capacity.String() != test.capacity {
			t.Fatalf("IDCapacity(%d, %t) = %s", test.length, test.uppercase, capacity)
		}
		entropy, err := IDEntropy(test.length, test.uppercase)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(entropy-test.entropy) >= 1e-12 {
			t.Fatalf("IDEntropy(%d, %t) = %.16f", test.length, test.uppercase, entropy)
		}
	}

	capacity, err := IDCapacity(8, false)
	if err != nil {
		t.Fatal(err)
	}
	capacity.SetInt64(0)
	second, err := IDCapacity(8, false)
	if err != nil {
		t.Fatal(err)
	}
	if second.String() != "7256313856" {
		t.Fatal("IDCapacity returned shared mutable state")
	}

	for _, length := range []int{0, 2, 257} {
		if _, err := IDCapacity(length, false); err == nil {
			t.Fatalf("IDCapacity(%d, false) succeeded", length)
		}
		if _, err := IDEntropy(length, false); err == nil {
			t.Fatalf("IDEntropy(%d, false) succeeded", length)
		}
	}
}

func TestCapacityMatchesReferenceFormulaAtEveryLength(t *testing.T) {
	for length := MinLength; length <= MaxLength; length++ {
		for _, uppercase := range []bool{false, true} {
			actual, err := IDCapacity(length, uppercase)
			if err != nil {
				t.Fatal(err)
			}
			digitCount := length / 3
			letterBase := int64(len(Letters))
			if uppercase {
				letterBase = int64(len(LettersWithUppercase))
			}
			letters := new(big.Int).Exp(
				big.NewInt(letterBase),
				big.NewInt(int64(length-digitCount)),
				nil,
			)
			digits := new(big.Int).Exp(
				big.NewInt(int64(len(Digits))),
				big.NewInt(int64(digitCount)),
				nil,
			)
			expected := letters.Mul(letters, digits)
			if actual.Cmp(expected) != 0 {
				t.Fatalf(
					"IDCapacity(%d, %t) = %s, want %s",
					length,
					uppercase,
					actual,
					expected,
				)
			}
		}
	}
}

func TestGenerateIsConcurrentSafe(t *testing.T) {
	const generatedPerWorker = 1_000
	configurations := []struct {
		length    int
		uppercase bool
	}{
		{MinLength, false},
		{MinLength, true},
		{DefaultLength, false},
		{DefaultLength, true},
		{MaxLength, false},
		{MaxLength, true},
	}
	const workersPerConfiguration = 4
	workers := len(configurations) * workersPerConfiguration
	var wait sync.WaitGroup
	wait.Add(workers)
	errors := make(chan error, workers)
	for worker := 0; worker < workers; worker++ {
		go func(worker int) {
			defer wait.Done()
			configuration := configurations[worker%len(configurations)]
			for sample := 0; sample < generatedPerWorker; sample++ {
				value, err := Generate(
					configuration.length,
					configuration.uppercase,
				)
				if err != nil {
					errors <- err
					return
				}
				if !IsValidIDOfLength(
					value,
					configuration.length,
					configuration.uppercase,
				) {
					errors <- fmt.Errorf("invalid generated value %q", value)
					return
				}
			}
		}(worker)
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		t.Error(err)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for _, mode := range []struct {
		name      string
		uppercase bool
	}{
		{"lowercase", false},
		{"uppercase", true},
	} {
		b.Run(mode.name, func(b *testing.B) {
			for _, length := range []int{8, 10, 16, 32, 64, 256} {
				b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
					b.ReportAllocs()
					for index := 0; index < b.N; index++ {
						value, err := Generate(length, mode.uppercase)
						if err != nil {
							b.Fatal(err)
						}
						benchmarkValue = value
					}
				})
			}
		})
	}
}

func BenchmarkGenerateParallel(b *testing.B) {
	for _, mode := range []struct {
		name      string
		uppercase bool
	}{
		{"lowercase", false},
		{"uppercase", true},
	} {
		b.Run(mode.name, func(b *testing.B) {
			for _, length := range []int{
				MinLength,
				DefaultLength,
				MaxLength,
			} {
				b.Run(fmt.Sprintf("length=%d", length), func(b *testing.B) {
					b.ReportAllocs()
					b.RunParallel(func(parallel *testing.PB) {
						for parallel.Next() {
							_, err := Generate(length, mode.uppercase)
							if err != nil {
								b.Error(err)
								return
							}
						}
					})
				})
			}
		})
	}
}

func BenchmarkValidation(b *testing.B) {
	value := "ab2cd3ef4gh5jk6mn7pq8rt9uv2wx3yz"
	b.ReportAllocs()
	for index := 0; index < b.N; index++ {
		benchmarkValid = IsValidID(value, false)
	}
}

func BenchmarkIDEntropy(b *testing.B) {
	b.ReportAllocs()
	for index := 0; index < b.N; index++ {
		benchmarkEntropy, _ = IDEntropy(DefaultLength, false)
	}
}

func BenchmarkIDCapacity(b *testing.B) {
	b.ReportAllocs()
	for index := 0; index < b.N; index++ {
		benchmarkCapacity, _ = IDCapacity(DefaultLength, false)
	}
}

var (
	benchmarkValue    string
	benchmarkValid    bool
	benchmarkEntropy  float64
	benchmarkCapacity *big.Int
)
