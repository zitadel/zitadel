package usergrant

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueUserGrant             = "user_grant"
	userGrantEventTypePrefix    = eventstore.EventType("user.grant.")
	UserGrantAddedType          = userGrantEventTypePrefix + "added"
	UserGrantChangedType        = userGrantEventTypePrefix + "changed"
	UserGrantCascadeChangedType = userGrantEventTypePrefix + "cascade.changed"
	UserGrantRemovedType        = userGrantEventTypePrefix + "removed"
	UserGrantCascadeRemovedType = userGrantEventTypePrefix + "cascade.removed"
	UserGrantDeactivatedType    = userGrantEventTypePrefix + "deactivated"
	UserGrantReactivatedType    = userGrantEventTypePrefix + "reactivated"
)

func NewAddUserGrantUniqueConstraint(resourceOwner, userID, projectID, projectGrantID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUserGrant,
		fmt.Sprintf("%s:%s:%s:%v", resourceOwner, userID, projectID, projectGrantID),
		"Errors.UserGrant.AlreadyExists")
}

func NewRemoveUserGrantUniqueConstraint(resourceOwner, userID, projectID, projectGrantID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUserGrant,
		fmt.Sprintf("%s:%s:%s:%s", resourceOwner, userID, projectID, projectGrantID))
}

type UserGrantAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID         string   `json:"userId,omitempty"`
	ProjectID      string   `json:"projectId,omitempty"`
	ProjectGrantID string   `json:"grantId,omitempty"`
	RoleKeys       []string `json:"roleKeys,omitempty"`
}

func (e *UserGrantAddedEvent) Payload() interface{} {
	return e
}

func (e *UserGrantAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUserGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.UserID, e.ProjectID, e.ProjectGrantID)}
}

func NewUserGrantAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	projectID,
	projectGrantID string,
	roleKeys []string) *UserGrantAddedEvent {
	return &UserGrantAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantAddedType,
		),
		UserID:         userID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
		RoleKeys:       roleKeys,
	}
}

func UserGrantAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserGrantAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "UGRANT-0p9ol", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys"`
}

func (e *UserGrantChangedEvent) Payload() interface{} {
	return e
}

func (e *UserGrantChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserGrantChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	roleKeys []string) *UserGrantChangedEvent {
	return &UserGrantChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func UserGrantChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "UGRANT-4M0sd", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantCascadeChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys,omitempty"`
}

func (e *UserGrantCascadeChangedEvent) Payload() interface{} {
	return e
}

func (e *UserGrantCascadeChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserGrantCascadeChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	roleKeys []string) *UserGrantCascadeChangedEvent {
	return &UserGrantCascadeChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantCascadeChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func UserGrantCascadeChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "UGRANT-Gs9df", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	userID               string `json:"-"`
	projectID            string `json:"-"`
	projectGrantID       string `json:"-"`
}

func (e *UserGrantRemovedEvent) Payload() interface{} {
	return nil
}

func (e *UserGrantRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.userID, e.projectID, e.projectGrantID)}
}

func NewUserGrantRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	projectID,
	projectGrantID string,
) *UserGrantRemovedEvent {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantRemovedType,
		),
		userID:         userID,
		projectID:      projectID,
		projectGrantID: projectGrantID,
	}
}

func UserGrantRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	userID               string `json:"-"`
	projectID            string `json:"-"`
	projectGrantID       string `json:"-"`
}

func (e *UserGrantCascadeRemovedEvent) Payload() interface{} {
	return nil
}

func (e *UserGrantCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.userID, e.projectID, e.projectGrantID)}
}

func NewUserGrantCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	projectID,
	projectGrantID string,
) *UserGrantCascadeRemovedEvent {
	return &UserGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantCascadeRemovedType,
		),
		userID:         userID,
		projectID:      projectID,
		projectGrantID: projectGrantID,
	}
}

func UserGrantCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *UserGrantDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserGrantDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserGrantDeactivatedEvent {
	return &UserGrantDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantDeactivatedType,
		),
	}
}

func UserGrantDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserGrantDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *UserGrantReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserGrantReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserGrantReactivatedEvent {
	return &UserGrantReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGrantReactivatedType,
		),
	}
}

func UserGrantReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserGrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
