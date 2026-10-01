// Package tidyid generates secure, human-friendly identifiers with a fixed
// two-letters-and-one-digit (LLD) rhythm.
//
// Every ID uses fresh bytes from the operating system CSPRNG through
// crypto/rand. Random bytes are never pooled or reused, letter sampling is
// unbiased, and a random-source failure terminates the process without a
// predictable fallback.
package tidyid
