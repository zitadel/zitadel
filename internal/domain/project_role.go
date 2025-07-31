package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"slices"
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
		if !containsRoleKey(roleKey, validRoles) {
			return true
		}
	}
	return false
}

func containsRoleKey(roleKey string, validRoles []string) bool {
	return slices.Contains(validRoles, roleKey)
}
