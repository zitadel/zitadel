package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
)

// HMACValue is a randomly salted HMAC-SHA256 of a payload.
// It allows checking whether a payload was seen before,
// without having to store the payload itself.
//
// It is intended for low entropy secrets, such as 6 digit TOTP codes:
// the random salt makes each value unique, even for an equal payload,
// which counters rainbow table usage and prevents that two stored values
// can be recognized as being derived from the same payload.
//
// As the salt is stored alongside the hash, an HMACValue is not suitable
// as a replacement for retrievable secrets. Use [CryptoValue] for that.
//
// Both fields are JSON serializable, so a value can be stored in an event payload.
type HMACValue struct {
	Hash []byte `json:"hash"`
	Salt []byte `json:"salt"`
}

// NewHMACValue returns the HMAC of payload, using a newly generated random salt.
// Two calls with an equal payload always return different values.
func NewHMACValue(payload string) *HMACValue {
	salt := make([]byte, 16)
	// [rand.Read] never returns an error, it panics if the system source fails.
	rand.Read(salt)
	hasher := hmac.New(sha256.New, salt)
	hasher.Write([]byte(payload))

	return &HMACValue{
		Hash: hasher.Sum(nil),
		Salt: salt,
	}
}

// Equal reports whether payload is the one v was created from.
// The hashes are compared in constant time.
func (v *HMACValue) Equal(payload string) bool {
	hasher := hmac.New(sha256.New, v.Salt)
	hasher.Write([]byte(payload))
	expectedHash := hasher.Sum(nil)
	return hmac.Equal(v.Hash, expectedHash)
}
