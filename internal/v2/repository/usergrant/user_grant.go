package usergrant

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	uniqueUserGrant             = "user_grant"
	userGrantEventTypePrefix    = eventstore.EventType("user.grant")
	UserGrantAddedType          = userGrantEventTypePrefix + "added"
	UserGrantChangedType        = userGrantEventTypePrefix + "changed"
	UserGrantCascadeChangedType = userGrantEventTypePrefix + "cascade.changed"
	UserGrantRemovedType        = userGrantEventTypePrefix + "removed"
	UserGrantCascadeRemovedType = userGrantEventTypePrefix + "cascade.removed"
	UserGrantDeactivatedType    = userGrantEventTypePrefix + "deactivated"
	UserGrantReactivatedType    = userGrantEventTypePrefix + "reactivated"
)

func NewAddUserGrantUniqueConstraint(resourceOwner, userID, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueUserGrant,
		fmt.Sprintf("%s:%s:%s", resourceOwner, userID, projectID),
		"Errors.UserGrant.AlreadyExists")
}

func NewRemoveUserGrantUniqueConstraint(resourceOwner, userID, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueUserGrant,
		fmt.Sprintf("%s:%s:%s", resourceOwner, userID, projectID))
}

type UserGrantAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID         string   `json:"userId,omitempty"`
	ProjectID      string   `json:"projectId,omitempty"`
	ProjectGrantID string   `json:"grantId,omitempty"`
	RoleKeys       []string `json:"roleKeys,omitempty"`
}

func (e *UserGrantAddedEvent) Data() interface{} {
	return e
}

func (e *UserGrantAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddUserGrantUniqueConstraint(e.ResourceOwner(), e.UserID, e.ProjectID)}
}

func NewUserGrantAddedEvent(
	ctx context.Context,
	resourceOwner,
	userID,
	projectID,
	projectGrantID string,
	roleKeys []string) *UserGrantAddedEvent {
	return &UserGrantAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			UserGrantAddedType,
			resourceOwner,
		),
		UserID:         userID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
		RoleKeys:       roleKeys,
	}
}

func UserGrantAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &UserGrantAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-0p9ol", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys,omitempty"`
}

func (e *UserGrantChangedEvent) Data() interface{} {
	return e
}

func (e *UserGrantChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewUserGrantChangedEvent(
	ctx context.Context,
	roleKeys []string) *UserGrantChangedEvent {
	return &UserGrantChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func UserGrantChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &UserGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-4M0sd", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantCascadeChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys,omitempty"`
}

func (e *UserGrantCascadeChangedEvent) Data() interface{} {
	return e
}

func (e *UserGrantCascadeChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewUserGrantCascadeChangedEvent(
	ctx context.Context,
	roleKeys []string) *UserGrantCascadeChangedEvent {
	return &UserGrantCascadeChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantCascadeChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func UserGrantCascadeChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &UserGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-Gs9df", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	userID               string
	projectID            string
}

func (e *UserGrantRemovedEvent) Data() interface{} {
	return nil
}

func (e *UserGrantRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveUserGrantUniqueConstraint(e.ResourceOwner(), e.userID, e.projectID)}
}

func NewUserGrantRemovedEvent(ctx context.Context, resourceOwner, userID, projectID string) *UserGrantRemovedEvent {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			UserGrantRemovedType,
			resourceOwner,
		),
		userID:    userID,
		projectID: projectID,
	}
}

func UserGrantRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	userID               string
	projectID            string
}

func (e *UserGrantCascadeRemovedEvent) Data() interface{} {
	return nil
}

func (e *UserGrantCascadeRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveUserGrantUniqueConstraint(e.ResourceOwner(), e.userID, e.projectID)}
}

func NewUserGrantCascadeRemovedEvent(ctx context.Context, resourceOwner, userID, projectID string) *UserGrantCascadeRemovedEvent {
	return &UserGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			UserGrantCascadeRemovedType,
			resourceOwner,
		),
		userID:    userID,
		projectID: projectID,
	}
}

func UserGrantCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantDeactivatedEvent) Data() interface{} {
	return nil
}

func (e *UserGrantDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewUserGrantDeactivatedEvent(ctx context.Context) *UserGrantDeactivatedEvent {
	return &UserGrantDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantDeactivatedType,
		),
	}
}

func UserGrantDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserGrantDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserGrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantReactivatedEvent) Data() interface{} {
	return nil
}

func (e *UserGrantReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewUserGrantReactivatedEvent(ctx context.Context) *UserGrantReactivatedEvent {
	return &UserGrantReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantReactivatedType,
		),
	}
}

func UserGrantReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserGrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
