package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"unicode/utf8"

	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ EncryptionAlgorithm = (*AESCrypto)(nil)

type AESCrypto struct {
	keys            map[string]string
	encryptionKeyID string
	keyIDs          []string
}

func NewAESCrypto(config *KeyConfig, keyStorage KeyStorage) (*AESCrypto, error) {
	keys, ids, err := LoadKeys(config, keyStorage)
	if err != nil {
		return nil, err
	}
	return &AESCrypto{
		keys:            keys,
		encryptionKeyID: config.EncryptionKeyID,
		keyIDs:          ids,
	}, nil
}

func (a *AESCrypto) Algorithm() string {
	return "aes"
}

func (a *AESCrypto) Encrypt(value []byte) ([]byte, error) {
	return EncryptAES(value, a.encryptionKey())
}

func (a *AESCrypto) Decrypt(value []byte, keyID string) ([]byte, error) {
	key, err := a.decryptionKey(keyID)
	if err != nil {
		return nil, err
	}
	return DecryptAES(value, key)
}

// DecryptString decrypts the value using the key identified by keyID.
// When the decrypted value contains non-UTF8 characters an error is returned.
func (a *AESCrypto) DecryptString(value []byte, keyID string) (string, error) {
	b, err := a.Decrypt(value, keyID)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", zerrors.ThrowPreconditionFailed(err, "CRYPT-hiCh0", "non-UTF-8 in decrypted string")
	}

	return string(b), nil
}

func (a *AESCrypto) EncryptionKeyID() string {
	return a.encryptionKeyID
}

func (a *AESCrypto) DecryptionKeyIDs() []string {
	return a.keyIDs
}

func (a *AESCrypto) encryptionKey() string {
	return a.keys[a.encryptionKeyID]
}

func (a *AESCrypto) decryptionKey(keyID string) (string, error) {
	key, ok := a.keys[keyID]
	if !ok {
		return "", zerrors.ThrowNotFound(nil, "CRYPT-nkj1s", "unknown key id")
	}
	return key, nil
}

func EncryptAESString(data string, key string) (string, error) {
	encrypted, err := EncryptAES([]byte(data), key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func EncryptAES(plainText []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	maxSize := 64 * 1024 * 1024
	if len(plainText) > maxSize {
		return nil, zerrors.ThrowPreconditionFailedf(nil, "CRYPT-AGg4t3", "data too large, max bytes: %v", maxSize)
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func DecryptAESString(data string, key string) (string, error) {
	text, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", nil
	}
	decrypted, err := DecryptAES(text, key)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func DecryptAES(text []byte, key string) ([]byte, error) {
	cipherText := make([]byte, len(text))
	copy(cipherText, text)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = zerrors.ThrowPreconditionFailed(nil, "CRYPT-23kH1", "cipher text block too short")
		return nil, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, err
}

// AES256GCMCrypto implements the EncryptionAlgorithm interface using AES-256-GCM for encryption and decryption.
type AES256GCMCrypto struct {
	keys            map[string]string
	encryptionKeyID string
	keyIDs          []string
	fallbackDecrypt func(value []byte, key string) ([]byte, error)
}

var _ EncryptionAlgorithm = (*AES256GCMCrypto)(nil)

func NewAES256GCMCrypto(config *KeyConfig, keyStorage KeyStorage, options ...AES256GCMCryptoOption) (*AES256GCMCrypto, error) {
	keys, ids, err := LoadKeys(config, keyStorage)
	if err != nil {
		return nil, err
	}
	crypto := &AES256GCMCrypto{
		keys:            keys,
		encryptionKeyID: config.EncryptionKeyID,
		keyIDs:          ids,
	}
	for _, option := range options {
		option(crypto)
	}
	return crypto, nil
}

func (a *AES256GCMCrypto) Algorithm() string {
	return string(jose.A256GCMKW)
}

func (a *AES256GCMCrypto) EncryptString(value []byte) (string, error) {
	encrypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{
		Algorithm: jose.A256GCMKW,
		Key:       a.encryptionKey(),
		KeyID:     a.encryptionKeyID,
	}, nil)
	if err != nil {
		return "", zerrors.ThrowInternal(err, "CRYPTO-Woox2", "Errors.Internal")
	}

	encrypted, err := encrypter.Encrypt(value)
	if err != nil {
		return "", zerrors.ThrowInternal(err, "CRYPTO-Woox3", "Errors.Internal")
	}

	serialized, err := encrypted.CompactSerialize()
	if err != nil {
		return "", zerrors.ThrowInternal(err, "CRYPT-Woox4", "Errors.Internal")
	}
	return serialized, nil
}

func (a *AES256GCMCrypto) Encrypt(value []byte) ([]byte, error) {
	encrypted, err := a.EncryptString(value)
	if err != nil {
		return nil, err
	}
	return []byte(encrypted), nil
}

// DecryptString decrypts the value using the key identified by keyID.
// When the decrypted value contains non-UTF8 characters an error is returned.
func (a *AES256GCMCrypto) DecryptString(value []byte, keyID string) (string, error) {
	b, err := a.Decrypt(value, keyID)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(b) {
		return "", zerrors.ThrowPreconditionFailed(err, "CRYPT-Koh1u", "non-UTF-8 in decrypted string")
	}
	return string(b), nil
}

func (a *AES256GCMCrypto) Decrypt(value []byte, keyID string) ([]byte, error) {
	key, err := a.decryptionKey(keyID)
	if err != nil {
		return nil, err
	}
	decrypted, err := a.decryptAS256GCM(value, key)
	if err == nil {
		return decrypted, nil
	}
	if a.fallbackDecrypt == nil {
		return nil, err
	}
	decrypted, fbErr := a.fallbackDecrypt(value, key)
	if fbErr == nil {
		return decrypted, nil
	}
	return nil, errors.Join(err, fbErr)
}

func (a *AES256GCMCrypto) decryptAS256GCM(value []byte, key string) ([]byte, error) {
	jwe, err := jose.ParseEncrypted(string(value), []jose.KeyAlgorithm{jose.A256GCMKW}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "CRYPT-ha6Oh", "malformed encypted value")
	}
	decrypted, err := jwe.Decrypt(key)
	if err != nil {
		return nil, zerrors.ThrowUnauthenticated(err, "CRYPT-OhN2u", "failed to decrypt value")
	}
	return decrypted, nil
}

func (a *AES256GCMCrypto) EncryptionKeyID() string {
	return a.encryptionKeyID
}

func (a *AES256GCMCrypto) DecryptionKeyIDs() []string {
	return a.keyIDs
}

func (a *AES256GCMCrypto) encryptionKey() string {
	return a.keys[a.encryptionKeyID]
}

func (a *AES256GCMCrypto) decryptionKey(keyID string) (string, error) {
	key, ok := a.keys[keyID]
	if !ok {
		return "", zerrors.ThrowNotFound(nil, "CRYPT-nkj1s", "unknown key id")
	}
	return key, nil
}

type AES256GCMCryptoOption func(*AES256GCMCrypto)

// WithAES256GCMCryptoFallbackDecrypt adds a fallback decryption function to the AES256GCMCrypto.
// This is used to support decrypting values that were encrypted with a different algorithm, such as AES-CFB, before migrating to AES-256-GCM.
func WithAES256GCMCryptoFallbackDecrypt(fallback func(value []byte, key string) ([]byte, error)) AES256GCMCryptoOption {
	return func(c *AES256GCMCrypto) {
		c.fallbackDecrypt = fallback
	}
}
