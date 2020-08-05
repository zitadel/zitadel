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
		return GetOIDCV1NativeApplicationCompliance(grantTypes, authMethod, redirectUris)
	case OIDCApplicationTypeWeb:
		return GetOIDCV1WebApplicationCompliance(grantTypes, redirectUris)
	case OIDCApplicationTypeUserAgent:
		return GetOIDCV1UserAgentApplicationCompliance(grantTypes, authMethod, redirectUris)
	}
	return nil
}

func GetOIDCV1NativeApplicationCompliance(grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if len(grantTypes) != 1 {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.GrantType.MultipleTypes")
	} else if !containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.GrantType.NotAuthorizationCodeFlow")
	}
	if authMethod != OIDCAuthMethodTypeNone {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.AuthMethodType.NotNone")
	}
	if !onlyLocalhostIsHttp(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.RediredtUris.HttpOnlyForLocalhost")
	}
	if hasHttpsUrl(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.RediredtUris.HttpsNotAllowed")
	}
	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
	}
	return compliance
}

func GetOIDCV1WebApplicationCompliance(grantTypes []OIDCGrantType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if len(grantTypes) != 1 {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Web.GrantType.MultipleTypes")
	} else if !containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Web.GrantType.NotAuthorizationCodeFlow")
	}
	if !urlsAreHttps(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Web.RediredtUris.NotHttps")
	}
	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
	}
	return compliance
}

func GetOIDCV1UserAgentApplicationCompliance(grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		if authMethod != OIDCAuthMethodTypeNone {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.UserAgent.AuthorizationCodeFlow.AuthMethodType.NotNone")
		}
	}
	if containsOIDCGrantType(grantTypes, OIDCGrantTypeImplicit) {
		if authMethod != OIDCAuthMethodTypeNone {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.UserAgent.Implicit.AuthMethodType.NotNone")
		}
	}
	if !urlsAreHttps(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.UserAgent.RediredtUris.NotHttps")
	}
	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
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

func hasHttpsUrl(uris []string) bool {
	for _, uri := range uris {
		if strings.HasPrefix(uri, https) {
			return true
		}
	}
	return false
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
