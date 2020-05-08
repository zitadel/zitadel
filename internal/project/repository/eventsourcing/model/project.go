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
		ObjectRoot:   project.ObjectRoot,
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
		ObjectRoot:   project.ObjectRoot,
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
	case ProjectAdded, ProjectChanged:
		return p.AppendAddProjectEvent(event)
	case ProjectDeactivated:
		return p.appendDeactivatedEvent()
	case ProjectReactivated:
		return p.appendReactivatedEvent()
	case ProjectMemberAdded:
		return p.appendAddMemberEvent(event)
	case ProjectMemberChanged:
		return p.appendChangeMemberEvent(event)
	case ProjectMemberRemoved:
		return p.appendRemoveMemberEvent(event)
	case ProjectRoleAdded:
		return p.appendAddRoleEvent(event)
	case ProjectRoleChanged:
		return p.appendChangeRoleEvent(event)
	case ProjectRoleRemoved:
		return p.appendRemoveRoleEvent(event)
	case ApplicationAdded:
		return p.appendAddAppEvent(event)
	case ApplicationChanged:
		return p.appendChangeAppEvent(event)
	case ApplicationRemoved:
		return p.appendRemoveAppEvent(event)
	case ApplicationDeactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_INACTIVE)
	case ApplicationReactivated:
		return p.appendAppStateEvent(event, model.APPSTATE_ACTIVE)
	case OIDCConfigAdded:
		return p.appendAddOIDCConfigEvent(event)
	case OIDCConfigChanged, OIDCConfigSecretChanged:
		return p.appendChangeOIDCConfigEvent(event)
	case ProjectGrantAdded:
		return p.appendAddGrantEvent(event)
	case ProjectGrantChanged:
		return p.appendChangeGrantEvent(event)
	case ProjectGrantDeactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_INACTIVE)
	case ProjectGrantReactivated:
		return p.appendGrantStateEvent(event, model.PROJECTGRANTSTATE_ACTIVE)
	case ProjectGrantRemoved:
		return p.appendRemoveGrantEvent(event)
	case ProjectGrantMemberAdded:
		return p.appendAddGrantMemberEvent(event)
	case ProjectGrantMemberChanged:
		return p.appendChangeGrantMemberEvent(event)
	case ProjectGrantMemberRemoved:
		return p.appendRemoveGrantMemberEvent(event)
	}
	return nil
}

func (p *Project) AppendAddProjectEvent(event *es_models.Event) error {
	p.setData(event)
	p.State = int32(model.PROJECTSTATE_ACTIVE)
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

func (p *Project) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-lo9sr").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
