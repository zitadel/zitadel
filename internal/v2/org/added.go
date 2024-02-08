package org

import "github.com/zitadel/zitadel/internal/eventstore"

var (
	_     eventstore.Event = (*added)(nil)
	Added *added
)

type added struct {
	eventstore.BaseEvent

	Name string
}

func NewAdded() *added {
	return new(added)
}

func (*added) Type() eventstore.EventType {
	return "org.added"
}
