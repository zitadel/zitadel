package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	in_model "github.com/caos/zitadel/internal/model"
)

type Org struct {
	es_models.ObjectRoot

	State  OrgState
	Name   string
	Domain string

	Members []*OrgMember
}

type OrgState in_model.Enum

var states = []string{"Active", "Inactive"}

func NewOrg(id string) *Org {
	return &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: id}, State: Active}
}

func (o *Org) IsActive() bool {
	return o.State == Active
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
