package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	_ eventstore.Command = (*DeactivatedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Deactivated *DeactivatedEvent
)

type DeactivatedEvent struct {
	creator string
}

func NewDeactivatedEvent() *DeactivatedEvent {
	return new(DeactivatedEvent)
}

// Creator implements [eventstore.action].
func (e *DeactivatedEvent) Creator() string {
	return e.creator
}

// Payload implements [eventstore.Command].
func (*DeactivatedEvent) Payload() any {
	return nil
}

// Revision implements [eventstore.action].
func (*DeactivatedEvent) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (*DeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	panic("unimplemented")
}

// Type implements [eventstore.action].
func (*DeactivatedEvent) Type() string {
	return "org.deactivated"
}
