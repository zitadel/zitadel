package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Application struct {
	models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	Type       AppType
	OIDCConfig *OIDCConfig
}

type AppState int32

const (
	AppStateUnspecified AppState = iota
	AppStateActive
	AppStateInactive
	AppStateRemoved
)

type AppType int32

const (
	AppTypeUnspecified AppType = iota
	AppTypeOIDC
	AppTypeSAML
)

func NewApplication(projectID, appID string) *Application {
	return &Application{ObjectRoot: models.ObjectRoot{AggregateID: projectID}, AppID: appID, State: AppStateActive}
}

func (a *Application) IsValid(includeConfig bool) bool {
	if a.Name == "" || a.AggregateID == "" {
		return false
	}
	if !includeConfig {
		return true
	}
	if a.Type == AppTypeOIDC && !a.OIDCConfig.IsValid() {
		return false
	}
	return true
}
