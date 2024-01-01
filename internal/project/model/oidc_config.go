package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type OIDCConfig struct {
	es_models.ObjectRoot
	AppID                    string
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

type Token struct {
	es_models.ObjectRoot

	TokenID    string
	ClientID   string
	Audience   []string
	Expiration time.Time
	Scopes     []string
}

func GetOIDCCompliance(version OIDCVersion, appType OIDCApplicationType, grantTypes []OIDCGrantType, responseTypes []OIDCResponseType, authMethod OIDCAuthMethodType, redirectUris []string) *Compliance {
	switch version {
	case OIDCVersionV1:
		domainGrantTypes := make([]domain.OIDCGrantType, len(grantTypes))
		for i, grantType := range grantTypes {
			domainGrantTypes[i] = domain.OIDCGrantType(grantType)
		}
		compliance := domain.GetOIDCV1Compliance(domain.OIDCApplicationType(appType), domainGrantTypes, domain.OIDCAuthMethodType(authMethod), redirectUris)
		return &Compliance{
			NoneCompliant: compliance.NoneCompliant,
			Problems:      compliance.Problems,
		}
	}
	return nil
}
