package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Application struct {
	es_models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	Type       AppType
	OIDCConfig *OIDCConfig
}

type AppState int32

const (
	APPSTATE_ACTIVE AppState = iota
	APPSTATE_INACTIVE
)

type AppType int32

const (
	APPTYPE_UNDEFINED AppType = iota
	APPTYPE_OIDC
	APPTYPE_SAML
)

func NewApplication(projectID, appID string) *Application {
	return &Application{ObjectRoot: es_models.ObjectRoot{ID: projectID}, AppID: appID, State: APPSTATE_ACTIVE}
}

func (a *Application) IsValid(includeConfig bool) bool {
	if a.Name == "" || a.ID == "" {
		return false
	}
	if !includeConfig {
		return true
	}
	if a.Type == APPTYPE_OIDC && !a.OIDCConfig.IsValid() {
		return false
	}
	return true
}
