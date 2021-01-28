package domain

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type ProjectGrantMember struct {
	es_models.ObjectRoot

	GrantID string
	UserID  string
	Roles   []string
}

func NewProjectGrantMember(aggregateID, userID, grantID string, roles ...string) *ProjectGrantMember {
	return &ProjectGrantMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: aggregateID,
		},
		GrantID: grantID,
		UserID:  userID,
		Roles:   roles,
	}
}

func (i *ProjectGrantMember) IsValid() bool {
	return i.AggregateID != "" && i.GrantID != "" && i.UserID != "" && len(i.Roles) != 0
}
