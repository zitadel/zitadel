package systemdefaults

import (
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/providers/chat"
	"github.com/caos/zitadel/internal/notification/providers/email"
	"github.com/caos/zitadel/internal/notification/providers/twilio"
	"github.com/caos/zitadel/internal/notification/templates"
	org_model "github.com/caos/zitadel/internal/org/model"
	pol "github.com/caos/zitadel/internal/policy"
	"golang.org/x/text/language"
)

type SystemDefaults struct {
	DefaultLanguage       language.Tag
	SecretGenerators      SecretGenerators
	UserVerificationKey   *crypto.KeyConfig
	Multifactors          MultifactorConfig
	VerificationLifetimes VerificationLifetimes
	DefaultPolicies       DefaultPolicies
	IamID                 string
	SetUp                 types.IAMSetUp
	Notifications         Notifications
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

type DefaultPolicies struct {
	Age        pol.PasswordAgePolicyDefault
	Complexity pol.PasswordComplexityPolicyDefault
	Lockout    pol.PasswordLockoutPolicyDefault
	OrgIam     org_model.OrgIamPolicy
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
}
