package start

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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

type encryptionKeys struct {
	DomainVerification crypto.EncryptionAlgorithm
	IDPConfig          crypto.EncryptionAlgorithm
	OIDC               crypto.EncryptionAlgorithm
	SAML               crypto.EncryptionAlgorithm
	OTP                crypto.EncryptionAlgorithm
	SMS                crypto.EncryptionAlgorithm
	SMTP               crypto.EncryptionAlgorithm
	User               crypto.EncryptionAlgorithm
	CSRFCookieKey      []byte
	UserAgentCookieKey []byte
	OIDCKey            []byte
}

func ensureEncryptionKeys(ctx context.Context, keyConfig *encryptionKeyConfig, keyStorage crypto.KeyStorage) (keys *encryptionKeys, err error) {
	if err := verifyDefaultKeys(ctx, keyStorage); err != nil {
		return nil, err
	}
	keys = new(encryptionKeys)
	keys.DomainVerification, err = crypto.NewAESCrypto(keyConfig.DomainVerification, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.IDPConfig, err = crypto.NewAESCrypto(keyConfig.IDPConfig, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.OIDC, err = crypto.NewAESCrypto(keyConfig.OIDC, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.SAML, err = crypto.NewAESCrypto(keyConfig.SAML, keyStorage)
	if err != nil {
		return nil, err
	}
	key, err := crypto.LoadKey(keyConfig.OIDC.EncryptionKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.OIDCKey = []byte(key)
	keys.OTP, err = crypto.NewAESCrypto(keyConfig.OTP, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.SMS, err = crypto.NewAESCrypto(keyConfig.SMS, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.SMTP, err = crypto.NewAESCrypto(keyConfig.SMTP, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.User, err = crypto.NewAESCrypto(keyConfig.User, keyStorage)
	if err != nil {
		return nil, err
	}
	key, err = crypto.LoadKey(keyConfig.CSRFCookieKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.CSRFCookieKey = []byte(key)
	key, err = crypto.LoadKey(keyConfig.UserAgentCookieKeyID, keyStorage)
	if err != nil {
		return nil, err
	}
	keys.UserAgentCookieKey = []byte(key)
	return keys, nil
}

func verifyDefaultKeys(ctx context.Context, keyStorage crypto.KeyStorage) (err error) {
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
	if err := keyStorage.CreateKeys(ctx, keys...); err != nil {
		return zerrors.ThrowInternal(err, "START-aGBq2", "cannot create default keys")
	}
	return nil
}
