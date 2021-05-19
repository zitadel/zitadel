package domain

import (
	"strings"
	"time"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	http                          = "http://"
	httpLocalhostWithPort         = "http://localhost:"
	httpLocalhostWithoutPort      = "http://localhost/"
	httpLoopbackV4WithPort        = "http://127.0.0.1:"
	httpLoopbackV4WithoutPort     = "http://127.0.0.1/"
	httpLoopbackV6WithPort        = "http://[::1]:"
	httpLoopbackV6WithoutPort     = "http://[::1]/"
	httpLoopbackV6LongWithPort    = "http://[0:0:0:0:0:0:0:1]:"
	httpLoopbackV6LongWithoutPort = "http://[0:0:0:0:0:0:0:1]/"
	https                         = "https://"
)

type OIDCApp struct {
	models.ObjectRoot

	AppID                    string
	AppName                  string
	ClientID                 string
	ClientSecret             *crypto.CryptoValue
	ClientSecretString       string
	RedirectUris             []string
	ResponseTypes            []OIDCResponseType
	GrantTypes               []OIDCGrantType
	ApplicationType          OIDCApplicationType
	AuthMethodType           OIDCAuthMethodType
	PostLogoutRedirectUris   []string
	OIDCVersion              OIDCVersion
	Compliance               *Compliance
	DevMode                  bool
	AccessTokenType          OIDCTokenType
	AccessTokenRoleAssertion bool
	IDTokenRoleAssertion     bool
	IDTokenUserinfoAssertion bool
	ClockSkew                time.Duration
	AdditionalOrigins        []string

	State AppState
}

func (a *OIDCApp) GetApplicationName() string {
	return a.AppName
}

func (a *OIDCApp) GetState() AppState {
	return a.State
}

func (a *OIDCApp) setClientID(clientID string) {
	a.ClientID = clientID
}

func (a *OIDCApp) setClientSecret(clientSecret *crypto.CryptoValue) {
	a.ClientSecret = clientSecret
}

func (a *OIDCApp) requiresClientSecret() bool {
	return a.AuthMethodType == OIDCAuthMethodTypeBasic || a.AuthMethodType == OIDCAuthMethodTypePost
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
	OIDCAuthMethodTypePrivateKeyJWT
)

type Compliance struct {
	NoneCompliant bool
	Problems      []string
}

type OIDCTokenType int32

const (
	OIDCTokenTypeBearer OIDCTokenType = iota
	OIDCTokenTypeJWT
)

func (a *OIDCApp) IsValid() bool {
	if a.ClockSkew > time.Second*5 || a.ClockSkew < time.Second*0 || !a.OriginsValid() {
		return false
	}
	grantTypes := a.getRequiredGrantTypes()
	if len(grantTypes) == 0 {
		return false
	}
	for _, grantType := range grantTypes {
		ok := containsOIDCGrantType(a.GrantTypes, grantType)
		if !ok {
			return false
		}
	}
	return true
}

func (a *OIDCApp) OriginsValid() bool {
	for _, origin := range a.AdditionalOrigins {
		if !http_util.IsOrigin(origin) {
			return false
		}
	}
	return true
}

func (a *OIDCApp) getRequiredGrantTypes() []OIDCGrantType {
	grantTypes := make([]OIDCGrantType, 0)
	implicit := false
	for _, r := range a.ResponseTypes {
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

func (a *OIDCApp) FillCompliance() {
	a.Compliance = GetOIDCCompliance(a.OIDCVersion, a.ApplicationType, a.GrantTypes, a.ResponseTypes, a.AuthMethodType, a.RedirectUris)
}

func GetOIDCCompliance(version OIDCVersion, appType OIDCApplicationType, grantTypes []OIDCGrantType, responseTypes []OIDCResponseType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	switch version {
	case OIDCVersionV1:
		return GetOIDCV1Compliance(appType, grantTypes, authMethod, redirectUris)
	}
	return nil
}

func GetOIDCV1Compliance(appType OIDCApplicationType, grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}
	if redirectUris == nil || len(redirectUris) == 0 {
		compliance.NoneCompliant = true
		compliance.Problems = append([]string{"Application.OIDC.V1.NoRedirectUris"}, compliance.Problems...)
	}
	if containsOIDCGrantType(grantTypes, OIDCGrantTypeImplicit) && containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		CheckRedirectUrisImplicitAndCode(compliance, appType, redirectUris)
	} else {
		if containsOIDCGrantType(grantTypes, OIDCGrantTypeImplicit) {
			CheckRedirectUrisImplicit(compliance, appType, redirectUris)
		}
		if containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
			CheckRedirectUrisCode(compliance, appType, redirectUris)
		}
	}

	switch appType {
	case OIDCApplicationTypeNative:
		GetOIDCV1NativeApplicationCompliance(compliance, authMethod)
	case OIDCApplicationTypeUserAgent:
		GetOIDCV1UserAgentApplicationCompliance(compliance, authMethod)
	}
	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
	}
	return compliance
}

