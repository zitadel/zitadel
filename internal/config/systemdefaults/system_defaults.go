package systemdefaults

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/providers/chat"
	"github.com/caos/zitadel/internal/notification/providers/email"
	"github.com/caos/zitadel/internal/notification/providers/twilio"
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
	IamID                    string
	Notifications            Notifications
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
	MachineKeySize           uint32
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
	Providers    Providers
	TemplateData TemplateData
}

type Endpoints struct {
	InitCode      string
	PasswordReset string
	VerifyEmail   string
	DomainClaimed string
}

type Providers struct {
	Chat   chat.ChatConfig
	Email  email.EmailConfig
	Twilio twilio.TwilioConfig
}

type TemplateData struct {
	InitCode      templates.TemplateData
	PasswordReset templates.TemplateData
	VerifyEmail   templates.TemplateData
	VerifyPhone   templates.TemplateData
	DomainClaimed templates.TemplateData
}
