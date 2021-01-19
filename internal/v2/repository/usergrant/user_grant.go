package usergrant

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	userGrantEventTypePrefix    = eventstore.EventType("user.grant")
	UserGrantAddedType          = userGrantEventTypePrefix + "added"
	UserGrantChangedType        = userGrantEventTypePrefix + "changed"
	UserGrantCascadeChangedType = userGrantEventTypePrefix + "cascade.changed"
	UserGrantRemovedType        = userGrantEventTypePrefix + "removed"
	UserGrantCascadeRemovedType = userGrantEventTypePrefix + "cascade.removed"
	UserGrantDeactivatedType    = userGrantEventTypePrefix + "deactivated"
	UserGrantReactivatedType    = userGrantEventTypePrefix + "reactivated"
)

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

func NewUserGrantAddedEvent(
	ctx context.Context,
	userID,
	projectID,
	projectGrantID string,
	roleKeys []string) *UserGrantAddedEvent {
	return &UserGrantAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantAddedType,
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
		return nil, errors.ThrowInternal(err, "UGRANT-2M9fs", "unable to unmarshal user grant")
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
}

func (e *UserGrantRemovedEvent) Data() interface{} {
	return e
}

func NewUserGrantRemovedEvent(ctx context.Context) *UserGrantRemovedEvent {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantRemovedType,
		),
	}
}

func UserGrantRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &UserGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-M0sdf", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantCascadeRemovedEvent) Data() interface{} {
	return e
}

func NewUserGrantCascadeRemovedEvent(ctx context.Context) *UserGrantRemovedEvent {
	return &UserGrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserGrantRemovedType,
		),
	}
}

func UserGrantCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &UserGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-E7urs", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantDeactivatedEvent) Data() interface{} {
	return e
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
	e := &UserGrantDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-pL0ds", "unable to unmarshal user grant")
	}

	return e, nil
}

type UserGrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserGrantReactivatedEvent) Data() interface{} {
	return e
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
	e := &UserGrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "UGRANT-M0sdf", "unable to unmarshal user grant")
	}

	return e, nil
}
