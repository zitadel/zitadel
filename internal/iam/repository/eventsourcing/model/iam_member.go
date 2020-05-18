package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

type IamMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func GetIamMember(members []*IamMember, id string) (int, *IamMember) {
	for i, m := range members {
		if m.UserID == id {
			return i, m
		}
	}
	return -1, nil
}

func IamMembersToModel(members []*IamMember) []*model.IamMember {
	convertedMembers := make([]*model.IamMember, len(members))
	for i, m := range members {
		convertedMembers[i] = IamMemberToModel(m)
	}
	return convertedMembers
}

func IamMembersFromModel(members []*model.IamMember) []*IamMember {
	convertedMembers := make([]*IamMember, len(members))
	for i, m := range members {
		convertedMembers[i] = IamMemberFromModel(m)
	}
	return convertedMembers
}

func IamMemberFromModel(member *model.IamMember) *IamMember {
	return &IamMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func IamMemberToModel(member *IamMember) *model.IamMember {
	return &model.IamMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func (iam *Iam) appendAddMemberEvent(event *es_models.Event) error {
	member := &IamMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	iam.Members = append(iam.Members, member)
	return nil
}

func (iam *Iam) appendChangeMemberEvent(event *es_models.Event) error {
	member := &IamMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	if i, m := GetIamMember(iam.Members, member.UserID); m != nil {
		iam.Members[i] = member
	}
	return nil
}

func (iam *Iam) appendRemoveMemberEvent(event *es_models.Event) error {
	member := &IamMember{}
	err := member.setData(event)
	if err != nil {
		return err
	}
	if i, m := GetIamMember(iam.Members, member.UserID); m != nil {
		iam.Members[i] = iam.Members[len(iam.Members)-1]
		iam.Members[len(iam.Members)-1] = nil
		iam.Members = iam.Members[:len(iam.Members)-1]
	}
	return nil
}

func (m *IamMember) setData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
