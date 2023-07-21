package systemdefaults

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
)

type SystemDefaults struct {
	SecretGenerators   SecretGenerators
	PasswordHasher     crypto.PasswordHashConfig
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
	FileSystemPath string
}

type KeyConfig struct {
	Size                int
	PrivateKeyLifetime  time.Duration
	PublicKeyLifetime   time.Duration
	CertificateSize     int
	CertificateLifetime time.Duration
}
