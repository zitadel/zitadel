package eventsourcing

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"

	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

const (
	projectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name    string           `json:"name,omitempty"`
	State   int32            `json:"-"`
	Members []*ProjectMember `json:"-"`
}

type ProjectMember struct {
	es_models.ObjectRoot
	UserID string   `json:"userId,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

func (p *Project) Changes(changed *Project) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func ProjectFromModel(project *model.Project) *Project {
	members := ProjectMembersFromModel(project.Members)
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:    project.Name,
		State:   model.ProjectStateToInt(project.State),
		Members: members,
	}
}

func ProjectToModel(project *Project) *model.Project {
	members := ProjectMembersToModel(project.Members)
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:    project.Name,
		State:   model.ProjectStateFromInt(project.State),
		Members: members,
	}
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
		ObjectRoot: es_models.ObjectRoot{
			ID:           member.ObjectRoot.ID,
			Sequence:     member.Sequence,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectMemberToModel(member *ProjectMember) *model.ProjectMember {
	return &model.ProjectMember{
		ObjectRoot: es_models.ObjectRoot{
			ID:           member.ID,
			ChangeDate:   member.ChangeDate,
			CreationDate: member.CreationDate,
			Sequence:     member.Sequence,
		},
		UserID: member.UserID,
		Roles:  member.Roles,
	}
}

func ProjectFromEvents(project *Project, events ...*es_models.Event) (*Project, error) {
	if project == nil {
		project = &Project{}
	}

	return project, project.AppendEvents(events...)
}

func (p *Project) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.ProjectAdded, model.ProjectChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		p.State = model.ProjectStateToInt(model.Active)
		return nil
	case model.ProjectDeactivated:
		return p.appendDeactivatedEvent()
	case model.ProjectReactivated:
		return p.appendReactivatedEvent()
	case model.ProjectMemberAdded:
		return p.appendAddMemberEvent(event)
	case model.ProjectMemberChanged:
		return p.appendChangeMemberEvent(event)
	case model.ProjectMemberRemoved:
		return p.appendRemoveMemberEvent(event)
	}
	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.Inactive)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = model.ProjectStateToInt(model.Active)
	return nil
}

func (p *Project) appendAddMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return err
	}
	member.ObjectRoot.CreationDate = event.CreationDate
	p.Members = append(p.Members, member)
	return nil
}

func (p *Project) appendChangeMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return err
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = member
		}
	}
	return nil
}

func (p *Project) appendRemoveMemberEvent(event *es_models.Event) error {
	member, err := getMemberData(event)
	if err != nil {
		return err
	}
	for i, m := range p.Members {
		if m.UserID == member.UserID {
			p.Members[i] = p.Members[len(p.Members)-1]
			p.Members[len(p.Members)-1] = nil
			p.Members = p.Members[:len(p.Members)-1]
		}
	}
	return nil
}

func getMemberData(event *es_models.Event) (*ProjectMember, error) {
	member := &ProjectMember{}
	member.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, member); err != nil {
		logging.Log("EVEN-e4dkp").WithError(err).Error("could not unmarshal event data")
		return nil, errors.ThrowInternal(err, "EVENT-83js6", "could not unmarshal event data")
	}
	return member, nil
}
