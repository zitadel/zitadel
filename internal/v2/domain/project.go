package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Project struct {
	models.ObjectRoot

	State   ProjectState
	Name    string
	Members []*Member
	Roles   []*ProjectRole
	//Applications         []*Application
	//Grants               []*ProjectGrant
	ProjectRoleAssertion bool
	ProjectRoleCheck     bool
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
