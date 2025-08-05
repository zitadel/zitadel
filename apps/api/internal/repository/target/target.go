package target

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix  eventstore.EventType = "target."
	AddedEventType                        = eventTypePrefix + "added"
	ChangedEventType                      = eventTypePrefix + "changed"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name             string              `json:"name"`
	TargetType       domain.TargetType   `json:"targetType"`
	Endpoint         string              `json:"endpoint"`
	Timeout          time.Duration       `json:"timeout"`
	InterruptOnError bool                `json:"interruptOnError"`
	SigningKey       *crypto.CryptoValue `json:"signingKey"`
}

func (e *AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *AddedEvent) Payload() any {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(e.Name)}
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	targetType domain.TargetType,
	endpoint string,
	timeout time.Duration,
	interruptOnError bool,
	signingKey *crypto.CryptoValue,
) *AddedEvent {
	return &AddedEvent{
		*eventstore.NewBaseEventForPush(
			ctx, aggregate, AddedEventType,
		),
		name, targetType, endpoint, timeout, interruptOnError, signingKey}
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name             *string             `json:"name,omitempty"`
	TargetType       *domain.TargetType  `json:"targetType,omitempty"`
	Endpoint         *string             `json:"endpoint,omitempty"`
	Timeout          *time.Duration      `json:"timeout,omitempty"`
	InterruptOnError *bool               `json:"interruptOnError,omitempty"`
	SigningKey       *crypto.CryptoValue `json:"signingKey,omitempty"`

	oldName string
}

func (e *ChangedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldName == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveUniqueConstraint(e.oldName),
		NewAddUniqueConstraint(*e.Name),
	}
}

func NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []Changes,
) *ChangedEvent {
	changeEvent := &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent
}

type Changes func(event *ChangedEvent)

func ChangeName(oldName, name string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Name = &name
		e.oldName = oldName
	}
}

func ChangeTargetType(targetType domain.TargetType) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.TargetType = &targetType
	}
}

func ChangeEndpoint(endpoint string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Endpoint = &endpoint
	}
}

func ChangeTimeout(timeout time.Duration) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Timeout = &timeout
	}
}

func ChangeInterruptOnError(interruptOnError bool) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.InterruptOnError = &interruptOnError
	}
}

func ChangeSigningKey(signingKey *crypto.CryptoValue) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.SigningKey = signingKey
	}
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	name string
}

func (e *RemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *RemovedEvent) Payload() any {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUniqueConstraint(e.name)}
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *RemovedEvent {
	return &RemovedEvent{*eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType), name}
}
