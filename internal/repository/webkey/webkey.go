package webkey

import (
	"context"
	"encoding/json"

	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueWebKeyType = "web_key"
)

const (
	eventTypePrefix      = eventstore.EventType("web_key.")
	AddedEventType       = eventTypePrefix + "added"
	ActivatedEventType   = eventTypePrefix + "activated"
	DeactivatedEventType = eventTypePrefix + "deactivated"
	RemovedEventType     = eventTypePrefix + "removed"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	PrivateKey *crypto.CryptoValue     `json:"privateKey"`
	PublicKey  *jose.JSONWebKey        `json:"publicKey"`
	Config     json.RawMessage         `json:"config"`
	ConfigType crypto.WebKeyConfigType `json:"configType"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewAddEventUniqueConstraint(UniqueWebKeyType, e.Agg.ID, "Errors.WebKey.Duplicate"),
	}
}

func (e *AddedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = base
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	privateKey *crypto.CryptoValue,
	publicKey *jose.JSONWebKey,
	config crypto.WebKeyConfig,
) (*AddedEvent, error) {
	configJson, err := json.Marshal(config)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "WEBKEY-IY9fa", "Errors.Internal")
	}
	return &AddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedEventType,
		),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Config:     configJson,
		ConfigType: config.Type(),
	}, nil
}

type ActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *ActivatedEvent) Payload() interface{} {
	return e
}

func (e *ActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *ActivatedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = base
}

func NewActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *ActivatedEvent {
	return &ActivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ActivatedEventType,
		),
	}
}

type DeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *DeactivatedEvent) Payload() interface{} {
	return e
}

func (e *DeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *DeactivatedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = base
}

func NewDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DeactivatedEvent {
	return &DeactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DeactivatedEventType,
		),
	}
}

type RemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) Payload() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		eventstore.NewRemoveUniqueConstraint(UniqueWebKeyType, e.Agg.ID),
	}
}

func (e *RemovedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = base
}

func NewRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RemovedEventType,
		),
	}
}
