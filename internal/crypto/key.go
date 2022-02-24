package crypto

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
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

type EncryptionKeys struct {
	DomainVerification   *KeyConfig
	IDPConfig            *KeyConfig
	OIDC                 *KeyConfig
	OTP                  *KeyConfig
	SMS                  *KeyConfig
	SMTP                 *KeyConfig
	User                 *KeyConfig
	CSRFCookieKeyID      string
	UserAgentCookieKeyID string
}

func LoadKey(keyStorage KeyStorage, id string) (string, error) {
	key, err := keyStorage.ReadKey(id)
	if err != nil {
		return "", err
	}
	return key.Value, nil
}

func LoadKeys(config *KeyConfig, keyStorage KeyStorage) (map[string]string, []string, error) {
	if config == nil {
		return nil, nil, errors.ThrowInvalidArgument(nil, "CRYPT-dJK8s", "config must not be nil")
	}
	readKeys, err := keyStorage.ReadKeys()
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
			logging.Errorf("description key %s not found", id)
			continue
		}
		keys[id] = key
		ids = append(ids, id)
	}
	return keys, ids, nil
}
