package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
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
