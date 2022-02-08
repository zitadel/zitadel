package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	personalAccessTokenEventPrefix = userEventTypePrefix + "pat."
	PersonalAccessTokenAddedType   = personalAccessTokenEventPrefix + "added"
	PersonalAccessTokenRemovedType = personalAccessTokenEventPrefix + "removed"
)

type PersonalAccessTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID    string    `json:"tokenId"`
	Expiration time.Time `json:"expiration"`
	Scopes     []string  `json:"scopes"`
}

func (e *PersonalAccessTokenAddedEvent) Data() interface{} {
	return e
}

func (e *PersonalAccessTokenAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPersonalAccessTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
	expiration time.Time,
	scopes []string,
) *PersonalAccessTokenAddedEvent {
	return &PersonalAccessTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PersonalAccessTokenAddedType,
		),
		TokenID:    tokenID,
		Expiration: expiration,
		Scopes:     scopes,
	}
}

func PersonalAccessTokenAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	tokenAdded := &PersonalAccessTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Dbges", "unable to unmarshal token added")
	}

	return tokenAdded, nil
}

type PersonalAccessTokenRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId"`
}

func (e *PersonalAccessTokenRemovedEvent) Data() interface{} {
	return e
}

func (e *PersonalAccessTokenRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPersonalAccessTokenRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *PersonalAccessTokenRemovedEvent {
	return &PersonalAccessTokenRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PersonalAccessTokenRemovedType,
		),
		TokenID: tokenID,
	}
}

func PersonalAccessTokenRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	tokenRemoved := &PersonalAccessTokenRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Dbneg", "unable to unmarshal token removed")
	}

	return tokenRemoved, nil
}
