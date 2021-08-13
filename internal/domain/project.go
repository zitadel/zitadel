package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Project struct {
	models.ObjectRoot

	State                ProjectState
	Name                 string
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
	OrgGrantCheck        bool
}

type ProjectState int32

const (
	ProjectStateUnspecified ProjectState = iota
	ProjectStateActive
	ProjectStateInactive
	ProjectStateRemoved
)

func (o *Project) IsValid() bool {
	return o.Name != ""
}
