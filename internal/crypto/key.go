package crypto

import (
	"os"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/config"
	"github.com/zitadel/zitadel/internal/errors"
)

const (
	ZitadelKeyPath = "ZITADEL_KEY_PATH"
)

type KeyConfig struct {
	EncryptionKeyID  string
	DecryptionKeyIDs []string
	Path             string
}

type Keys map[string]string

func ReadKeys(path string) (Keys, error) {
	if path == "" {
		path = os.Getenv(ZitadelKeyPath)
		if path == "" {
			return nil, errors.ThrowInvalidArgument(nil, "CRYPT-56lka", "no path set")
		}
	}
	keys := new(Keys)
	err := config.Read(keys, path)
	return *keys, err
}

func LoadKey(config *KeyConfig, id string) (string, error) {
	keys, _, err := LoadKeys(config)
	if err != nil {
		return "", err
	}
	return keys[id], nil
}

func LoadKeys(config *KeyConfig) (map[string]string, []string, error) {
	if config == nil {
		return nil, nil, errors.ThrowInvalidArgument(nil, "CRYPT-dJK8s", "config must not be nil")
	}
	readKeys, err := ReadKeys(config.Path)
	if err != nil {
		return nil, nil, err
	}
	keys := make(map[string]string)
	ids := make([]string, 0, len(config.DecryptionKeyIDs)+1)
	if config.EncryptionKeyID != "" {
		key, ok := readKeys[config.EncryptionKeyID]
		if !ok {
			return nil, nil, errors.ThrowInternalf(nil, "CRYPT-v2Kas", "encryption key not found")
		}
		keys[config.EncryptionKeyID] = key
		ids = append(ids, config.EncryptionKeyID)
	}
	for _, id := range config.DecryptionKeyIDs {
		key, ok := readKeys[id]
		if !ok {
			logging.Log("CRYPT-s23rf").Warnf("description key %s not found", id)
			continue
		}
		keys[id] = key
		ids = append(ids, id)
	}
	return keys, ids, nil
}
