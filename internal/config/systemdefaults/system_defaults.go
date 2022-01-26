package systemdefaults

import (
	"github.com/caos/zitadel/internal/notification/channels/log"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/channels/chat"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/notification/templates"
)

type SystemDefaults struct {
	DefaultLanguage          language.Tag
	Domain                   string
	ZitadelDocs              ZitadelDocs
	SecretGenerators         SecretGenerators
	UserVerificationKey      *crypto.KeyConfig
	IDPConfigVerificationKey *crypto.KeyConfig
	Multifactors             MultifactorConfig
	VerificationLifetimes    VerificationLifetimes
	DomainVerification       DomainVerification
	Notifications            Notifications
	WebAuthN                 WebAuthN
	KeyConfig                KeyConfig
}

type ZitadelDocs struct {
	Issuer            string
	DiscoveryEndpoint string
}

type SecretGenerators struct {
	PasswordSaltCost         int
	ClientSecretGenerator    crypto.GeneratorConfig
	InitializeUserCode       crypto.GeneratorConfig
	EmailVerificationCode    crypto.GeneratorConfig
	PhoneVerificationCode    crypto.GeneratorConfig
	PasswordVerificationCode crypto.GeneratorConfig
	PasswordlessInitCode     crypto.GeneratorConfig
	MachineKeySize           uint32
	ApplicationKeySize       uint32
}

type MultifactorConfig struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer          string
	VerificationKey *crypto.KeyConfig
}

type VerificationLifetimes struct {
	PasswordCheck      types.Duration
	ExternalLoginCheck types.Duration
	MFAInitSkip        types.Duration
	SecondFactorCheck  types.Duration
	MultiFactorCheck   types.Duration
}

type DomainVerification struct {
	VerificationKey       *crypto.KeyConfig
	VerificationGenerator crypto.GeneratorConfig
}

type Notifications struct {
	DebugMode    bool
	Endpoints    Endpoints
	Providers    Channels
	TemplateData TemplateData
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
	Email      smtp.EmailConfig
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

type WebAuthN struct {
	ID            string
	OriginLogin   string
	OriginConsole string
	DisplayName   string
}

type KeyConfig struct {
	Size                     int
	PrivateKeyLifetime       types.Duration
	PublicKeyLifetime        types.Duration
	EncryptionConfig         *crypto.KeyConfig
	SigningKeyRotationCheck  types.Duration
	SigningKeyGracefulPeriod types.Duration
}
