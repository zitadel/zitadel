package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type OTP struct {
	es_models.ObjectRoot

	Secret       *crypto.CryptoValue
	SecretString string
	Url          string
	State        MFAState
}

type OTPState int32

const (
	OTPStateUnspecified OTPState = iota
	OTPStateActive
	OTPStateRemoved

	otpStateCount
)

func (s OTPState) Valid() bool {
	return s >= 0 && s < otpStateCount
}

func NewOTPKey(issuer, accountName string, cryptoAlg crypto.EncryptionAlgorithm) (*otp.Key, *crypto.CryptoValue, error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: issuer, AccountName: accountName})
	if err != nil {
		return nil, nil, err
	}
	encryptedSecret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	if err != nil {
		return nil, nil, err
	}
	return key, encryptedSecret, nil
}
