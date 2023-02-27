package domain

import (
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type LoginPolicy struct {
	models.ObjectRoot

	Default                    bool
	AllowUsernamePassword      bool
	AllowRegister              bool
	AllowExternalIDP           bool
	IDPProviders               []*IDPProvider
	ForceMFA                   bool
	SecondFactors              []SecondFactorType
	MultiFactors               []MultiFactorType
	PasswordlessType           PasswordlessType
	HidePasswordReset          bool
	IgnoreUnknownUsernames     bool
	AllowDomainDiscovery       bool
	DefaultRedirectURI         string
	PasswordCheckLifetime      time.Duration
	ExternalLoginCheckLifetime time.Duration
	MFAInitSkipLifetime        time.Duration
	SecondFactorCheckLifetime  time.Duration
	MultiFactorCheckLifetime   time.Duration
	DisableLoginWithEmail      bool
	DisableLoginWithPhone      bool
}

func ValidateDefaultRedirectURI(rawURL string) bool {
	if rawURL == "" {
		return true
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	switch parsedURL.Scheme {
	case "":
		return false
	case "http", "https":
		return parsedURL.Host != ""
	default:
		return true
	}
}

type IDPProvider struct {
	models.ObjectRoot
	Type        IdentityProviderType
	IDPConfigID string

	Name        string
	StylingType IDPConfigStylingType // deprecated
	IDPType     IDPType
	IDPState    IDPConfigState
}

func (p IDPProvider) IsValid() bool {
	return p.IDPConfigID != ""
}

type PasswordlessType int32

const (
	PasswordlessTypeNotAllowed PasswordlessType = iota
	PasswordlessTypeAllowed

	passwordlessCount
)

func (f PasswordlessType) Valid() bool {
	return f >= 0 && f < passwordlessCount
}

func (p *LoginPolicy) HasSecondFactors() bool {
	return len(p.SecondFactors) > 0
}

func (p *LoginPolicy) HasMultiFactors() bool {
	return len(p.MultiFactors) > 0
}
