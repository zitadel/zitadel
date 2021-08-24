package model

import (
	"encoding/json"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Application struct {
	es_models.ObjectRoot
	AppID      string      `json:"appId"`
	State      int32       `json:"-"`
	Name       string      `json:"name,omitempty"`
	Type       int32       `json:"appType,omitempty"`
	OIDCConfig *OIDCConfig `json:"-"`
	APIConfig  *APIConfig  `json:"-"`
}

type ApplicationID struct {
	es_models.ObjectRoot
	AppID string `json:"appId"`
}

func GetApplication(apps []*Application, id string) (int, *Application) {
	for i, a := range apps {
		if a.AppID == id {
			return i, a
		}
	}
	return -1, nil
}

func (a *Application) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-8die3").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
