package domain

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type GroupGrant struct {
	es_models.ObjectRoot

	State          GroupGrantState
	GroupID        string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
}

type GroupGrantState int32

const (
	GroupGrantStateUnspecified GroupGrantState = iota
	GroupGrantStateActive
	GroupGrantStateInactive
	GroupGrantStateRemoved
)

func (g *GroupGrant) IsValid() bool {
	return g.ProjectID != "" && g.GroupID != ""
}

func (g *GroupGrant) HasInvalidRoles(validRoles []string) bool {
	for _, roleKey := range g.RoleKeys {
		if !containsRoleKey(roleKey, validRoles) {
			return true
		}
	}
	return false
}
