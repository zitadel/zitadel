package org

import "github.com/zitadel/zitadel/internal/eventstore"

var (
	_           eventstore.Event = (*deactivated)(nil)
	Deactivated *deactivated
)

type deactivated struct {
	eventstore.BaseEvent
}

func NewDeactivated() *deactivated {
	return new(deactivated)
}

func (*deactivated) Type() eventstore.EventType {
	return "org.deactivated"
}
