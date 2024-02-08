package org

import "github.com/zitadel/zitadel/internal/eventstore"

var (
	_       eventstore.Event = (*changed)(nil)
	Changed *changed
)

type changed struct {
	eventstore.BaseEvent
}

func NewChanged() *changed {
	return new(changed)
}

func (*changed) Type() eventstore.EventType {
	return "org.changed"
}
