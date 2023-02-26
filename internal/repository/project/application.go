package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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

	AppID string `json:"appId,omitempty"`
	Name  string `json:"name,omitempty"`
	ExternalURL string `json:"external_url,omitempty"`
	IsVisibleToEndUser bool `json:"is_visible_to_end_user,omitempty"`
}

func (e *ApplicationAddedEvent) Data() interface{} {
	return e
}

func (e *ApplicationAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddApplicationUniqueConstraint(e.Name, e.Aggregate().ID)}
}

func NewApplicationAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	name, externalURL string,
	isVisibleToEndUser bool,
) *ApplicationAddedEvent {
	return &ApplicationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationAddedType,
		),
		AppID: appID,
		Name:  name,
		ExternalURL: externalURL,
		IsVisibleToEndUser: isVisibleToEndUser,
	}
}

func ApplicationAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
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

	AppID   string `json:"appId,omitempty"`
	Name    string `json:"name,omitempty"`
	ExternalURL string `json:"external_url,omitempty"`
	IsVisibleToEndUser bool `json:"is_visible_to_end_user,omitempty"`
	oldName string
}

func (e *ApplicationChangedEvent) Data() interface{} {
	return e
}

func (e *ApplicationChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{
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
	externalURL string,	
	isVisibleToEndUser bool,
) *ApplicationChangedEvent {
	return &ApplicationChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationChangedType,
		),
		AppID:   appID,
		Name:    newName,
		ExternalURL: externalURL,
		IsVisibleToEndUser: isVisibleToEndUser,
		oldName: oldName,
	}
}

func ApplicationChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
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

func ApplicationDeactivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
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

func ApplicationReactivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
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

	AppID    string `json:"appId,omitempty"`
	name     string
	entityID string
}

func (e *ApplicationRemovedEvent) Data() interface{} {
	return e
}

func (e *ApplicationRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	remove := []*eventstore.EventUniqueConstraint{NewRemoveApplicationUniqueConstraint(e.name, e.Aggregate().ID)}
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

func ApplicationRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &ApplicationRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "APPLICATION-1m9e3", "unable to unmarshal application")
	}

	return e, nil
}
