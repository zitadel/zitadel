package model

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type Project struct {
	es_models.ObjectRoot

	State                  ProjectState
	Name                   string
	Members                []*ProjectMember
	Roles                  []*ProjectRole
	Applications           []*Application
	Grants                 []*ProjectGrant
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

type ProjectState int32

const (
	ProjectStateActive ProjectState = iota
	ProjectStateInactive
	ProjectStateRemoved
)

func (p *Project) IsActive() bool {
	return p.State == ProjectStateActive
}

func (p *Project) IsValid() bool {
	return p.Name != ""
}

func (p *Project) ContainsRole(role *ProjectRole) bool {
	for _, r := range p.Roles {
		if r.Key == role.Key {
			return true
		}
	}
	return false
}

func (p *Project) GetApp(appID string) (int, *Application) {
	for i, a := range p.Applications {
		if a.AppID == appID {
			return i, a
		}
	}
	return -1, nil
}

func (p *Project) GetGrant(grantID string) (int, *ProjectGrant) {
	for i, g := range p.Grants {
		if g.GrantID == grantID {
			return i, g
		}
	}
	return -1, nil
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

func (p *Project) ContainsGrantMember(member *ProjectGrantMember) bool {
	for _, g := range p.Grants {
		if g.GrantID != member.GrantID {
			continue
		}
		if _, m := g.GetMember(member.UserID); m != nil {
			return true
		}
	}
	return false
}
