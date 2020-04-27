package systemdefaults

import "github.com/caos/zitadel/internal/crypto"

type SystemDefaults struct {
	SecretGenerators    SecretGenerators
	UserVerificationKey *crypto.KeyConfig
	Multifactors        Multifactors
}

type SecretGenerators struct {
	PasswordSaltCost         int
	ClientSecretGenerator    crypto.GeneratorConfig
	InitializeUserCode       crypto.GeneratorConfig
	EmailVerificationCode    crypto.GeneratorConfig
	PhoneVerificationCode    crypto.GeneratorConfig
	PasswordVerificationCode crypto.GeneratorConfig
}

type Multifactors struct {
	OTP OTP
}

type OTP struct {
	Issuer          string
	VerificationKey *crypto.KeyConfig
	CryptoMFA       crypto.EncryptionAlgorithm
}
