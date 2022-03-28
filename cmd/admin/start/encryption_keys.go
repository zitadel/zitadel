package start

import (
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

var (
	defaultKeyIDs = []string{
		"domainVerificationKey",
		"idpConfigKey",
		"oidcKey",
		"otpKey",
		"smsKey",
		"smtpKey",
		"userKey",
		"csrfCookieKey",
		"userAgentCookieKey",
	}
)

type encryptionKeys struct {
	DomainVerification crypto.EncryptionAlgorithm
	IDPConfig          crypto.EncryptionAlgorithm
	OIDC               crypto.EncryptionAlgorithm
	OTP                crypto.EncryptionAlgorithm
	SMS                crypto.EncryptionAlgorithm
	SMTP               crypto.EncryptionAlgorithm
	User               crypto.EncryptionAlgorithm
	CSRFCookieKey      []byte
	UserAgentCookieKey []byte
	OIDCKey            []byte
}

func ensureEncryptionKeys(keyConfig *encryptionKeyConfig, keyStorage crypto.KeyStorage) (*encryptionKeys, error) {
	keys, err := keyStorage.ReadKeys()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		if err := createDefaultKeys(keyStorage); err != nil {
			return nil, err
		}
	}
	encryptionKeys := new(encryptionKeys)
	encryptionKeys.DomainVerification, err = crypto.NewAESCrypto(keyConfig.DomainVerification, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.IDPConfig, err = crypto.NewAESCrypto(keyConfig.IDPConfig, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.OIDC, err = crypto.NewAESCrypto(keyConfig.OIDC, keyStorage)
	if err != nil {
		return nil, err
	}
	key, err := crypto.LoadKey(keyConfig.OIDC.EncryptionKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.OIDCKey = []byte(key)
	encryptionKeys.OTP, err = crypto.NewAESCrypto(keyConfig.OTP, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.SMS, err = crypto.NewAESCrypto(keyConfig.SMS, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.SMTP, err = crypto.NewAESCrypto(keyConfig.SMTP, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.User, err = crypto.NewAESCrypto(keyConfig.User, keyStorage)
	if err != nil {
		return nil, err
	}
	key, err = crypto.LoadKey(keyConfig.CSRFCookieKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.CSRFCookieKey = []byte(key)
	key, err = crypto.LoadKey(keyConfig.UserAgentCookieKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	encryptionKeys.UserAgentCookieKey = []byte(key)
	return encryptionKeys, nil
}

func createDefaultKeys(keyStorage crypto.KeyStorage) error {
	keys := make([]*crypto.Key, len(defaultKeyIDs))
	for i, keyID := range defaultKeyIDs {
		key, err := crypto.NewKey(keyID)
		if err != nil {
			return err
		}
		keys[i] = key
	}
	if err := keyStorage.CreateKeys(keys...); err != nil {
		return caos_errs.ThrowInternal(err, "START-aGBq2", "cannot create default keys")
	}
	return nil
}
