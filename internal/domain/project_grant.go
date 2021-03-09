package domain

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type ProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	State        ProjectGrantState
	RoleKeys     []string
}

type ProjectGrantIDs struct {
	ProjectID string
	GrantID   string
}

type ProjectGrantState int32

const (
	ProjectGrantStateUnspecified ProjectGrantState = iota
	ProjectGrantStateActive
	ProjectGrantStateInactive
	ProjectGrantStateRemoved
)

func (p *ProjectGrant) IsValid() bool {
	return p.GrantedOrgID != ""
}

func (g *ProjectGrant) HasInvalidRoles(validRoles []string) bool {
	for _, roleKey := range g.RoleKeys {
		if !containsRoleKey(roleKey, validRoles) {
			return true
		}
	}
	return false
}

func GetRemovedRoles(existingRoles, newRoles []string) []string {
	removed := make([]string, 0)
	for _, role := range existingRoles {
		if !containsRoleKey(role, newRoles) {
			removed = append(removed, role)
		}
	}
	return removed
}
