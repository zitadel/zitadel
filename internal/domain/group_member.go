package domain

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type GroupMember struct {
	es_models.ObjectRoot

	GroupID string
	Roles   []string
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
	return i.AggregateID != "" && i.GroupID != "" && len(i.Roles) != 0
}

func NewGroupMember(aggregateID, groupID string, roles ...string) *GroupMember {
	return &GroupMember{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: aggregateID,
		},
		GroupID: groupID,
		Roles:   roles,
	}
}
