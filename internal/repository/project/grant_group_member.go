package project

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/groupmember"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	UniqueProjectGrantGroupMemberType  = "project_grant_group_member"
	GrantGroupMemberAddedType          = grantEventTypePrefix + groupmember.GroupAddedEventType
	GrantGroupMemberChangedType        = grantEventTypePrefix + groupmember.GroupChangedEventType
	GrantGroupMemberRemovedType        = grantEventTypePrefix + groupmember.GroupRemovedEventType
	GrantGroupMemberCascadeRemovedType = grantEventTypePrefix + groupmember.GroupCascadeRemovedEventType
)

func NewAddProjectGrantGroupMemberUniqueConstraint(projectID, groupID, grantID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueProjectGrantGroupMemberType,
		fmt.Sprintf("%s:%s:%s", projectID, groupID, grantID),
		"Errors.Project.GroupMember.AlreadyExists")
}

func NewRemoveProjectGrantGroupMemberUniqueConstraint(projectID, groupID, grantID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueProjectGrantGroupMemberType,
		fmt.Sprintf("%s:%s:%s", projectID, groupID, grantID),
	)
}

type GrantGroupMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	GroupID string   `json:"groupId"`
	GrantID string   `json:"grantId"`
}

func (e *GrantGroupMemberAddedEvent) Payload() interface{} {
	return e
}

func (e *GrantGroupMemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddProjectGrantGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID, e.GrantID)}
}

func NewProjectGrantGroupMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	grantID string,
	roles ...string,
) *GrantGroupMemberAddedEvent {
	return &GrantGroupMemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantGroupMemberAddedType,
		),
		GroupID: groupID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func GrantGroupMemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantGroupMemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-9f0sf", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantGroupMemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	GrantID string   `json:"grantId"`
	GroupID string   `json:"userId"`
}

func (e *GrantGroupMemberChangedEvent) Payload() interface{} {
	return e
}

func (e *GrantGroupMemberChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewProjectGrantGroupMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	grantID string,
	roles ...string,
) *GrantGroupMemberChangedEvent {
	return &GrantGroupMemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantGroupMemberChangedType,
		),
		GroupID: groupID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func GrantGroupMemberChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantGroupMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-39fi8", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantGroupMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string `json:"userId"`
	GrantID string `json:"grantId"`
}

func (e *GrantGroupMemberRemovedEvent) Payload() interface{} {
	return e
}

func (e *GrantGroupMemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectGrantGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID, e.GrantID)}
}

func NewProjectGrantGroupMemberRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	grantID string,
) *GrantGroupMemberRemovedEvent {
	return &GrantGroupMemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantGroupMemberRemovedType,
		),
		GroupID: groupID,
		GrantID: grantID,
	}
}

func GrantGroupMemberRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantGroupMemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-173fM", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantGroupMemberCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string `json:"userId"`
	GrantID string `json:"grantId"`
}

func (e *GrantGroupMemberCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *GrantGroupMemberCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectGrantGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID, e.GrantID)}
}

func NewProjectGrantGroupMemberCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	grantID string,
) *GrantGroupMemberCascadeRemovedEvent {
	return &GrantGroupMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantGroupMemberCascadeRemovedType,
		),
		GroupID: groupID,
		GrantID: grantID,
	}
}

func GrantGroupMemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantGroupMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-3kfs3", "unable to unmarshal label policy")
	}

	return e, nil
}
