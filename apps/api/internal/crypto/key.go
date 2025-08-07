package crypto

import (
	"crypto/rand"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type KeyConfig struct {
	EncryptionKeyID  string
	DecryptionKeyIDs []string
}

type Keys map[string]string

type Key struct {
	ID    string
	Value string
}

func NewKey(id string) (*Key, error) {
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return nil, err
	}
	return &Key{
		ID:    id,
		Value: string(randBytes),
	}, nil
}

func LoadKey(id string, keyStorage KeyStorage) (string, error) {
	key, err := keyStorage.ReadKey(id)
	if err != nil {
		return "", err
	}
	return key.Value, nil
}

func LoadKeys(config *KeyConfig, keyStorage KeyStorage) (Keys, []string, error) {
	if config == nil {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "CRYPT-dJK8s", "config must not be nil")
	}
	readKeys, err := keyStorage.ReadKeys()
	if err != nil {
		return nil, nil, err
	}
	keys := make(Keys)
	ids := make([]string, 0, len(config.DecryptionKeyIDs)+1)
	if config.EncryptionKeyID != "" {
		key, ok := readKeys[config.EncryptionKeyID]
		if !ok {
			return nil, nil, zerrors.ThrowInternalf(nil, "CRYPT-v2Kas", "encryption key %s not found", config.EncryptionKeyID)
		}
		keys[config.EncryptionKeyID] = key
		ids = append(ids, config.EncryptionKeyID)
	}
	for _, id := range config.DecryptionKeyIDs {
		key, ok := readKeys[id]
		if !ok {
			logging.Errorf("description key %s not found", id)
			continue
		}
		keys[id] = key
		ids = append(ids, id)
	}
	return keys, ids, nil
}
