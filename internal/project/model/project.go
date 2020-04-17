package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Project struct {
	es_models.ObjectRoot

	State        ProjectState
	Name         string
	Members      []*ProjectMember
	Roles        []*ProjectRole
	Applications []*Application
	Grants       []*ProjectGrant
}

type ProjectState int32

const (
	PROJECTSTATE_ACTIVE ProjectState = iota
	PROJECTSTATE_INACTIVE
)

func NewProject(id string) *Project {
	return &Project{ObjectRoot: es_models.ObjectRoot{ID: id}, State: PROJECTSTATE_ACTIVE}
}

func (p *Project) IsActive() bool {
	if p.State == PROJECTSTATE_ACTIVE {
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

func (p *Project) ContainsApp(app *Application) (bool, *Application) {
	for _, a := range p.Applications {
		if a.AppID == app.AppID {
			return true, a
		}
	}
	return false, nil
}

func (p *Project) ContainsGrant(grant *ProjectGrant) bool {
	for _, g := range p.Grants {
		if g.GrantID == grant.GrantID {
			return true
		}
	}
	return false
}

func (p *Project) ContainsGrantForOrg(orgID string) bool {
	for _, g := range p.Grants {
		if g.GrantedOrgID == orgID {
			return true
		}
	}
	return false
}

func (p *Project) ContainsRoles(roleKeys []string) bool {
	for _, r := range roleKeys {
		if !p.ContainsRole(&ProjectRole{Key: r}) {
			return false
		}
	}
	return true
}

func ProjectStateToInt(s ProjectState) int32 {
	return int32(s)
}

func ProjectStateFromInt(index int32) ProjectState {
	return ProjectState(index)
}
