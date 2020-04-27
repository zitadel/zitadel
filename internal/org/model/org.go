package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Org struct {
	es_models.ObjectRoot

	State  OrgState
	Name   string
	Domain string

	Members []*OrgMember
}

type OrgState int32

const (
	ORGSTATE_ACTIVE OrgState = iota
	ORGSTATE_INACTIVE
)

func NewOrg(id string) *Org {
	return &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: id}, State: ORGSTATE_ACTIVE}
}

func (o *Org) IsActive() bool {
	return o.State == ORGSTATE_ACTIVE
}

func (o *Org) IsValid() bool {
	return o.Name != "" && o.Domain != ""
}

func (o *Org) ContainsMember(userID string) bool {
	for _, member := range o.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}
