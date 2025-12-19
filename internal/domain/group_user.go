package domain

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type GroupUser struct {
	es_models.ObjectRoot

	UserID     string
	Attributes []string
}

type GroupUserState int32

const (
	GroupUserStateUnspecified GroupUserState = iota
	GroupUserStateActive
	GroupUserStateInactive
	GroupUserStateRemoved

	groupUserStateMax
)

func (s GroupUserState) Valid() bool {
	return s > GroupUserStateUnspecified && s < groupUserStateMax
}

func (i *GroupUser) IsValid() bool {
	return i.AggregateID != "" && i.UserID != ""
}

func NewGroupUser(aggregateID, userID string, attributes ...string) *GroupUser {
	return &GroupUser{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: aggregateID,
		},
		UserID:     userID,
		Attributes: attributes,
	}
}
