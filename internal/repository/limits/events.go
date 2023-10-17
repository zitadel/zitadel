package limits

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	eventTypePrefix = eventstore.EventType("limits.")
	SetEventType    = eventTypePrefix + "set"
	ResetEventType  = eventTypePrefix + "reset"
)

// SetEvent describes that limits are added or modified and contains only changed properties
type SetEvent struct {
	eventstore.BaseEvent `json:"-"`
	AuditLogRetention    *time.Duration `json:"auditLogRetention,omitempty"`
}

func (e *SetEvent) Data() interface{} {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewSetEvent(
	base *eventstore.BaseEvent,
	changes ...LimitsChange,
) *SetEvent {
	changedEvent := &SetEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent
}

type LimitsChange func(*SetEvent)

func ChangeAuditLogRetention(auditLogRetention time.Duration) LimitsChange {
	return func(e *SetEvent) {
		e.AuditLogRetention = &auditLogRetention
	}
}

func SetEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "LIMITS-YwPkZ", "unable to unmarshal limits set")
	}
	return e, nil
}

type ResetEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ResetEvent) Data() interface{} {
	return e
}

func (e *ResetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewResetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *ResetEvent {
	return &ResetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ResetEventType,
		),
	}
}

func ResetEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &ResetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "LIMITS-a9O4Q", "unable to unmarshal limits reset")
	}
	return e, nil
}
