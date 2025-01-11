package domain

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type ProjectGrantGroupMember struct {
	es_models.ObjectRoot

	GrantID string
	GroupID string
	Roles   []string
}

func (i *ProjectGrantGroupMember) IsValid() bool {
	return i.AggregateID != "" && i.GrantID != "" && i.GroupID != "" && len(i.Roles) != 0
}
