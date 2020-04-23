package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectGrantMember struct {
	es_models.ObjectRoot
	GrantID string   `json:"grantId,omitempty"`
	UserID  string   `json:"userId,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

func GetProjectGrantMember(members []*ProjectGrantMember, id string) (int, *ProjectGrantMember) {
	for i, m := range members {
		if m.UserID == id {
			return i, m
		}
	}
	return -1, nil
}

func GrantMembersToModel(members []*ProjectGrantMember) []*model.ProjectGrantMember {
	convertedMembers := make([]*model.ProjectGrantMember, len(members))
	for i, g := range members {
		convertedMembers[i] = GrantMemberToModel(g)
	}
	return convertedMembers
}

func GrantMembersFromModel(members []*model.ProjectGrantMember) []*ProjectGrantMember {
	convertedMembers := make([]*ProjectGrantMember, len(members))
	for i, g := range members {
		convertedMembers[i] = GrantMemberFromModel(g)
	}
	return convertedMembers
}

func GrantMemberFromModel(member *model.ProjectGrantMember) *ProjectGrantMember {
	return &ProjectGrantMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.ObjectRoot.AggregateID,
			Sequence:     member.Sequence,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
		},
		GrantID: member.GrantID,
		UserID:  member.UserID,
		Roles:   member.Roles,
	}
}

func GrantMemberToModel(member *ProjectGrantMember) *model.ProjectGrantMember {
	return &model.ProjectGrantMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  member.AggregateID,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
			Sequence:     member.Sequence,
		},
		GrantID: member.GrantID,
		UserID:  member.UserID,
		Roles:   member.Roles,
	}
}

func (p *Project) appendAddGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate

	if _, g := GetProjectGrant(p.Grants, member.GrantID); g != nil {
		g.Members = append(g.Members, member)
	}
	return nil
}

func (p *Project) appendChangeGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}
	if _, g := GetProjectGrant(p.Grants, member.GrantID); g != nil {
		if i, m := GetProjectGrantMember(g.Members, member.UserID); m != nil {
			g.Members[i].getData(event)
		}
	}
	return nil
}

func (p *Project) appendRemoveGrantMemberEvent(event *es_models.Event) error {
	member := &ProjectGrantMember{}
	err := member.getData(event)
	if err != nil {
		return err
	}

	if _, g := GetProjectGrant(p.Grants, member.GrantID); g != nil {
		if i, member := GetProjectGrantMember(g.Members, member.UserID); member != nil {
			g.Members[i] = g.Members[len(g.Members)-1]
			g.Members[len(g.Members)-1] = nil
			g.Members = g.Members[:len(g.Members)-1]
		}
	}
	return nil
}

func (m *ProjectGrantMember) getData(event *es_models.Event) error {
	m.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, m); err != nil {
		logging.Log("EVEN-8die2").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
