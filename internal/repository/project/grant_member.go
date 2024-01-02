package project

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	UniqueProjectGrantMemberType  = "project_grant_member"
	GrantMemberAddedType          = grantEventTypePrefix + member.AddedEventType
	GrantMemberChangedType        = grantEventTypePrefix + member.ChangedEventType
	GrantMemberRemovedType        = grantEventTypePrefix + member.RemovedEventType
	GrantMemberCascadeRemovedType = grantEventTypePrefix + member.CascadeRemovedEventType
)

func NewAddProjectGrantMemberUniqueConstraint(projectID, userID, grantID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueProjectGrantMemberType,
		fmt.Sprintf("%s:%s:%s", projectID, userID, grantID),
		"Errors.Project.Member.AlreadyExists")
}

func NewRemoveProjectGrantMemberUniqueConstraint(projectID, userID, grantID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueProjectGrantMemberType,
		fmt.Sprintf("%s:%s:%s", projectID, userID, grantID),
	)
}

type GrantMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	UserID  string   `json:"userId"`
	GrantID string   `json:"grantId"`
}

func (e *GrantMemberAddedEvent) Payload() interface{} {
	return e
}

func (e *GrantMemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddProjectGrantMemberUniqueConstraint(e.Aggregate().ID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	grantID string,
	roles ...string,
) *GrantMemberAddedEvent {
	return &GrantMemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberAddedType,
		),
		UserID:  userID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func GrantMemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantMemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-9f0sf", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantMemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	GrantID string   `json:"grantId"`
	UserID  string   `json:"userId"`
}

func (e *GrantMemberChangedEvent) Payload() interface{} {
	return e
}

func (e *GrantMemberChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewProjectGrantMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	grantID string,
	roles ...string,
) *GrantMemberChangedEvent {
	return &GrantMemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberChangedType,
		),
		UserID:  userID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func GrantMemberChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-39fi8", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID  string `json:"userId"`
	GrantID string `json:"grantId"`
}

func (e *GrantMemberRemovedEvent) Payload() interface{} {
	return e
}

func (e *GrantMemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectGrantMemberUniqueConstraint(e.Aggregate().ID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	grantID string,
) *GrantMemberRemovedEvent {
	return &GrantMemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberRemovedType,
		),
		UserID:  userID,
		GrantID: grantID,
	}
}

func GrantMemberRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantMemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-173fM", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantMemberCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID  string `json:"userId"`
	GrantID string `json:"grantId"`
}

func (e *GrantMemberCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *GrantMemberCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectGrantMemberUniqueConstraint(e.Aggregate().ID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	grantID string,
) *GrantMemberCascadeRemovedEvent {
	return &GrantMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberCascadeRemovedType,
		),
		UserID:  userID,
		GrantID: grantID,
	}
}

func GrantMemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-3kfs3", "unable to unmarshal label policy")
	}

	return e, nil
}
