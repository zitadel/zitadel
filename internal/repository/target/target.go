package target

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	keyEventTypePrefix    eventstore.EventType = "target.key."
	KeyAddedEventType                          = keyEventTypePrefix + "added"
	KeyActivatedEventType                      = keyEventTypePrefix + "activated"
	KeyRemovedEventType                        = keyEventTypePrefix + "removed"
)

type KeyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID     string `json:"keyId,omitempty"`
	PublicKey []byte `json:"publicKey,omitempty"`
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
) *KeyAddedEvent {
	return &KeyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx, aggregate, KeyAddedEventType,
		),
		KeyID:     keyID,
		PublicKey: publicKey,
	}
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
