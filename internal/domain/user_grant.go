package domain

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type UserGrant struct {
	es_models.ObjectRoot

	State          UserGrantState
	UserID         string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
}

type UserGrantState int32

const (
	UserGrantStateUnspecified UserGrantState = iota
	UserGrantStateActive
	UserGrantStateInactive
	UserGrantStateRemoved
)

func (u *UserGrant) IsValid() bool {
	return u.ProjectID != "" && u.UserID != ""
}

func (g *UserGrant) HasInvalidRoles(validRoles []string) bool {
	for _, roleKey := range g.RoleKeys {
		if !containsRoleKey(roleKey, validRoles) {
			return true
		}
	}
	return false
}

func containsRoleKey(roleKey string, validRoles []string) bool {
	for _, validRole := range validRoles {
		if roleKey == validRole {
			return true
		}
	}
	return false
}
