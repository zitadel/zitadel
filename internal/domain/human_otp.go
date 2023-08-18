package domain

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

type TOTP struct {
	*ObjectDetails

	Secret string
	URI    string
}

func NewTOTPKey(issuer, accountName string, cryptoAlg crypto.EncryptionAlgorithm) (*otp.Key, *crypto.CryptoValue, error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: issuer, AccountName: accountName})
	if err != nil {
		return nil, nil, caos_errs.ThrowInternal(err, "TOTP-ieY3o", "Errors.Internal")
	}
	encryptedSecret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	if err != nil {
		return nil, nil, err
	}
	return key, encryptedSecret, nil
}

func VerifyTOTP(code string, secret *crypto.CryptoValue, cryptoAlg crypto.EncryptionAlgorithm) error {
	decrypt, err := crypto.DecryptString(secret, cryptoAlg)
	if err != nil {
		return err
	}

	valid := totp.Validate(code, decrypt)
	if !valid {
		return caos_errs.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode")
	}
	return nil
}
