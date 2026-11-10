package domain

import (
	"slices"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	State        ProjectGrantState
	RoleKeys     []string
}

type ProjectGrantState int32

const (
	ProjectGrantStateUnspecified ProjectGrantState = iota
	ProjectGrantStateActive
	ProjectGrantStateInactive
	ProjectGrantStateRemoved

	projectGrantStateMax
)

func (s ProjectGrantState) Valid() bool {
	return s > ProjectGrantStateUnspecified && s < projectGrantStateMax
}

func (s ProjectGrantState) Exists() bool {
	return s != ProjectGrantStateUnspecified && s != ProjectGrantStateRemoved
}

func (p *ProjectGrant) IsValid() bool {
	return p.GrantedOrgID != ""
}

func GetRemovedRoles(existingRoles, newRoles []string) map[string]bool {
	removed := make(map[string]bool, 0)
	for _, role := range existingRoles {
		if !slices.Contains(newRoles, role) {
			removed[role] = true
		}
	}
	return removed
}
