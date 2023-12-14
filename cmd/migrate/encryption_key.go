package migrate

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type encryptionKeyConfig struct {
	OIDC *crypto.KeyConfig
	SAML *crypto.KeyConfig
}

type encryptionKeys struct {
	OIDC crypto.EncryptionAlgorithm
	SAML crypto.EncryptionAlgorithm
}

func ensureEncryptionKeys(keyConfig *encryptionKeyConfig, keyStorage crypto.KeyStorage) (keys *encryptionKeys, err error) {
	if err := verifyDefaultKeys(keyStorage); err != nil {
		return nil, err
	}
	keys = new(encryptionKeys)
	keys.OIDC, err = crypto.NewAESCrypto(keyConfig.OIDC, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.SAML, err = crypto.NewAESCrypto(keyConfig.SAML, keyStorage)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

var (
	defaultKeyIDs = []string{
		"domainVerificationKey",
		"idpConfigKey",
		"oidcKey",
		"samlKey",
		"otpKey",
		"smsKey",
		"smtpKey",
		"userKey",
		"csrfCookieKey",
		"userAgentCookieKey",
	}
)

func verifyDefaultKeys(keyStorage crypto.KeyStorage) (err error) {
	keys := make([]*crypto.Key, 0, len(defaultKeyIDs))
	for _, keyID := range defaultKeyIDs {
		_, err := crypto.LoadKey(keyID, keyStorage)
		if err == nil {
			continue
		}
		key, err := crypto.NewKey(keyID)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := keyStorage.CreateKeys(keys...); err != nil {
		return zerrors.ThrowInternal(err, "MIGRA-aGBq2", "cannot create default keys")
	}
	return nil
}
