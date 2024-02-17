package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	_ eventstore.Command = (*ChangedEvent)(nil)
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Changed *ChangedEvent
)

type ChangedEvent struct {
	creator string
}

func NewChangedEvent() *ChangedEvent {
	return new(ChangedEvent)
}

// Creator implements [eventstore.action].
func (e *ChangedEvent) Creator() string {
	return e.creator
}

// Payload implements [eventstore.Command].
func (*ChangedEvent) Payload() any {
	panic("unimplemented")
}

// Revision implements [eventstore.action].
func (*ChangedEvent) Revision() uint16 {
	return 1
}

// UniqueConstraints implements [eventstore.Command].
func (*ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	panic("unimplemented")
}

// UniqueConstraints implements [eventstore.action].
func (*ChangedEvent) Type() string {
	return "org.changed"
}
