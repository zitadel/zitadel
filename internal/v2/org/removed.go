package org

import "github.com/zitadel/zitadel/internal/eventstore"

var (
	_       eventstore.Event = (*removed)(nil)
	Removed *removed
)

type removed struct {
	eventstore.BaseEvent
}

func NewRemoved() *removed {
	return new(removed)
}

func (*removed) Type() eventstore.EventType {
	return "org.removed"
}
