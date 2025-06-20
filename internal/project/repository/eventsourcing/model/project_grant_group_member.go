package model

import (
	"encoding/json"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectGrantGroupMember struct {
	es_models.ObjectRoot
	GrantID string   `json:"grantId,omitempty"`
	GroupID string   `json:"groupId,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

func (m *ProjectGrantGroupMember) SetData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-8die2").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
