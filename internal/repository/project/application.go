package project

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueAppNameType          = "appname"
	applicationEventTypePrefix = projectEventTypePrefix + "application."
	ApplicationAddedType       = applicationEventTypePrefix + "added"
	ApplicationChangedType     = applicationEventTypePrefix + "changed"
	ApplicationDeactivatedType = applicationEventTypePrefix + "deactivated"
	ApplicationReactivatedType = applicationEventTypePrefix + "reactivated"
	ApplicationRemovedType     = applicationEventTypePrefix + "removed"
)

func NewAddApplicationUniqueConstraint(name, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueAppNameType,
		fmt.Sprintf("%s:%s", name, projectID),
		"Errors.Project.App.AlreadyExists")
}

func NewRemoveApplicationUniqueConstraint(name, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueAppNameType,
		fmt.Sprintf("%s:%s", name, projectID))
}

type ApplicationAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId,omitempty"`
	Name  string `json:"name,omitempty"`
}

func (e *ApplicationAddedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddApplicationUniqueConstraint(e.Name, e.Aggregate().ID)}
}

func NewApplicationAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	name string,
) *ApplicationAddedEvent {
	return &ApplicationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationAddedType,
		),
		AppID: appID,
		Name:  name,
	}
}

func ApplicationAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-Nffg2", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID   string `json:"appId,omitempty"`
	Name    string `json:"name,omitempty"`
	oldName string
}

func (e *ApplicationChangedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveApplicationUniqueConstraint(e.oldName, e.Aggregate().ID),
		NewAddApplicationUniqueConstraint(e.Name, e.Aggregate().ID),
	}
}

func NewApplicationChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	oldName,
	newName string,
) *ApplicationChangedEvent {
	return &ApplicationChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationChangedType,
		),
		AppID:   appID,
		Name:    newName,
		oldName: oldName,
	}
}

func ApplicationChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-9l0cs", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewApplicationDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
) *ApplicationDeactivatedEvent {
	return &ApplicationDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationDeactivatedType,
		),
		AppID: appID,
	}
}

func ApplicationDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-0p9fB", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationReactivatedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewApplicationReactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
) *ApplicationReactivatedEvent {
	return &ApplicationReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationReactivatedType,
		),
		AppID: appID,
	}
}

func ApplicationReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID    string `json:"appId,omitempty"`
	name     string
	entityID string
}

func (e *ApplicationRemovedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	remove := []*eventstore.UniqueConstraint{NewRemoveApplicationUniqueConstraint(e.name, e.Aggregate().ID)}
	if e.entityID != "" {
		remove = append(remove, NewRemoveSAMLConfigEntityIDUniqueConstraint(e.entityID))
	}
	return remove
}

func NewApplicationRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	name string,
	entityID string,
) *ApplicationRemovedEvent {
	return &ApplicationRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationRemovedType,
		),
		AppID:    appID,
		name:     name,
		entityID: entityID,
	}
}

func ApplicationRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}
