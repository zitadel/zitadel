package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	in_model "github.com/caos/zitadel/internal/model"
)

type Org struct {
	es_models.ObjectRoot

	State OrgState
	Name  string
}

type OrgState in_model.Enum

var states = []string{"Active", "Inactive"}

func NewOrg(id string) *Org {
	return &Org{ObjectRoot: es_models.ObjectRoot{ID: id}, State: Active}
}

func (p *Org) IsActive() bool {
	if p.State == Active {
		return true
	}
	return false
}

func (p *Org) IsValid() bool {
	if p.Name == "" {
		return false
	}
	return true
}
