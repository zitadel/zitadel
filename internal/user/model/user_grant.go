package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type UserGrant struct {
	es_models.ObjectRoot

	GrantID   string
	State     UserGrantState
	ProjectID string
	RoleKeys  []string
}

type UserGrantState int32

const (
	USERGRANTSTATE_ACTIVE UserGrantState = iota
	USERGRANTSTATE_INACTIVE
)

func (u *User) GetGrant(grantID string) (int, *UserGrant) {
	for i, g := range u.Grants {
		if g.GrantID == grantID {
			return i, g
		}
	}
	return -1, nil
}

func (u *User) ContainsGrantForProject(projectID string) (int, *UserGrant) {
	for i, g := range u.Grants {
		if g.ProjectID == projectID {
			return i, g
		}
	}
	return -1, nil
}

func (u *UserGrant) IsValid() bool {
	return u.ProjectID != ""
}
