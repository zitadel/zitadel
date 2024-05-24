package domain

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type TOTP struct {
	*ObjectDetails

	Secret string
	URI    string
}

func NewTOTPKey(issuer, accountName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: issuer, AccountName: accountName})
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TOTP-ieY3o", "Errors.Internal")
	}
	return key, nil
}

func VerifyTOTP(code string, secret *crypto.CryptoValue, cryptoAlg crypto.EncryptionAlgorithm) error {
	decrypt, err := crypto.DecryptString(secret, cryptoAlg)
	if err != nil {
		return err
	}

	valid := totp.Validate(code, decrypt)
	if !valid {
		return zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode")
	}
	return nil
}