func GetOIDCV1NativeApplicationCompliance(compliance *Compliance, authMethod OIDCAuthMethodType) {
	if authMethod != OIDCAuthMethodTypeNone {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.AuthMethodType.NotNone")
	}
}

func GetOIDCV1UserAgentApplicationCompliance(compliance *Compliance, authMethod OIDCAuthMethodType) {
	if authMethod != OIDCAuthMethodTypeNone {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.UserAgent.AuthMethodType.NotNone")
	}
}

func CheckRedirectUrisCode(compliance *Compliance, appType OIDCApplicationType, redirectUris []string) {
	if urlsAreHttps(redirectUris) {
		return
	}
	if urlContainsPrefix(redirectUris, http) {
		if appType == OIDCApplicationTypeUserAgent {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb")
		}
		if appType == OIDCApplicationTypeNative && !onlyLocalhostIsHttp(redirectUris) {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Code.RedirectUris.NativeShouldBeHttpLocalhost")
		}
	}
	if containsCustom(redirectUris) && appType != OIDCApplicationTypeNative {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Code.RedirectUris.CustomOnlyForNative")
	}
}

func CheckRedirectUrisImplicit(compliance *Compliance, appType OIDCApplicationType, redirectUris []string) {
	if urlsAreHttps(redirectUris) {
		return
	}
	if containsCustom(redirectUris) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed")
	}
	if urlContainsPrefix(redirectUris, http) {
		if appType == OIDCApplicationTypeNative {
			if !onlyLocalhostIsHttp(redirectUris) {
				compliance.NoneCompliant = true
				compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Implicit.RedirectUris.NativeShouldBeHttpLocalhost")
			}
			return
		}
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed")
	}
}

func CheckRedirectUrisImplicitAndCode(compliance *Compliance, appType OIDCApplicationType, redirectUris []string) {
	if urlsAreHttps(redirectUris) {
		return
	}
	if containsCustom(redirectUris) && appType != OIDCApplicationTypeNative {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed")
	}
	if urlContainsPrefix(redirectUris, http) {
		if appType == OIDCApplicationTypeUserAgent {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb")
		}
		if !onlyLocalhostIsHttp(redirectUris) && appType == OIDCApplicationTypeNative {
			compliance.NoneCompliant = true
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Implicit.RedirectUris.NativeShouldBeHttpLocalhost")
		}
	}
	if !compliance.NoneCompliant {
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.NotAllCombinationsAreAllowed")
	}
}

func urlsAreHttps(uris []string) bool {
	for _, uri := range uris {
		if !strings.HasPrefix(uri, https) {
			return false
		}
	}
	return true
}

func urlContainsPrefix(uris []string, prefix string) bool {
	for _, uri := range uris {
		if strings.HasPrefix(uri, prefix) {
			return true
		}
	}
	return false
}

func containsCustom(uris []string) bool {
	for _, uri := range uris {
		if !strings.HasPrefix(uri, http) && !strings.HasPrefix(uri, https) {
			return true
		}
	}
	return false
}

func onlyLocalhostIsHttp(uris []string) bool {
	for _, uri := range uris {
		if strings.HasPrefix(uri, http) && !isHTTPLoopbackLocalhost(uri) {
			return false
		}
	}
	return true
}

func isHTTPLoopbackLocalhost(uri string) bool {
	return strings.HasPrefix(uri, httpLocalhostWithoutPort) ||
		strings.HasPrefix(uri, httpLocalhostWithPort) ||
		strings.HasPrefix(uri, httpLoopbackV4WithoutPort) ||
		strings.HasPrefix(uri, httpLoopbackV4WithPort) ||
		strings.HasPrefix(uri, httpLoopbackV6WithoutPort) ||
		strings.HasPrefix(uri, httpLoopbackV6WithPort) ||
		strings.HasPrefix(uri, httpLoopbackV6LongWithoutPort) ||
		strings.HasPrefix(uri, httpLoopbackV6LongWithPort)
}
