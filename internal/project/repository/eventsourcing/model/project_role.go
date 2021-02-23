package model

import (
	"encoding/json"
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
)

type ProjectRole struct {
	es_models.ObjectRoot
	Key         string `json:"key,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Group       string `json:"group,omitempty"`
}

func GetProjectRole(roles []*ProjectRole, key string) (int, *ProjectRole) {
	for i, r := range roles {
		if r.Key == key {
			return i, r
		}
	}
	return -1, nil
}

func ProjectRolesToModel(roles []*ProjectRole) []*model.ProjectRole {
	convertedRoles := make([]*model.ProjectRole, len(roles))
	for i, r := range roles {
		convertedRoles[i] = ProjectRoleToModel(r)
	}
	return convertedRoles
}

func ProjectRolesFromModel(roles []*model.ProjectRole) []*ProjectRole {
	convertedRoles := make([]*ProjectRole, len(roles))
	for i, r := range roles {
		convertedRoles[i] = ProjectRoleFromModel(r)
	}
	return convertedRoles
}

func ProjectRoleFromModel(role *model.ProjectRole) *ProjectRole {
	return &ProjectRole{
		ObjectRoot:  role.ObjectRoot,
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func ProjectRoleToModel(role *ProjectRole) *model.ProjectRole {
	return &model.ProjectRole{
		ObjectRoot:  role.ObjectRoot,
		Key:         role.Key,
		DisplayName: role.DisplayName,
		Group:       role.Group,
	}
}

func (p *Project) appendAddRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	role.ObjectRoot.CreationDate = event.CreationDate
	p.Roles = append(p.Roles, role)
	return nil
}

func (p *Project) appendChangeRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	if i, r := GetProjectRole(p.Roles, role.Key); r != nil {
		p.Roles[i] = role
	}
	return nil
}

func (p *Project) appendRemoveRoleEvent(event *es_models.Event) error {
	role := new(ProjectRole)
	err := role.setData(event)
	if err != nil {
		return err
	}
	if i, r := GetProjectRole(p.Roles, role.Key); r != nil {
		p.Roles[i] = p.Roles[len(p.Roles)-1]
		p.Roles[len(p.Roles)-1] = nil
		p.Roles = p.Roles[:len(p.Roles)-1]
	}
	return nil
}

func (r *ProjectRole) setData(event *es_models.Event) error {
	r.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-d9euw").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
