package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/caos/zitadel/internal/errors"
)

var _ EncryptionAlgorithm = (*AESCrypto)(nil)

type AESCrypto struct {
	keys            map[string]string
	encryptionKeyID string
	keyIDs          []string
}

func NewAESCrypto(config *KeyConfig) (*AESCrypto, error) {
	keys, ids, err := LoadKeys(config)
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

func (a *AESCrypto) DecryptString(value []byte, keyID string) (string, error) {
	key, err := a.decryptionKey(keyID)
	if err != nil {
		return "", err
	}
	b, err := DecryptAES(value, key)
	if err != nil {
		return "", err
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
		return "", errors.ThrowNotFound(nil, "CRYPT-nkj1s", "unknown key id")
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
		err = errors.ThrowPreconditionFailed(nil, "CRYPT-23kH1", "cipher text block too short")
		return nil, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, err
}
