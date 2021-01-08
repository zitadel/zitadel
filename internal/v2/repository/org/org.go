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
