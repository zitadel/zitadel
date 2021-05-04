package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
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
	ClientKeys               []*ClientKey
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

type ClientKey struct {
	es_models.ObjectRoot

	ApplicationID  string
	ClientID       string
	KeyID          string
	Type           key_model.AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
}

type Token struct {
	es_models.ObjectRoot

	TokenID    string
	ClientID   string
	Audience   []string
	Expiration time.Time
	Scopes     []string
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

//ClientID random_number@projectname (eg. 495894098234@zitadel)
func (c *OIDCConfig) GenerateNewClientID(idGenerator id.Generator, project *Project) error {
	rndID, err := idGenerator.Next()
	if err != nil {
		return err
	}

	c.ClientID = fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower(project.Name), " ", "_"))
	return nil
}

func (c *OIDCConfig) GenerateClientSecretIfNeeded(generator crypto.Generator) (string, error) {
	if c.AuthMethodType == OIDCAuthMethodTypeBasic || c.AuthMethodType == OIDCAuthMethodTypePost {
		return c.GenerateNewClientSecret(generator)
	}
	return "", nil
}

func (c *OIDCConfig) GenerateNewClientSecret(generator crypto.Generator) (string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		logging.Log("MODEL-UpnTI").OnError(err).Error("unable to create client secret")
		return "", errors.ThrowInternal(err, "MODEL-gH2Wl", "Errors.Project.CouldNotGenerateClientSecret")
	}
	c.ClientSecret = cryptoValue
	return stringSecret, nil
}

func (c *OIDCConfig) FillCompliance() {
	c.Compliance = GetOIDCCompliance(c.OIDCVersion, c.ApplicationType, c.GrantTypes, c.ResponseTypes, c.AuthMethodType, c.RedirectUris)
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
