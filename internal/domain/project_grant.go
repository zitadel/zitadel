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

func GetRemovedRoles(existingRoles, newRoles []string) []string {
	removed := make([]string, 0)
	for _, role := range existingRoles {
		if !containsKey(newRoles, role) {
			removed = append(removed, role)
		}
	}
	return removed
}

func containsKey(roles []string, key string) bool {
	for _, role := range roles {
		if role == key {
			return true
		}
	}
	return false
}
