package domain

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Member struct {
	es_models.ObjectRoot

	UserID string
	Roles  []string
}

func NewMember(aggregateID, userID string, roles ...string) *Member {
	return &Member{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: aggregateID,
		},
		UserID: userID,
		Roles:  roles,
	}
}

func (i *Member) IsValid() bool {
	return i.AggregateID != "" && i.UserID != "" && len(i.Roles) != 0
}

type MemberState int32

const (
	MemberStateUnspecified MemberState = iota
	MemberStateActive
	MemberStateRemoved

	memberStateCount
)

func (f MemberState) Valid() bool {
	return f >= 0 && f < memberStateCount
}
