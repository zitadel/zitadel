package limits

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("limits.")
	SetEventType    = eventTypePrefix + "set"
	ResetEventType  = eventTypePrefix + "reset"
)

// SetEvent describes that limits are added or modified and contains only changed properties
type SetEvent struct {
	*eventstore.BaseEvent      `json:"-"`
	AuditLogRetention          *time.Duration `json:"auditLogRetention,omitempty"`
	AllowPublicOrgRegistration *bool          `json:"allowPublicOrgRegistration,omitempty"`
}

func (e *SetEvent) Payload() any {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewSetEvent(
	base *eventstore.BaseEvent,
	changes ...LimitsChange,
) *SetEvent {
	changedEvent := &SetEvent{
		BaseEvent: base,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent
}

type LimitsChange func(*SetEvent)

func ChangeAuditLogRetention(auditLogRetention *time.Duration) LimitsChange {
	return func(e *SetEvent) {
		e.AuditLogRetention = auditLogRetention
	}
}

func ChangeAllowPublicOrgRegistration(allow *bool) LimitsChange {
	return func(e *SetEvent) {
		e.AllowPublicOrgRegistration = allow
	}
}

var SetEventMapper = eventstore.GenericEventMapper[SetEvent]

type ResetProperty int

const (
	ResetUnknownProperty ResetProperty = iota
	ResetAuditLogRetention
	ResetAllowPublicOrgRegistration
)

type ResetEvent struct {
	*eventstore.BaseEvent `json:"-"`
	Properties            []ResetProperty `json:"properties,omitempty"`
}

func (e *ResetEvent) Payload() any {
	return e
}

func (e *ResetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *ResetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewResetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	properties ...ResetProperty,
) *ResetEvent {
	return &ResetEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ResetEventType,
		),
		Properties: properties,
	}
}

var ResetEventMapper = eventstore.GenericEventMapper[ResetEvent]
