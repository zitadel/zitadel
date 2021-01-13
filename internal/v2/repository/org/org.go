package org

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	OrgAdded       = orgEventTypePrefix + "added"
	OrgChanged     = orgEventTypePrefix + "changed"
	OrgDeactivated = orgEventTypePrefix + "deactivated"
	OrgReactivated = orgEventTypePrefix + "reactivated"
	OrgRemoved     = orgEventTypePrefix + "removed"
)

type OrgAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *OrgAddedEvent) Data() interface{} {
	return e
}

func NewOrgAddedEvent(ctx context.Context, name string) *OrgAddedEvent {
	return &OrgAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgAdded,
		),
		Name: name,
	}
}

func OrgAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgAdded := &OrgAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgAdded, nil
}

type OrgChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *OrgChangedEvent) Data() interface{} {
	return e
}

func NewOrgChangedEvent(ctx context.Context, name string) *OrgChangedEvent {
	return &OrgChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgChanged,
		),
		Name: name,
	}
}

func OrgChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgChanged, nil
}

type OrgDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgDeactivatedEvent) Data() interface{} {
	return e
}

func NewOrgDeactivatedEvent(ctx context.Context) *OrgDeactivatedEvent {
	return &OrgDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDeactivated,
		),
	}
}

func OrgDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-DAfbs", "unable to unmarshal org deactivated")
	}

	return orgChanged, nil
}

type OrgReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgReactivatedEvent) Data() interface{} {
	return e
}

func NewOrgReactivatedEvent(ctx context.Context) *OrgReactivatedEvent {
	return &OrgReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgReactivated,
		),
	}
}

func OrgReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-DAfbs", "unable to unmarshal org deactivated")
	}

	return orgChanged, nil
}
