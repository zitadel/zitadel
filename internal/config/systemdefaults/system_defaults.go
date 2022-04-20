package systemdefaults

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/crypto"
)

type SystemDefaults struct {
	DefaultLanguage    language.Tag
	Domain             string
	ZitadelDocs        ZitadelDocs
	SecretGenerators   SecretGenerators
	Multifactors       MultifactorConfig
	DomainVerification DomainVerification
	Notifications      Notifications
	KeyConfig          KeyConfig
}

type ZitadelDocs struct {
	Issuer            string
	DiscoveryEndpoint string
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
