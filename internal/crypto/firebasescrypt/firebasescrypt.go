// Package firebasescrypt verifies the modified scrypt of Firebase
// Authentication, as defined by https://github.com/firebase/scrypt.
package firebasescrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/zitadel/passwap/verifier"
	"golang.org/x/crypto/scrypt"
)

const (
	Identifier = "firebasescrypt"
	Prefix     = "$" + Identifier + "$"
)

// Firebase caps mem_cost at 14 and rounds at 8.
const (
	DefaultMinLN = 14
	DefaultMaxLN = 14
	DefaultMinR  = 8
	DefaultMaxR  = 8
)

// $firebasescrypt$ln=<mem_cost>,r=<rounds>$<salt>$<hash>$<salt_separator>$<signer_key>
const paramsFormat = "ln=%d,r=%d"

const (
	keyLen = 32
	// Firebase always uses p=1.
	parallelization = 1
)

type checker struct {
	ln int // log2 of N, the CPU/memory cost parameter
	r  int

	salt          []byte
	saltSeparator []byte
	hash          []byte
	signerKey     []byte
}

func parse(encoded string) (*checker, error) {
	if !strings.HasPrefix(encoded, Prefix) {
		return nil, nil
	}

	parts := strings.Split(encoded, "$")
	if len(parts) != 7 {
		return nil, fmt.Errorf("firebasescrypt parse: expected 7 parts, got %d", len(parts))
	}

	var c checker
	_, err := fmt.Sscanf(parts[2], paramsFormat, &c.ln, &c.r)
	if err != nil {
		return nil, fmt.Errorf("firebasescrypt parse: %w", err)
	}
	if fmt.Sprintf(paramsFormat, c.ln, c.r) != parts[2] {
		return nil, fmt.Errorf("firebasescrypt parse: params %q not in canonical form", parts[2])
	}
	if c.ln < 0 {
		return nil, fmt.Errorf("firebasescrypt parse: negative ln %d", c.ln)
	}

	c.salt, err = decodeBase64(parts[3])
	if err != nil {
		return nil, fmt.Errorf("firebasescrypt parse salt: %w", err)
	}
	c.hash, err = decodeBase64(parts[4])
	if err != nil {
		return nil, fmt.Errorf("firebasescrypt parse hash: %w", err)
	}
	c.saltSeparator, err = decodeBase64(parts[5])
	if err != nil {
		return nil, fmt.Errorf("firebasescrypt parse salt separator: %w", err)
	}
	c.signerKey, err = decodeBase64(parts[6])
	if err != nil {
		return nil, fmt.Errorf("firebasescrypt parse signer key: %w", err)
	}

	if len(c.hash) == 0 || len(c.hash) != len(c.signerKey) {
		return nil, errors.New("firebasescrypt parse: hash and signer key must be non-empty and of equal length")
	}

	return &c, nil
}

// Firebase CLI and console exports use standard base64.
// The Admin SDK uses URL-safe base64.
func decodeBase64(s string) ([]byte, error) {
	s = strings.TrimRight(s, "=")
	if strings.ContainsAny(s, "-_") {
		return base64.RawURLEncoding.Strict().DecodeString(s)
	}
	return base64.RawStdEncoding.Strict().DecodeString(s)
}

func (c *checker) validate(opts *ValidationOpts) error {
	if c.ln < opts.MinLN || c.ln > opts.MaxLN {
		return &verifier.BoundsError{
			Algorithm: Identifier,
			Param:     "LN",
			Min:       opts.MinLN,
			Max:       opts.MaxLN,
			Actual:    c.ln,
		}
	}
	if c.r < opts.MinR || c.r > opts.MaxR {
		return &verifier.BoundsError{
			Algorithm: Identifier,
			Param:     "R",
			Min:       opts.MinR,
			Max:       opts.MaxR,
			Actual:    c.r,
		}
	}
	return nil
}

// Makes the same library calls in the same order as CompareFirebaseScrypt
// of Ory Kratos (Apache-2.0), which was used as the reference.
// https://github.com/ory/kratos/blob/b86338da04a040247a07f46100a86dcfb3875909/hash/hash_comparator.go#L320-L351
func (c *checker) verify(pw string) (verifier.Result, error) {
	key, err := scrypt.Key([]byte(pw), slices.Concat(c.salt, c.saltSeparator), 1<<c.ln, c.r, parallelization, keyLen)
	if err != nil {
		return verifier.Fail, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return verifier.Fail, err
	}

	hash := make([]byte, len(c.signerKey))
	cipher.NewCTR(block, make([]byte, aes.BlockSize)).XORKeyStream(hash, c.signerKey)
	res := subtle.ConstantTimeCompare(hash, c.hash)

	return verifier.Result(res), nil
}

type ValidationOpts struct {
	MinLN int
	MaxLN int
	MinR  int
	MaxR  int
}

var DefaultValidationOpts = &ValidationOpts{
	MinLN: DefaultMinLN,
	MaxLN: DefaultMaxLN,
	MinR:  DefaultMinR,
	MaxR:  DefaultMaxR,
}

func checkValidationOpts(opts *ValidationOpts) *ValidationOpts {
	if opts == nil {
		return DefaultValidationOpts
	}
	o := *opts
	if o.MinLN <= 0 {
		o.MinLN = DefaultValidationOpts.MinLN
	}
	if o.MaxLN <= 0 {
		o.MaxLN = DefaultValidationOpts.MaxLN
	}
	if o.MinR <= 0 {
		o.MinR = DefaultValidationOpts.MinR
	}
	if o.MaxR <= 0 {
		o.MaxR = DefaultValidationOpts.MaxR
	}
	return &o
}

type Verifier struct {
	opts *ValidationOpts
}

func NewVerifier(opts *ValidationOpts) *Verifier {
	return &Verifier{
		opts: checkValidationOpts(opts),
	}
}

func (v *Verifier) Validate(encoded string) (verifier.Result, error) {
	c, err := parse(encoded)
	if err != nil || c == nil {
		return verifier.Skip, err
	}
	err = c.validate(v.opts)
	if err != nil {
		return verifier.Fail, err
	}
	return verifier.OK, nil
}

// Verify does not check parameter bounds, same as passwap verifiers.
func (v *Verifier) Verify(encoded, password string) (verifier.Result, error) {
	c, err := parse(encoded)
	if err != nil || c == nil {
		return verifier.Skip, err
	}

	return c.verify(password)
}
