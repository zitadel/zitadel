package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Application struct {
	es_models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	Type       AppType
	OIDCConfig *OIDCConfig
}
type ApplicationChanges struct {
	Changes      []*ApplicationChange
	LastSequence uint64
}

type ApplicationChange struct {
	ChangeDate *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType  string               `json:"eventType,omitempty"`
	Sequence   uint64               `json:"sequence,omitempty"`
	Modifier   string               `json:"modifierUser,omitempty"`
	Data       interface{}          `json:"data,omitempty"`
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
	return &Application{ObjectRoot: es_models.ObjectRoot{AggregateID: projectID}, AppID: appID, State: APPSTATE_ACTIVE}
}

func (a *Application) IsValid(includeConfig bool) bool {
	if a.Name == "" || a.AggregateID == "" {
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
