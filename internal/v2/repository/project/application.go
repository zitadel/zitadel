package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func NewAddApplicationUniqueConstraint(name, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueAppNameType,
		fmt.Sprintf("%s:%s", name, projectID),
		"Errors.Project.App.AlreadyExists")
}

func NewRemoveApplicationUniqueConstraint(name, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueAppNameType,
		fmt.Sprintf("%s:%s", name, projectID))
}

type ApplicationAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID     string `json:"appId,omitempty"`
	Name      string `json:"name,omitempty"`
	projectID string
}

func (e *ApplicationAddedEvent) Data() interface{} {
	return e
}

func (e *ApplicationAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddApplicationUniqueConstraint(e.Name, e.projectID)}
}

func NewApplicationAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	name,
	projectID string,
) *ApplicationAddedEvent {
	return &ApplicationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationAddedType,
		),
		AppID:     appID,
		Name:      name,
		projectID: projectID,
	}
}

func ApplicationAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ApplicationAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-Nffg2", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID     string `json:"appId,omitempty"`
	Name      string `json:"name,omitempty"`
	oldName   string
	projectID string
}

func (e *ApplicationChangedEvent) Data() interface{} {
	return e
}

func (e *ApplicationChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{
		NewRemoveApplicationUniqueConstraint(e.oldName, e.projectID),
		NewAddApplicationUniqueConstraint(e.Name, e.projectID),
	}
}

func NewApplicationChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	oldName,
	newName,
	projectID string,
) *ApplicationChangedEvent {
	return &ApplicationChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationChangedType,
		),
		AppID:     appID,
		Name:      newName,
		oldName:   oldName,
		projectID: projectID,
	}
}

func ApplicationChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ApplicationChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-9l0cs", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationDeactivatedEvent) Data() interface{} {
	return e
}

func (e *ApplicationDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func ApplicationDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ApplicationDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-0p9fB", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId,omitempty"`
}

func (e *ApplicationReactivatedEvent) Data() interface{} {
	return e
}

func (e *ApplicationReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func ApplicationReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ApplicationReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}

type ApplicationRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID     string `json:"appId,omitempty"`
	name      string
	projectID string
}

func (e *ApplicationRemovedEvent) Data() interface{} {
	return e
}

func (e *ApplicationRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveApplicationUniqueConstraint(e.name, e.projectID)}
}

func NewApplicationRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	name,
	projectID string,
) *ApplicationRemovedEvent {
	return &ApplicationRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationRemovedType,
		),
		AppID:     appID,
		name:      name,
		projectID: projectID,
	}
}

func ApplicationRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ApplicationRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}
