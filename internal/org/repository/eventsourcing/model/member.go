package model

import (
	"encoding/json"
	"reflect"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OrgMember struct {
	es_models.ObjectRoot `json:"-"`

	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func (m *OrgMember) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := m.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *OrgMember) AppendEvent(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)

	return m.SetData(event)
}

func (m *OrgMember) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, m)
	if err != nil {
		return zerrors.ThrowInternal(err, "EVENT-Hz7Mb", "unable to unmarshal data")
	}
	return nil
}

func (m *OrgMember) Changes(updatedMember *OrgMember) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if !reflect.DeepEqual(m.Roles, updatedMember.Roles) {
		changes["roles"] = updatedMember.Roles
		changes["userId"] = m.UserID
	}

	return changes
}
