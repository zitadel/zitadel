package domain

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type GroupMember struct {
	es_models.ObjectRoot

	UserID  string
	GroupID string
}

type GroupMemberState int32

const (
	GroupMemberStateUnspecified GroupMemberState = iota
	GroupMemberStateActive
	GroupMemberStateInactive
	GroupMemberStateRemoved

	groupMemberStateMax
)

func (s GroupMemberState) Valid() bool {
	return s > GroupMemberStateUnspecified && s < groupMemberStateMax
}

func (i *GroupMember) IsValid() bool {
	return i.AggregateID != "" && i.UserID != "" && i.GroupID != ""
}
