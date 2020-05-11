package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
)

const (
	ProjectVersion = "v1"
)

type Project struct {
	es_models.ObjectRoot
	Name         string           `json:"name,omitempty"`
	State        int32            `json:"-"`
	Members      []*ProjectMember `json:"-"`
	Roles        []*ProjectRole   `json:"-"`
	Applications []*Application   `json:"-"`
	Grants       []*ProjectGrant  `json:"-"`
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
	roles := ProjectRolesFromModel(project.Roles)
	apps := AppsFromModel(project.Applications)
	grants := GrantsFromModel(project.Grants)
	return &Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.ObjectRoot.AggregateID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		Name:         project.Name,
		State:        int32(project.State),
		Members:      members,
		Roles:        roles,
		Applications: apps,
		Grants:       grants,
	}
}

func ProjectToModel(project *Project) *model.Project {
	members := ProjectMembersToModel(project.Members)
	roles := ProjectRolesToModel(project.Roles)
	apps := AppsToModel(project.Applications)
	grants := GrantsToModel(project.Grants)
	return &model.Project{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.AggregateID,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
			Sequence:     project.Sequence,
		},
		Name:         project.Name,
		State:        model.ProjectState(project.State),
		Members:      members,
		Roles:        roles,
		Applications: apps,
		Grants:       grants,
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
		p.State = int32(model.PROJECTSTATE_ACTIVE)
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
	case model.ProjectRoleAdded:
		return p.appendAddRoleEvent(event)
	case model.ProjectRoleChanged:
		return p.appendChangeRoleEvent(event)
	case model.ProjectRoleRemoved:
		return p.appendRemoveRoleEvent(event)
	case model.ApplicationAdded:
		return p.appendAddAppEvent(event)
	case model.ApplicationChanged:
		return p.appendChangeAppEvent(event)
	case model.ApplicationRemoved:
		return p.appendRemoveAppEvent(event)
	case model.ApplicationDeactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_INACTIVE)
	case model.ApplicationReactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_ACTIVE)
	case model.OIDCConfigAdded:
		return p.appendAddOIDCConfigEvent(event)
	case model.OIDCConfigChanged, model.OIDCConfigSecretChanged:
		return p.appendChangeOIDCConfigEvent(event)
	case model.ProjectGrantAdded:
		return p.appendAddGrantEvent(event)
	case model.ProjectGrantChanged:
		return p.appendChangeGrantEvent(event)
	case model.ProjectGrantDeactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_INACTIVE)
	case model.ProjectGrantReactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_ACTIVE)
	case model.ProjectGrantRemoved:
		return p.appendRemoveGrantEvent(event)
	case model.ProjectGrantMemberAdded:
		return p.appendAddGrantMemberEvent(event)
	case model.ProjectGrantMemberChanged:
		return p.appendChangeGrantMemberEvent(event)
	case model.ProjectGrantMemberRemoved:
		return p.appendRemoveGrantMemberEvent(event)
	}
	return nil
}

func (p *Project) appendDeactivatedEvent() error {
	p.State = int32(model.PROJECTSTATE_INACTIVE)
	return nil
}

func (p *Project) appendReactivatedEvent() error {
	p.State = int32(model.PROJECTSTATE_ACTIVE)
	return nil
}
