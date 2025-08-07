package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Project struct {
	es_models.ObjectRoot

	State                  ProjectState
	Name                   string
	Members                []*ProjectMember
	Roles                  []*ProjectRole
	Applications           []*Application
	Grants                 []*ProjectGrant
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

type ProjectState int32

const (
	ProjectStateActive ProjectState = iota
	ProjectStateInactive
	ProjectStateRemoved
)
