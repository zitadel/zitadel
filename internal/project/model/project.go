package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Project struct {
	es_models.ObjectRoot

	State                ProjectState
	Name                 string
	Members              []*ProjectMember
	Roles                []*ProjectRole
	Applications         []*Application
	Grants               []*ProjectGrant
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
}
type ProjectChanges struct {
	Changes      []*ProjectChange
	LastSequence uint64
}

type ProjectChange struct {
	ChangeDate        *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType         string               `json:"eventType,omitempty"`
	Sequence          uint64               `json:"sequence,omitempty"`
	ModifierId        string               `json:"modifierUser,omitempty"`
	ModifierName      string               `json:"-"`
	ModifierLoginName string               `json:"-"`
	Data              interface{}          `json:"data,omitempty"`
}

type ProjectState int32

const (
	ProjectStateActive ProjectState = iota
	ProjectStateInactive
	ProjectStateRemoved
)

func NewProject(id string) *Project {
	return &Project{ObjectRoot: es_models.ObjectRoot{AggregateID: id}, State: ProjectStateActive}
}

func (p *Project) IsActive() bool {
	return p.State == ProjectStateActive
}

func (p *Project) IsValid() bool {
	return p.Name != ""
}

func (p *Project) GetMember(userID string) (int, *ProjectMember) {
	for i, m := range p.Members {
		if m.UserID == userID {
			return i, m
		}
	}
	return -1, nil
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
