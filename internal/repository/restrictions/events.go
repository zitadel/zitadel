package restrictions

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("restrictions.")
	SetEventType    = eventTypePrefix + "set"
)

// SetEvent describes that restrictions are added or modified and contains only changed properties
type SetEvent struct {
	*eventstore.BaseEvent          `json:"-"`
	DisallowPublicOrgRegistrations *bool `json:"disallowPublicOrgRegistrations,omitempty"`
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
	changes ...RestrictionsChange,
) *SetEvent {
	changedEvent := &SetEvent{
		BaseEvent: base,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent
}

type RestrictionsChange func(*SetEvent)

func ChangePublicOrgRegistrations(disallow bool) RestrictionsChange {
	return func(e *SetEvent) {
		e.DisallowPublicOrgRegistrations = gu.Ptr(disallow)
	}
}

var SetEventMapper = eventstore.GenericEventMapper[SetEvent]
