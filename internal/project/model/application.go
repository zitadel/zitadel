package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type Application struct {
	es_models.ObjectRoot

	AppID        string
	State        AppState
	CreationDate time.Time
	ChangeDate   time.Time
	ProjectID    string
	Name         string
	OIDCConfig   *OIDCConfig
}

type OIDCConfig struct {
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

type AppState int32

const (
	APPSTATE_ACTIVE AppState = iota
	APPSTATE_INACTIVE
)

type AppType int32

const (
	APPTYPE_OIDC AppType = iota
	APPTYPE_SAML
)

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

func NewApp(projectID, appID string) *Application {
	return &Application{ObjectRoot: es_models.ObjectRoot{ID: projectID}, AppID: appID, State: APPSTATE_ACTIVE}
}

func (a *Application) IsValid() bool {
	if a.Name == "" && a.ID == "" && a.OIDCConfig == nil && a.OIDCConfig.IsValid() {
		return false
	}
	return true
}

func (c *OIDCConfig) IsValid() bool {
	//TODO: Implement Validation check
	return true
}

func AppStateToInt(s AppState) int32 {
	return int32(s)
}

func AppStateFromInt(index int32) AppState {
	return AppState(index)
}
