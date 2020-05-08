package systemdefaults

import (
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
)

type SystemDefaults struct {
	SecretGenerators      SecretGenerators
	UserVerificationKey   *crypto.KeyConfig
	Multifactors          MultifactorConfig
	VerificationLifetimes VerificationLifetimes
}

type SecretGenerators struct {
	PasswordSaltCost         int
	ClientSecretGenerator    crypto.GeneratorConfig
	InitializeUserCode       crypto.GeneratorConfig
	EmailVerificationCode    crypto.GeneratorConfig
	PhoneVerificationCode    crypto.GeneratorConfig
	PasswordVerificationCode crypto.GeneratorConfig
}

type MultifactorConfig struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer          string
	VerificationKey *crypto.KeyConfig
}

type VerificationLifetimes struct {
	PasswordCheck    types.Duration
	MfaInitSkip      types.Duration
	MfaSoftwareCheck types.Duration
	MfaHardwareCheck types.Duration
}
