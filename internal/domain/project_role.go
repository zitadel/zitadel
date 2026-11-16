package domain

import (
	"slices"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectRole struct {
	models.ObjectRoot

	Key         string
	DisplayName string
	Group       string
}

type ProjectRoleState int32

const (
	ProjectRoleStateUnspecified ProjectRoleState = iota
	ProjectRoleStateActive
	ProjectRoleStateRemoved
)

func (s ProjectRoleState) Exists() bool {
	return s != ProjectRoleStateUnspecified && s != ProjectRoleStateRemoved
}

func (p *ProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}

func HasInvalidRoles(validRoles, roles []string) bool {
	for _, roleKey := range roles {
		if !slices.Contains(validRoles, roleKey) {
			return true
		}
	}
	return false
}