package restrictions

import (
	"github.com/muhlemmer/gu"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("restrictions.")
	SetEventType    = eventTypePrefix + "set"
)

// SetEvent describes that restrictions are added or modified and contains only changed properties
type SetEvent struct {
	*eventstore.BaseEvent         `json:"-"`
	DisallowPublicOrgRegistration *bool           `json:"disallowPublicOrgRegistration,omitempty"`
	AllowedLanguages              *[]language.Tag `json:"allowedLanguages,omitempty"`
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

func ChangeDisallowPublicOrgRegistration(disallow bool) RestrictionsChange {
	return func(e *SetEvent) {
		e.DisallowPublicOrgRegistration = gu.Ptr(disallow)
	}
}

func ChangeAllowedLanguages(allowedLanguages []language.Tag) RestrictionsChange {
	return func(e *SetEvent) {
		e.AllowedLanguages = &allowedLanguages
	}
}

var SetEventMapper = eventstore.GenericEventMapper[SetEvent]
