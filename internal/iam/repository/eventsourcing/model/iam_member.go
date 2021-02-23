package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

type IAMMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func GetIAMMember(members []*IAMMember, id string) (int, *IAMMember) {
	for i, m := range members {
		if m.UserID == id {
			return i, m
		}
	}
	return -1, nil
}

func IAMMembersToModel(members []*IAMMember) []*model.IAMMember {
	convertedMembers := make([]*model.IAMMember, len(members))
	for i, m := range members {
		convertedMembers[i] = IAMMemberToModel(m)
	}
	return convertedMembers
}

func IAMMembersFromModel(members []*model.IAMMember) []*IAMMember {
	convertedMembers := make([]*IAMMember, len(members))
	for i, m := range members {
		convertedMembers[i] = IAMMemberFromModel(m)
	}
	return convertedMembers
}

func IAMMemberFromModel(member *model.IAMMember) *IAMMember {
	return &IAMMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func IAMMemberToModel(member *IAMMember) *model.IAMMember {
	return &model.IAMMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func (iam *IAM) appendAddMemberEvent(event *es_models.Event) error {
	member := &IAMMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	iam.Members = append(iam.Members, member)
	return nil
}

func (iam *IAM) appendChangeMemberEvent(event *es_models.Event) error {
	member := &IAMMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetIAMMember(iam.Members, member.UserID); m != nil {
		iam.Members[i] = member
	}
	return nil
}

func (iam *IAM) appendRemoveMemberEvent(event *es_models.Event) error {
	member := &IAMMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetIAMMember(iam.Members, member.UserID); m != nil {
		iam.Members[i] = iam.Members[len(iam.Members)-1]
		iam.Members[len(iam.Members)-1] = nil
		iam.Members = iam.Members[:len(iam.Members)-1]
	}
	return nil
}

func (m *IAMMember) SetData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
