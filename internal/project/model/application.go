package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Application struct {
	es_models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	Type       AppType
	OIDCConfig *OIDCConfig
	APIConfig  *APIConfig
}
type ApplicationChanges struct {
	Changes      []*ApplicationChange
	LastSequence uint64
}

type ApplicationChange struct {
	ChangeDate        *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType         string               `json:"eventType,omitempty"`
	Sequence          uint64               `json:"sequence,omitempty"`
	ModifierId        string               `json:"modifierUser,omitempty"`
	ModifierName      string               `json:"-"`
	ModifierLoginName string               `json:"-"`
	Data              interface{}          `json:"data,omitempty"`
}

type AppState int32

const (
	AppStateActive AppState = iota
	AppStateInactive
	AppStateRemoved
)

type AppType int32

const (
	AppTypeUnspecified AppType = iota
	AppTypeOIDC
	AppTypeSAML
	AppTypeAPI
)

func NewApplication(projectID, appID string) *Application {
	return &Application{ObjectRoot: es_models.ObjectRoot{AggregateID: projectID}, AppID: appID, State: AppStateActive}
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
	if a.Type == AppTypeAPI && !a.APIConfig.IsValid() {
		return false
	}
	return true
}

func (a *Application) GetKey(keyID string) (int, *ClientKey) {
	if a.OIDCConfig == nil {
		return -1, nil
	}
	for i, k := range a.OIDCConfig.ClientKeys {
		if k.KeyID == keyID {
			return i, k
		}
	}
	return -1, nil
}
