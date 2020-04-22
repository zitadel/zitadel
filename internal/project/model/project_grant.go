package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
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
	PROJECTGRANTSTATE_ACTIVE ProjectGrantState = iota
	PROJECTGRANTSTATE_INACTIVE
)

func NewProjectGrant(projectID, grantID string) *ProjectGrant {
	return &ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: projectID}, GrantID: grantID, State: PROJECTGRANTSTATE_ACTIVE}
}

func (p *ProjectGrant) IsActive() bool {
	if p.State == PROJECTGRANTSTATE_ACTIVE {
		return true
	}
	return false
}

func (p *ProjectGrant) IsValid() bool {
	if p.GrantedOrgID == "" {
		return false
	}
	return true
}

func (p *ProjectGrant) GetMember(userID string) (int, *ProjectGrantMember) {
	for i, m := range p.Members {
		if m.UserID == userID {
			return i, m
		}
	}
	return -1, nil
}
