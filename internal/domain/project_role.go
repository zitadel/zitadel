package domain

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
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

func NewProjectRole(projectID, key string) *ProjectRole {
	return &ProjectRole{ObjectRoot: models.ObjectRoot{AggregateID: projectID}, Key: key}
}

func (p *ProjectRole) IsValid() bool {
	return p.AggregateID != "" && p.Key != ""
}

func containsRoleKey(roleKey string, validRoles []string) bool {
	for _, validRole := range validRoles {
		if roleKey == validRole {
			return true
		}
	}
	return false
}
