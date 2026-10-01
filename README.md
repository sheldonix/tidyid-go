<p align="center">
  <img src="docs/media/logo-128.png" alt="TidyID" height="64">
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/sheldonix/tidyid-go/v2"><img src="https://pkg.go.dev/badge/github.com/sheldonix/tidyid-go/v2.svg" alt="Go Reference"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.24%2B-00ADD8.svg?logo=go&amp;logoColor=white" alt="Go"></a>
  <a href="https://github.com/sheldonix/tidyid-go/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
</p>

A tiny, secure, and human-friendly ID generator for Go.

- **Tiny** Zero third-party runtime dependencies and a compact package-level API.
- **Secure** Uses Go's `crypto/rand` system CSPRNG with unbiased sampling and no weak fallback. Generate independently across goroutines, processes, and hosts; enforce absolute uniqueness with a database constraint.
- **Human-friendly** Creates letter-first, lowercase-alphanumeric IDs by default with a fixed `LLD` rhythm—no accidental long words, punctuation, or ambiguous characters. Easy to read, type, and transcribe; ready for URLs, filenames, database/cache/object-storage keys, HTML/CSS IDs, command lines, logs, and more.

```go
import "github.com/sheldonix/tidyid-go/v2"

id1, _ := tidyid.Generate(32, false) // zu5wk9fb6bz2xf8mc8gx5cp3zz7nj4an (length = 32)
id2, _ := tidyid.Generate(16, false) // ud4wp9zw2rf4dw4t                 (length = 16)
id3, _ := tidyid.Generate(10, false) // kz9ac6qr8f                       (length = 10)
id4, _ := tidyid.Generate(10, true)  // NJ2Db5vz5k                       (length = 10, allowUppercase = true)
```

Explicitly passing the length in application code is strongly recommended, even when using the default of 32.

## Install

```sh
go get github.com/sheldonix/tidyid-go/v2
```

## CLI

Install TidyID:

```sh
go install github.com/sheldonix/tidyid-go/v2/cmd/tidyid@latest
```

Generate IDs:

```sh
tidyid
# rw6pa5nu2tf3nb7tg2kq2fk9ju3mr4tv (length = 32)

tidyid -s 16
# we6yc3xk6cg5wj2y (length = 16)

tidyid -s 10 -u
# yM9vT6NB2T (length = 10, allowUppercase = true)
```

Use `--size` or `-s` to set the length. Use `--allow-uppercase` or `-u`
to allow uppercase letters.

## Format

By default, IDs repeat two lowercase letters followed by one digit (`LLD`):

```text
gr7 gb5 kd4
```

| Characters | Positions | Alphabet |
| --- | --- | --- |
| Letters | First two of each group | `abcdefghjkmnpqrtuvwxyz` |
| Digits | Every third character | `23456789` |

- Every ID starts with a letter.
- `i`, `l`, `o`, `s`, `0`, and `1` are excluded to reduce visual and handwritten ambiguity.
- The pattern prevents long letter sequences and needs no escaping in URL paths, filenames, or HTML/CSS IDs.
- In default mode, typing needs no Shift key, `_`, `-`, or other punctuation.

`allowUppercase` defaults to `false`. Set it to `true` to sample letter
positions from the combined uppercase and lowercase alphabet.

## API

| API | Description |
| --- | --- |
| `New() (string, error)` | Generate a 32-character lowercase ID. |
| `Generate(length int, allowUppercase bool) (string, error)` | Generate a 3–256 character ID. |
| `IsValidID(value string, allowUppercase bool) bool` | Check the format and supported length range. |
| `IsValidIDOfLength(value string, length int, allowUppercase bool) bool` | Check the format and an exact valid length. |
| `EnsureValidID(value string, allowUppercase bool) error` | Return `InvalidIDFormatError` for an invalid value. |
| `EnsureValidIDOfLength(value string, length int, allowUppercase bool) error` | Return `InvalidIDLengthError` or `InvalidIDFormatError`. |
| `IDCapacity(length int, allowUppercase bool) (*big.Int, error)` | Return the exact ID space as a new `*big.Int`. |
| `IDEntropy(length int, allowUppercase bool) (float64, error)` | Return entropy in bits. |

All package-level generation and validation functions are safe for concurrent use.

Constants: `Letters`, `LettersWithUppercase`, `Digits`, `DefaultLength`, `MinLength`, `MaxLength`, `Version`.

Errors: `InvalidIDLengthError`, `InvalidIDFormatError`.

## Security

- **Unpredictability** Every generation calls `crypto/rand.Read`, backed by the operating system cryptographically secure random number generator. TidyID never uses `math/rand`.
- **Uniformity** Letter positions use rejection sampling, while digit positions use an exact eight-way mapping. Both avoid modulo bias, so every valid ID of the same length and mode has equal probability.
- **Collision-aware** Choose a length for your scale to make collisions extremely unlikely. Use a database `PRIMARY KEY` or `UNIQUE` constraint when absolute uniqueness must be enforced.

  > **Default mode (`allowUppercase = false`)**
  >
  > | Length | Capacity | Entropy |
  > | ---: | ---: | ---: |
  > | 8 | 7,256,313,856 | 32.76 bits |
  > | 10 | 1,277,111,238,656 | 40.22 bits |
  > | 12 | 224,771,578,003,456 | 47.68 bits |
  > | 16 | 19,146,942,100,646,395,904 | 64.05 bits |
  > | 23 | 6,315,282,784,770,463,143,393,492,992 | 92.35 bits |
  > | 32 | 366,605,391,805,505,419,895,548,144,464,707,977,216 | 128.11 bits |

  > **Uppercase allowed (`allowUppercase = true`)**
  >
  > | Length | Capacity | Entropy |
  > | ---: | ---: | ---: |
  > | 8 | 464,404,086,784 | 38.76 bits |
  > | 10 | 163,470,238,547,968 | 47.22 bits |
  > | 12 | 57,541,523,968,884,736 | 55.68 bits |
  > | 16 | 39,212,937,422,123,818,811,392 | 75.05 bits |
  > | 23 | 413,878,372,582,717,072,565,435,956,723,712 | 108.35 bits |
  > | 32 | 1,537,654,461,271,398,604,689,577,164,520,902,527,668,977,664 | 150.11 bits |

  Use 16 or more characters for large public datasets. For security tokens, choose the length based on your threat model. A 32-character TidyID provides 128.11 bits of entropy by default, or 150.11 bits with `allowUppercase = true`.

## Database uniqueness

Use a primary key or unique constraint. Insert first and retry only an ID conflict—never query before inserting:

```go
for attempt := 0; attempt < 128; attempt++ {
	id, err := tidyid.Generate(16, false)
	if err != nil {
		return "", err
	}

	result, err := db.ExecContext(
		ctx,
		"INSERT INTO resources (id) VALUES ($1) ON CONFLICT (id) DO NOTHING",
		id,
	)
	if err != nil {
		return "", err
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if inserted == 1 {
		return id, nil
	}
}
return "", errors.New("unable to insert a resource with a unique TidyID")
```

Propagate network, permission, transaction, and non-ID constraint errors.

## Requirements

- Go `>=1.24`

## License

MIT
