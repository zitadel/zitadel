package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func (m *ProjectMember) SetData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
