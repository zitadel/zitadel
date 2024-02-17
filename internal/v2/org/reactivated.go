package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	_ eventstore.Command = (*ReactivatedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Reactivated *ReactivatedEvent
)

type ReactivatedEvent struct {
	creator string
}

func NewReactivatedEvent() *ReactivatedEvent {
	return new(ReactivatedEvent)
}

// Creator implements [eventstore.action].
func (e *ReactivatedEvent) Creator() string {
	return e.creator
}

// Payload implements [eventstore.Command].
func (*ReactivatedEvent) Payload() any {
	return nil
}

// Revision implements [eventstore.action].
func (*ReactivatedEvent) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (*ReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	panic("unimplemented")
}

// Type implements [eventstore.action].
func (*ReactivatedEvent) Type() string {
	return "org.reactivated"
}
