package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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
	OIDCRESPONSETYPE_CODE OIDCResponseType = iota
	OIDCRESPONSETYPE_ID_TOKEN
	OIDCRESPONSETYPE_TOKEN
)

type OIDCGrantType int32

const (
	OIDCGRANTTYPE_AUTHORIZATION_CODE OIDCGrantType = iota
	OIDCGRANTTYPE_IMPLICIT
	OIDCGRANTTYPE_REFRESH_TOKEN
)

type OIDCApplicationType int32

const (
	OIDCAPPLICATIONTYPE_WEB OIDCApplicationType = iota
	OIDCAPPLICATIONTYPE_USER_AGENT
	OIDCAPPLICATIONTYPE_NATIVE
)

type OIDCAuthMethodType int32

const (
	OIDCAUTHMETHODTYPE_BASIC OIDCAuthMethodType = iota
	OIDCAUTHMETHODTYPE_POST
	OIDCAUTHMETHODTYPE_NONE
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

func (c *OIDCConfig) getRequiredGrantTypes() []OIDCGrantType {
	grantTypes := make([]OIDCGrantType, 0)
	implicit := false
	for _, r := range c.ResponseTypes {
		switch r {
		case OIDCRESPONSETYPE_CODE:
			grantTypes = append(grantTypes, OIDCGRANTTYPE_AUTHORIZATION_CODE)
		case OIDCRESPONSETYPE_ID_TOKEN, OIDCRESPONSETYPE_TOKEN:
			if !implicit {
				grantTypes = append(grantTypes, OIDCGRANTTYPE_IMPLICIT)
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
