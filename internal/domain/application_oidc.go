package domain

import (
	"strings"
	"time"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
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
	EncodedHash              string
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
	SkipNativeAppSuccessPage bool
	BackChannelLogoutURI     string
	LoginVersion             LoginVersion
	LoginBaseURI             string

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

func (a *OIDCApp) setClientSecret(encodedHash string) {
	a.EncodedHash = encodedHash
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
	OIDCResponseTypeUnspecified OIDCResponseType = iota - 1 // Negative offset not to break existing configs.
	OIDCResponseTypeCode
	OIDCResponseTypeIDToken
	OIDCResponseTypeIDTokenToken
)

//go:generate enumer -type OIDCResponseMode -transform snake -trimprefix OIDCResponseMode
type OIDCResponseMode int

const (
	OIDCResponseModeUnspecified OIDCResponseMode = iota
	OIDCResponseModeQuery
	OIDCResponseModeFragment
	OIDCResponseModeFormPost
)

type OIDCGrantType int32

const (
	OIDCGrantTypeAuthorizationCode OIDCGrantType = iota
	OIDCGrantTypeImplicit
	OIDCGrantTypeRefreshToken
	OIDCGrantTypeDeviceCode
	OIDCGrantTypeTokenExchange
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
		if !http_util.IsOrigin(strings.TrimSpace(origin)) {
			return false
		}
	}
	return true
}

func ContainsRequiredGrantTypes(responseTypes []OIDCResponseType, grantTypes []OIDCGrantType) bool {
	required := RequiredOIDCGrantTypes(responseTypes, grantTypes)
	return ContainsOIDCGrantTypes(required, grantTypes)
}

func RequiredOIDCGrantTypes(responseTypes []OIDCResponseType, grantTypesSet []OIDCGrantType) (grantTypes []OIDCGrantType) {
	var implicit bool

	for _, r := range responseTypes {
		switch r {
		case OIDCResponseTypeCode:
			// #5684 when "Device Code" is selected, "Authorization Code" is no longer a hard requirement
			if !containsOIDCGrantType(grantTypesSet, OIDCGrantTypeDeviceCode) {
				grantTypes = append(grantTypes, OIDCGrantTypeAuthorizationCode)
			} else {
				grantTypes = append(grantTypes, OIDCGrantTypeDeviceCode)
			}
		case OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken:
			if !implicit {
				implicit = true
				grantTypes = append(grantTypes, OIDCGrantTypeImplicit)
			}
		}
	}

	return grantTypes
}

func (a *OIDCApp) getRequiredGrantTypes() []OIDCGrantType {
	return RequiredOIDCGrantTypes(a.ResponseTypes, a.GrantTypes)
}

func ContainsOIDCGrantTypes(shouldContain, list []OIDCGrantType) bool {
	for _, should := range shouldContain {
		if !containsOIDCGrantType(list, should) {
			return false
		}
	}
	return true
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
	return &Compliance{
		NoneCompliant: true,
		Problems:      []string{"Application.OIDC.UnsupportedVersion"},
	}
}

func GetOIDCV1Compliance(appType OIDCApplicationType, grantTypes []OIDCGrantType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	compliance := &Compliance{NoneCompliant: false}

	checkGrantTypesCombination(compliance, grantTypes)
	checkRedirectURIs(compliance, grantTypes, appType, redirectUris)
	checkApplicationType(compliance, appType, authMethod)

	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
	}
	return compliance
}

func checkGrantTypesCombination(compliance *Compliance, grantTypes []OIDCGrantType) {
	if !containsOIDCGrantType(grantTypes, OIDCGrantTypeDeviceCode) && containsOIDCGrantType(grantTypes, OIDCGrantTypeRefreshToken) && !containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode) {
		compliance.NoneCompliant = true
		compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.GrantType.Refresh.NoAuthCode")
	}
}

func checkRedirectURIs(compliance *Compliance, grantTypes []OIDCGrantType, appType OIDCApplicationType, redirectUris []string) {
	// See #5684 for OIDCGrantTypeDeviceCode and redirectUris further explanation
	if len(redirectUris) == 0 && (!containsOIDCGrantType(grantTypes, OIDCGrantTypeDeviceCode) || (containsOIDCGrantType(grantTypes, OIDCGrantTypeDeviceCode) && containsOIDCGrantType(grantTypes, OIDCGrantTypeAuthorizationCode))) {
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
}

func checkApplicationType(compliance *Compliance, appType OIDCApplicationType, authMethod OIDCAuthMethodType) {
	switch appType {
	case OIDCApplicationTypeNative:
		GetOIDCV1NativeApplicationCompliance(compliance, authMethod)
	case OIDCApplicationTypeUserAgent:
		GetOIDCV1UserAgentApplicationCompliance(compliance, authMethod)
	}
	if compliance.NoneCompliant {
		compliance.Problems = append([]string{"Application.OIDC.V1.NotCompliant"}, compliance.Problems...)
	}
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
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost")
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
				compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost")
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
			compliance.Problems = append(compliance.Problems, "Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost")
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

func OIDCOriginAllowList(redirectURIs, additionalOrigins []string) ([]string, error) {
	allowList := make([]string, 0)
	for _, redirect := range redirectURIs {
		origin, err := http_util.GetOriginFromURLString(redirect)
		if err != nil {
			return nil, err
		}
		if !http_util.IsOriginAllowed(allowList, origin) {
			allowList = append(allowList, origin)
		}
	}
	for _, origin := range additionalOrigins {
		if !http_util.IsOriginAllowed(allowList, origin) {
			allowList = append(allowList, origin)
		}
	}
	return allowList, nil
}
