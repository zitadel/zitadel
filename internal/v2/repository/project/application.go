package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
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

	AppID   string         `json:"appId,omitempty"`
	Name    string         `json:"name,omitempty"`
	AppType domain.AppType `json:"appType,omitempty"`
}

func (e *ApplicationAddedEvent) Data() interface{} {
	return e
}

func (e *ApplicationAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewApplicationAddedEvent(ctx context.Context, appID, name string, appType domain.AppType) *ApplicationAddedEvent {
	return &ApplicationAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ApplicationAdded,
		),
		AppID:   appID,
		Name:    name,
		AppType: appType,
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
