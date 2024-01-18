package model

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	State        ProjectGrantState
	RoleKeys     []string
	Members      []*ProjectGrantMember
}

type ProjectGrantState int32

const (
	ProjectGrantStateActive ProjectGrantState = iota
	ProjectGrantStateInactive
)
