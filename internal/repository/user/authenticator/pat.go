package authenticator

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	patPrefix      = eventPrefix + "pat."
	PATCreatedType = patPrefix + "created"
	PATDeletedType = patPrefix + "deleted"
)

type PATCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID string `json:"userID"`

	ExpirationDate    time.Time `json:"expirationDate,omitempty"`
	Scopes            []string  `json:"scopes,omitempty"`
	TriggeredAtOrigin string    `json:"triggerOrigin,omitempty"`
}

func (e *PATCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PATCreatedEvent) Payload() interface{} {
	return e
}

func (e *PATCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPATCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	expirationDate time.Time,
	scopes []string,
) *PATCreatedEvent {
	return &PATCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PATCreatedType,
		),
		UserID:            userID,
		ExpirationDate:    expirationDate,
		Scopes:            scopes,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type PATDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *PATDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PATDeletedEvent) Payload() interface{} {
	return e
}

func (e *PATDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPATDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PATDeletedEvent {
	return &PATDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PATDeletedType,
		),
	}
}
