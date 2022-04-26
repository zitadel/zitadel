package crypto

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/errors"
)

const (
	TypeEncryption CryptoType = iota
	TypeHash
)

type Crypto interface {
	Algorithm() string
}

type EncryptionAlgorithm interface {
	Crypto
	EncryptionKeyID() string
	DecryptionKeyIDs() []string
	Encrypt(value []byte) ([]byte, error)
	Decrypt(hashed []byte, keyID string) ([]byte, error)
	DecryptString(hashed []byte, keyID string) (string, error)
}

type HashAlgorithm interface {
	Crypto
	Hash(value []byte) ([]byte, error)
	CompareHash(hashed, comparer []byte) error
}

type CryptoValue struct {
	CryptoType CryptoType
	Algorithm  string
	KeyID      string
	Crypted    []byte
}

func (c *CryptoValue) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CryptoValue) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, c)
	}
	if s, ok := src.(string); ok {
		return json.Unmarshal([]byte(s), c)
	}
	return nil
}

type CryptoType int

func Crypt(value []byte, c Crypto) (*CryptoValue, error) {
	switch alg := c.(type) {
	case EncryptionAlgorithm:
		return Encrypt(value, alg)
	case HashAlgorithm:
		return Hash(value, alg)
	}
	return nil, errors.ThrowInternal(nil, "CRYPT-r4IaHZ", "algorithm not supported")
}

func Encrypt(value []byte, alg EncryptionAlgorithm) (*CryptoValue, error) {
	encrypted, err := alg.Encrypt(value)
	if err != nil {
		return nil, errors.ThrowInternal(err, "CRYPT-qCD0JB", "error encrypting value")
	}
	return &CryptoValue{
		CryptoType: TypeEncryption,
		Algorithm:  alg.Algorithm(),
		KeyID:      alg.EncryptionKeyID(),
		Crypted:    encrypted,
	}, nil
}

func Decrypt(value *CryptoValue, alg EncryptionAlgorithm) ([]byte, error) {
	if err := checkEncryptionAlgorithm(value, alg); err != nil {
		return nil, err
	}
	return alg.Decrypt(value.Crypted, value.KeyID)
}

func DecryptString(value *CryptoValue, alg EncryptionAlgorithm) (string, error) {
	if err := checkEncryptionAlgorithm(value, alg); err != nil {
		return "", err
	}
	return alg.DecryptString(value.Crypted, value.KeyID)
}

func checkEncryptionAlgorithm(value *CryptoValue, alg EncryptionAlgorithm) error {
	if value.Algorithm != alg.Algorithm() {
		return errors.ThrowInvalidArgument(nil, "CRYPT-Nx7XlT", "value was encrypted with a different key")
	}
	for _, id := range alg.DecryptionKeyIDs() {
		if id == value.KeyID {
			return nil
		}
	}
	return errors.ThrowInvalidArgument(nil, "CRYPT-Kq12vn", "value was encrypted with a different key")
}

func Hash(value []byte, alg HashAlgorithm) (*CryptoValue, error) {
	hashed, err := alg.Hash(value)
	if err != nil {
		return nil, errors.ThrowInternal(err, "CRYPT-rBVaJU", "error hashing value")
	}
	return &CryptoValue{
		CryptoType: TypeHash,
		Algorithm:  alg.Algorithm(),
		Crypted:    hashed,
	}, nil
}

func CompareHash(value *CryptoValue, comparer []byte, alg HashAlgorithm) error {
	if value.Algorithm != alg.Algorithm() {
		return errors.ThrowInvalidArgument(nil, "CRYPT-HF32f", "value was hashed with a different algorithm")
	}
	return alg.CompareHash(value.Crypted, comparer)
}
