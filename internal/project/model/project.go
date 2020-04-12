package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	in_model "github.com/caos/zitadel/internal/model"
)

type Project struct {
	es_models.ObjectRoot

	State   ProjectState
	Name    string
	Members []*ProjectMember
	Roles   []*ProjectRole
}

type ProjectState in_model.Enum

var states = []string{"Active", "Inactive"}

func NewProject(id string) *Project {
	return &Project{ObjectRoot: es_models.ObjectRoot{ID: id}, State: Active}
}

func (p *Project) IsActive() bool {
	if p.State == Active {
		return true
	}
	return false
}

func (p *Project) IsValid() bool {
	if p.Name == "" {
		return false
	}
	return true
}

func (p *Project) ContainsMember(member *ProjectMember) bool {
	for _, m := range p.Members {
		if m.UserID == member.UserID {
			return true
		}
	}
	return false
}

func (p *Project) ContainsRole(role *ProjectRole) bool {
	for _, r := range p.Roles {
		if r.Key == role.Key {
			return true
		}
	}
	return false
}
