package model

import (
	"encoding/json"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/org/model"
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
		return errors.ThrowInternal(err, "EVENT-Hz7Mb", "unable to unmarshal data")
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

func OrgMemberFromEvent(member *OrgMember, event *es_models.Event) (*OrgMember, error) {
	if member == nil {
		member = new(OrgMember)
	}
	member.ObjectRoot.AppendEvent(event)
	err := json.Unmarshal(event.Data, member)
	if err != nil {
		return nil, errors.ThrowInternal(err, "EVENT-D4qxo", "invalid event data")
	}
	return member, nil
}

func OrgMembersFromModel(members []*model.OrgMember) []*OrgMember {
	convertedMembers := make([]*OrgMember, len(members))
	for i, m := range members {
		convertedMembers[i] = OrgMemberFromModel(m)
	}
	return convertedMembers
}

func OrgMemberFromModel(member *model.OrgMember) *OrgMember {
	return &OrgMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func OrgMembersToModel(members []*OrgMember) []*model.OrgMember {
	convertedMembers := make([]*model.OrgMember, len(members))
	for i, m := range members {
		convertedMembers[i] = OrgMemberToModel(m)
	}
	return convertedMembers
}

func OrgMemberToModel(member *OrgMember) *model.OrgMember {
	return &model.OrgMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}
