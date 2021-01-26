package model

import (
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"

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
	ChangeDate   *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType    string               `json:"eventType,omitempty"`
	Sequence     uint64               `json:"sequence,omitempty"`
	ModifierId   string               `json:"modifierUser,omitempty"`
	ModifierName string               `json:"-"`
	Data         interface{}          `json:"data,omitempty"`
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
	return true
}

type ApplicationKey struct {
	es_models.ObjectRoot

	AppID          string
	KeyID          string
	Type           key_model.AuthNKeyType
	ExpirationDate time.Time
	PrivateKey     []byte
}
