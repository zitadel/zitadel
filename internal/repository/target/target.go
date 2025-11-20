package target

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	keyEventTypePrefix      eventstore.EventType = "target.key."
	KeyAddedEventType                            = keyEventTypePrefix + "added"
	KeyActivatedEventType                        = keyEventTypePrefix + "activated"
	KeyDeactivatedEventType                      = keyEventTypePrefix + "deactivated"
	KeyRemovedEventType                          = keyEventTypePrefix + "removed"
)

type KeyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID          string    `json:"keyId,omitempty"`
	PublicKey      []byte    `json:"publicKey,omitempty"`
	Fingerprint    string    `json:"fingerprint,omitempty"`
	ExpirationDate time.Time `json:"expirationDate,omitempty"`
}

func (e *KeyAddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *KeyAddedEvent) Payload() any {
	return e
}

func (e *KeyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewKeyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	keyID string,
	publicKey []byte,
	fingerprint string,
	expiration time.Time,
) *KeyAddedEvent {
	return &KeyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx, aggregate, KeyAddedEventType,
		),
		KeyID:          keyID,
		PublicKey:      publicKey,
		Fingerprint:    fingerprint,
		ExpirationDate: expiration,
	}
}

type KeyActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *KeyActivatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *KeyActivatedEvent) Payload() any {
	return e
}

func (e *KeyActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewKeyActivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string) *KeyActivatedEvent {
	return &KeyActivatedEvent{*eventstore.NewBaseEventForPush(ctx, aggregate, KeyActivatedEventType), id}
}

type KeyDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *KeyDeactivatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *KeyDeactivatedEvent) Payload() any {
	return e
}

func (e *KeyDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewKeyDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string) *KeyDeactivatedEvent {
	return &KeyDeactivatedEvent{*eventstore.NewBaseEventForPush(ctx, aggregate, KeyDeactivatedEventType), id}
}

type KeyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *KeyRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *KeyRemovedEvent) Payload() any {
	return e
}

func (e *KeyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewKeyRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string) *KeyRemovedEvent {
	return &KeyRemovedEvent{*eventstore.NewBaseEventForPush(ctx, aggregate, KeyRemovedEventType), id}
}
