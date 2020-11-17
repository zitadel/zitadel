package idp

import "github.com/caos/zitadel/internal/eventstore/v2"

type AddedEvent struct {
	eventstore.BaseEvent

	ID   string `idpConfigId`
	Name string `name`
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent: *base,
		ID:        configID,
		Name:      name,
	}
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}
