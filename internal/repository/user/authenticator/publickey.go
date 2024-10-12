package authenticator

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	publicKeyPrefix      = eventPrefix + "public_key."
	PublicKeyCreatedType = publicKeyPrefix + "created"
	PublicKeyDeletedType = publicKeyPrefix + "deleted"
)

type PublicKeyCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID string `json:"userID"`

	ExpirationDate    time.Time `json:"expirationDate,omitempty"`
	PublicKey         []byte    `json:"publicKey,omitempty"`
	TriggeredAtOrigin string    `json:"triggerOrigin,omitempty"`
}

func (e *PublicKeyCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PublicKeyCreatedEvent) Payload() interface{} {
	return e
}

func (e *PublicKeyCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PublicKeyCreatedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewPublicKeyCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	expirationDate time.Time,
	publicKey []byte,
) *PublicKeyCreatedEvent {
	return &PublicKeyCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PublicKeyCreatedType,
		),
		UserID:            userID,
		ExpirationDate:    expirationDate,
		PublicKey:         publicKey,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type PublicKeyDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *PublicKeyDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PublicKeyDeletedEvent) Payload() interface{} {
	return e
}

func (e *PublicKeyDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPublicKeyDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PublicKeyDeletedEvent {
	return &PublicKeyDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PublicKeyDeletedType,
		),
	}
}
