package authenticator

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	passwordPrefix      = eventPrefix + "password."
	PasswordCreatedType = passwordPrefix + "created"
	PasswordDeletedType = passwordPrefix + "deleted"
)

type PasswordCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID         string `json:"userID"`
	EncodedHash    string `json:"encodedHash,omitempty"`
	ChangeRequired bool   `json:"changeRequired,omitempty"`
	TriggerOrigin  string `json:"triggerOrigin,omitempty"`
}

func (e *PasswordCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PasswordCreatedEvent) Payload() interface{} {
	return e
}

func (e *PasswordCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	encodeHash string,
	changeRequired bool,
) *PasswordCreatedEvent {
	return &PasswordCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordCreatedType,
		),
		UserID:         userID,
		EncodedHash:    encodeHash,
		ChangeRequired: changeRequired,
		TriggerOrigin:  http.DomainContext(ctx).Origin(),
	}
}

type PasswordDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *PasswordDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PasswordDeletedEvent) Payload() interface{} {
	return e
}

func (e *PasswordDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PasswordDeletedEvent {
	return &PasswordDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordDeletedType,
		),
	}
}
