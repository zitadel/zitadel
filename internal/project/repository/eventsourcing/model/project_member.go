package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func GetProjectMember(members []*ProjectMember, id string) (int, *ProjectMember) {
	for i, m := range members {
		if m.UserID == id {
			return i, m
		}
	}
	return -1, nil
}

func ProjectMembersToModel(members []*ProjectMember) []*model.ProjectMember {
	convertedMembers := make([]*model.ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberToModel(m)
	}
	return convertedMembers
}

func ProjectMembersFromModel(members []*model.ProjectMember) []*ProjectMember {
	convertedMembers := make([]*ProjectMember, len(members))
	for i, m := range members {
		convertedMembers[i] = ProjectMemberFromModel(m)
	}
	return convertedMembers
}

func ProjectMemberFromModel(member *model.ProjectMember) *ProjectMember {
	return &ProjectMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func ProjectMemberToModel(member *ProjectMember) *model.ProjectMember {
	return &model.ProjectMember{
		ObjectRoot: member.ObjectRoot,
		UserID:     member.UserID,
		Roles:      member.Roles,
	}
}

func (p *Project) appendAddMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	p.Members = append(p.Members, member)
	return nil
}

func (p *Project) appendChangeMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetProjectMember(p.Members, member.UserID); m != nil {
		p.Members[i] = member
	}
	return nil
}

func (p *Project) appendRemoveMemberEvent(event *es_models.Event) error {
	member := &ProjectMember{}
	err := member.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetProjectMember(p.Members, member.UserID); m != nil {
		p.Members[i] = p.Members[len(p.Members)-1]
		p.Members[len(p.Members)-1] = nil
		p.Members = p.Members[:len(p.Members)-1]
	}
	return nil
}

func (m *ProjectMember) SetData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
