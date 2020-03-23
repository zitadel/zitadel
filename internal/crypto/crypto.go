package crypto

import (
	"github.com/caos/zitadel/internal/errors"
)

const (
	TypeEncryption CryptoType = iota
	TypeHash
)

type Crypto interface {
	Algorithm() string
}

type EncryptionAlg interface {
	Crypto
	EncryptionKeyID() string
	DecryptionKeyIDs() []string
	Encrypt(value []byte) ([]byte, error)
	Decrypt(hashed []byte, keyID string) ([]byte, error)
	DecryptString(hashed []byte, keyID string) (string, error)
}

type HashAlg interface {
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

type CryptoType int

func Crypt(value []byte, c Crypto) (*CryptoValue, error) {
	switch alg := c.(type) {
	case EncryptionAlg:
		return Encrypt(value, alg)
	case HashAlg:
		return Hash(value, alg)
	}
	return nil, errors.ThrowInternal(nil, "CRYPT-r4IaHZ", "algorithm not supported")
}

func Encrypt(value []byte, alg EncryptionAlg) (*CryptoValue, error) {
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

func Decrypt(value *CryptoValue, alg EncryptionAlg) ([]byte, error) {
	if err := checkEncAlg(value, alg); err != nil {
		return nil, err
	}
	return alg.Decrypt(value.Crypted, value.KeyID)
}

func DecryptString(value *CryptoValue, alg EncryptionAlg) (string, error) {
	if err := checkEncAlg(value, alg); err != nil {
		return "", err
	}
	return alg.DecryptString(value.Crypted, value.KeyID)
}

func checkEncAlg(value *CryptoValue, alg EncryptionAlg) error {
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

func Hash(value []byte, alg HashAlg) (*CryptoValue, error) {
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

func CompareHash(value *CryptoValue, comparer []byte, alg HashAlg) error {
	return alg.CompareHash(value.Crypted, comparer)
}
