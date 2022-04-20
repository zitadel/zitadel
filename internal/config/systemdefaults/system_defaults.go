package systemdefaults

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
)

type SystemDefaults struct {
	Domain             string //TODO: remove with v2
	SecretGenerators   SecretGenerators
	Multifactors       MultifactorConfig
	DomainVerification DomainVerification
	Notifications      Notifications
	KeyConfig          KeyConfig
}

type SecretGenerators struct {
	PasswordSaltCost   int
	MachineKeySize     uint32
	ApplicationKeySize uint32
}

type MultifactorConfig struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer string
}

type DomainVerification struct {
	VerificationGenerator crypto.GeneratorConfig
}

type Notifications struct {
	Endpoints      Endpoints
	FileSystemPath string
}

type Endpoints struct {
	InitCode                 string
	PasswordReset            string
	VerifyEmail              string
	DomainClaimed            string
	PasswordlessRegistration string
}

type KeyConfig struct {
	Size               int
	PrivateKeyLifetime time.Duration
	PublicKeyLifetime  time.Duration
}
