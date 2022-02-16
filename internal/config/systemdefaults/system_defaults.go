package systemdefaults

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/channels/chat"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/notification/templates"
)

type SystemDefaults struct {
	DefaultLanguage             language.Tag
	Domain                      string
	ZitadelDocs                 ZitadelDocs
	SecretGenerators            SecretGenerators
	UserVerificationKey         *crypto.KeyConfig
	IDPConfigVerificationKey    *crypto.KeyConfig
	SMTPPasswordVerificationKey *crypto.KeyConfig
	Multifactors                MultifactorConfig
	VerificationLifetimes       VerificationLifetimes
	DomainVerification          DomainVerification
	Notifications               Notifications
	KeyConfig                   KeyConfig
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
	Issuer          string
	VerificationKey *crypto.KeyConfig
}

type VerificationLifetimes struct {
	PasswordCheck      time.Duration
	ExternalLoginCheck time.Duration
	MFAInitSkip        time.Duration
	SecondFactorCheck  time.Duration
	MultiFactorCheck   time.Duration
}

type DomainVerification struct {
	VerificationKey       *crypto.KeyConfig
	VerificationGenerator crypto.GeneratorConfig
}

type Notifications struct {
	DebugMode bool
	Endpoints Endpoints
	Providers Channels
}

type Endpoints struct {
	InitCode                 string
	PasswordReset            string
	VerifyEmail              string
	DomainClaimed            string
	PasswordlessRegistration string
}

type Channels struct {
	Chat       chat.ChatConfig
	Twilio     twilio.TwilioConfig
	FileSystem fs.FSConfig
	Log        log.LogConfig
}

type TemplateData struct {
	InitCode      templates.TemplateData
	PasswordReset templates.TemplateData
	VerifyEmail   templates.TemplateData
	VerifyPhone   templates.TemplateData
	DomainClaimed templates.TemplateData
}

type KeyConfig struct {
	Size                     int
	PrivateKeyLifetime       time.Duration
	PublicKeyLifetime        time.Duration
	SigningKeyRotationCheck  time.Duration
	SigningKeyGracefulPeriod time.Duration
}
