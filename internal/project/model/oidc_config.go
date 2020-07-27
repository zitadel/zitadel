package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"strings"
)

const (
	http          = "http://"
	httpLocalhost = "http://localhost"
	https         = "https://"
)

type OIDCConfig struct {
	es_models.ObjectRoot
	AppID                  string
	ClientID               string
	ClientSecret           *crypto.CryptoValue
	ClientSecretString     string
	RedirectUris           []string
	ResponseTypes          []OIDCResponseType
	GrantTypes             []OIDCGrantType
	ApplicationType        OIDCApplicationType
	AuthMethodType         OIDCAuthMethodType
	PostLogoutRedirectUris []string
	OIDCVersion            OIDCVersion
	Compliance             *Compliance
}

type OIDCVersion int32

const (
	OIDCVersionV1 OIDCVersion = iota
)

type OIDCResponseType int32

const (
	OIDCResponseTypeCode OIDCResponseType = iota
	OIDCResponseTypeIDToken
	OIDCResponseTypeIDTokenToken
)

type OIDCGrantType int32

const (
	OIDCGrantTypeAuthorizationCode OIDCGrantType = iota
	OIDCGrantTypeImplicit
	OIDCGrantTypeRefreshToken
)

type OIDCApplicationType int32

const (
	OIDCApplicationTypeWeb OIDCApplicationType = iota
	OIDCApplicationTypeUserAgent
	OIDCApplicationTypeNative
)

type OIDCAuthMethodType int32

const (
	OIDCAuthMethodTypeBasic OIDCAuthMethodType = iota
	OIDCAuthMethodTypePost
	OIDCAuthMethodTypeNone
)

type Compliance struct {
	NoneCompliant bool
	Problems      []string
}

func (c *OIDCConfig) IsValid() bool {
	grantTypes := c.getRequiredGrantTypes()
	for _, grantType := range grantTypes {
		ok := containsOIDCGrantType(c.GrantTypes, grantType)
		if !ok {
			return false
		}
	}
	return true
}

func (c *OIDCConfig) FillCompliance() {
	c.Compliance = GetOIDCCompliance(c.OIDCVersion, c.ApplicationType, c.GrantTypes, c.ResponseTypes, c.AuthMethodType, c.RedirectUris)
}

func GetOIDCCompliance(version OIDCVersion, appType OIDCApplicationType, grantTypes []OIDCGrantType, responseTypes []OIDCResponseType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	switch version {
	case OIDCVersionV1:
		return GetOIDCV1Compliance(appType, grantTypes, authMethod, redirectUris)
	}
	return nil
}

func GetOIDCV1Compliance(appType OIDCApplicationType, grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	switch appType {
	case OIDCApplicationTypeNative:
		return GetOIDCNativeApplicationCompliance(grantTypes, authMethod, redirectUris)
	case OIDCApplicationTypeWeb:
		return GetOIDCWebApplicationCompliance(grantTypes, redirectUris)
	case OIDCApplicationTypeUserAgent:
		return GetOIDCUserAgentApplicationCompliance(grantTypes, authMethod, redirectUris)
	}
	return nil
}

func GetOIDCNativeApplicationCompliance(grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if len(grantTypes) != 1 {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Native.GrantType.MultipleTypes")
	} else if !containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Native.GrantType.NotAuthorizationCodeFlow")
	}
	if authMethod != OIDCAuthMethodTypeNone {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Native.AuthMethodType.NotNone")
	}
	if !onlyLocalhostIsHttp(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Native.RediredtUris.HttpOnlyForLocalhost")
	}
	return compliance
}

func GetOIDCWebApplicationCompliance(grantTypes []OIDCGrantType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if len(grantTypes) != 1 {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Web.GrantType.NotAuthorizationCodeFlow")
	} else if !containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Web.GrantType.NotAuthorizationCodeFlow")
	}
	if !urlsAreHttps(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Web.RediredtUris.NotHttps")
	}
	return compliance
}

func GetOIDCUserAgentApplicationCompliance(grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		if authMethod != OIDCAuthMethodTypeNone {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.UserAgent.AuthorizationCodeFlow.AuthMethodType.NotNone")
		}
	}
	if containsOIDCGrantType(grantTypes, OIDCGrantTypeImplicit) {
		if authMethod != OIDCAuthMethodTypePost {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.UserAgent.Implicit.AuthMethodType.NotPost")
		}
	}
	if !urlsAreHttps(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.Web.RediredtUris.NotHttps")
	}
	return compliance
}

func (c *OIDCConfig) getRequiredGrantTypes() []OIDCGrantType {
	grantTypes := make([]OIDCGrantType, 0)
	implicit := false
	for _, r := range c.ResponseTypes {
		switch r {
		case OIDCResponseTypeCode:
			grantTypes = append(grantTypes, OIDCGrantTypeAuthorizationCode)
		case OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken:
			if !implicit {
				implicit = true
				grantTypes = append(grantTypes, OIDCGrantTypeImplicit)
			}
		}
	}
	return grantTypes
}

func containsOIDCGrantType(grantTypes []OIDCGrantType, grantType OIDCGrantType) bool {
	for _, gt := range grantTypes {
		if gt == grantType {
			return true
		}
	}
	return false
}

func urlsAreHttps(uris []string) bool {
	for _, uri := range uris {
		if !strings.HasPrefix(uri, https) {
			return false
		}
	}
	return true
}

func onlyLocalhostIsHttp(uris []string) bool {
	for _, uri := range uris {
		if strings.HasPrefix(uri, http) {
			if !strings.HasPrefix(uri, httpLocalhost) {
				return false
			}
		}
	}
	return true
}
