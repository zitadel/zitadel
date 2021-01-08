package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	applicationEventTypePrefix = projectEventTypePrefix + "application."
	ApplicationAdded           = applicationEventTypePrefix + "added"
	ApplicationChanged         = applicationEventTypePrefix + "changed"
	ApplicationDeactivated     = applicationEventTypePrefix + "deactivated"
	ApplicationReactivated     = applicationEventTypePrefix + "reactivated"
	ApplicationRemoved         = applicationEventTypePrefix + "removed"
)

type ApplicationAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *ApplicationAddedEvent) Data() interface{} {
	return e
}

func NewApplicationAddedEvent(ctx context.Context, name string) *ApplicationAddedEvent {
	return &ApplicationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ApplicationAdded,
		),
		Name: name,
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
