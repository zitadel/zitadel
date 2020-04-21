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
	OIDCRESPONSETYPE_TOKEN_ID_TOKEN
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
	code, idToken, tokenIdToken := c.getChosenResponseTypes()
	if code {
		ok := c.containsGrantType(OIDCGRANTTYPE_AUTHORIZATION_CODE)
		if !ok {
			return false
		}
	}
	if idToken {
		ok := c.containsGrantType(OIDCGRANTTYPE_IMPLICIT)
		if !ok {
			return false
		}
	}
	if tokenIdToken {
		ok := c.containsGrantType(OIDCGRANTTYPE_IMPLICIT)
		if !ok {
			return false
		}
	}
	return true
}

func (c *OIDCConfig) getChosenResponseTypes() (bool, bool, bool) {
	code := false
	idToken := false
	tokenIdToken := false
	for _, r := range c.ResponseTypes {
		switch r {
		case OIDCRESPONSETYPE_CODE:
			code = true
		case OIDCRESPONSETYPE_ID_TOKEN:
			idToken = true
		case OIDCRESPONSETYPE_TOKEN_ID_TOKEN:
			tokenIdToken = true
		}
	}
	return code, idToken, tokenIdToken
}

func (c *OIDCConfig) containsGrantType(grantType OIDCGrantType) bool {
	for _, t := range c.GrantTypes {
		if t == grantType {
			return true
		}
	}
	return false
}
