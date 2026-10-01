// Package tidyid generates secure, random, and human-friendly identifiers.
//
// Every ID uses fresh bytes from the operating system CSPRNG through
// crypto/rand. Random bytes are never pooled or reused, letter sampling is
// unbiased, and a random-source failure terminates the process without a
// predictable fallback.
package tidyid
