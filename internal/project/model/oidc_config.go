package model

import (
	"fmt"
	"strings"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
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
}

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

func (c *OIDCConfig) IsValid() bool {
	grantTypes := c.getRequiredGrantTypes()
	for _, grantType := range grantTypes {
		ok := c.containsGrantType(grantType)
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
	if c.AuthMethodType == OIDCAuthMethodTypeNone {
		return "", nil
	}
	return c.GenerateNewClientSecret(generator)
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

func (c *OIDCConfig) getRequiredGrantTypes() []OIDCGrantType {
	grantTypes := make([]OIDCGrantType, 0)
	implicit := false
	for _, r := range c.ResponseTypes {
		switch r {
		case OIDCResponseTypeCode:
			grantTypes = append(grantTypes, OIDCGrantTypeAuthorizationCode)
		case OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken:
			if !implicit {
				grantTypes = append(grantTypes, OIDCGrantTypeImplicit)
			}
		}
	}
	return grantTypes
}

func (c *OIDCConfig) containsGrantType(grantType OIDCGrantType) bool {
	for _, t := range c.GrantTypes {
		if t == grantType {
			return true
		}
	}
	return false
}
