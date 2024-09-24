package authenticator

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	jwtPrefix      = eventPrefix + "jwt."
	JWTCreatedType = jwtPrefix + "created"
	JWTDeletedType = jwtPrefix + "deleted"
)

type JWTCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID string `json:"userID"`

	ExpirationDate    time.Time `json:"expirationDate,omitempty"`
	PublicKey         []byte    `json:"publicKey,omitempty"`
	TriggeredAtOrigin string    `json:"triggerOrigin,omitempty"`
}

func (e *JWTCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *JWTCreatedEvent) Payload() interface{} {
	return e
}

func (e *JWTCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *JWTCreatedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewJWTCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	expirationDate time.Time,
	publicKey []byte,
) *JWTCreatedEvent {
	return &JWTCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			JWTCreatedType,
		),
		UserID:            userID,
		ExpirationDate:    expirationDate,
		PublicKey:         publicKey,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type JWTDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *JWTDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *JWTDeletedEvent) Payload() interface{} {
	return e
}

func (e *JWTDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewJWTDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *JWTDeletedEvent {
	return &JWTDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			JWTDeletedType,
		),
	}
}
