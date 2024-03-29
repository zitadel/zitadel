package crypto

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TypeEncryption CryptoType = iota
	TypeHash                  // Depcrecated: use [passwap.Swapper] instead
)

type EncryptionAlgorithm interface {
	Algorithm() string
	EncryptionKeyID() string
	DecryptionKeyIDs() []string
	Encrypt(value []byte) ([]byte, error)
	Decrypt(hashed []byte, keyID string) ([]byte, error)
	DecryptString(hashed []byte, keyID string) (string, error)
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

func Crypt(value []byte, alg EncryptionAlgorithm) (*CryptoValue, error) {
	return Encrypt(value, alg)
}

func Encrypt(value []byte, alg EncryptionAlgorithm) (*CryptoValue, error) {
	encrypted, err := alg.Encrypt(value)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "CRYPT-qCD0JB", "error encrypting value")
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
		return zerrors.ThrowInvalidArgument(nil, "CRYPT-Nx7XlT", "value was encrypted with a different key")
	}
	for _, id := range alg.DecryptionKeyIDs() {
		if id == value.KeyID {
			return nil
		}
	}
	return zerrors.ThrowInvalidArgument(nil, "CRYPT-Kq12vn", "value was encrypted with a different key")
}

func CheckToken(alg EncryptionAlgorithm, token string, content string) error {
	if token == "" {
		return zerrors.ThrowPermissionDenied(nil, "CRYPTO-Sfefs", "Errors.Intent.InvalidToken")
	}
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return zerrors.ThrowPermissionDenied(err, "CRYPTO-Swg31", "Errors.Intent.InvalidToken")
	}
	decryptedToken, err := alg.DecryptString(data, alg.EncryptionKeyID())
	if err != nil {
		return zerrors.ThrowPermissionDenied(err, "CRYPTO-Sf4gt", "Errors.Intent.InvalidToken")
	}
	if decryptedToken != content {
		return zerrors.ThrowPermissionDenied(nil, "CRYPTO-CRYPTO", "Errors.Intent.InvalidToken")
	}
	return nil
}

// SecretOrEncodedHash returns the Crypted value from legacy [CryptoValue] if it is not nil.
// otherwise it will returns the encoded hash string.
func SecretOrEncodedHash(secret *CryptoValue, encoded string) string {
	if secret != nil {
		return string(secret.Crypted)
	}
	return encoded
}
