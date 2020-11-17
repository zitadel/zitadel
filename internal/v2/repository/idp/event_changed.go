package idp

import "github.com/caos/zitadel/internal/eventstore/v2"

type ChangedEdvent struct {
	eventstore.BaseEvent `json:"-"`

	current *ConfigAggregate
	changed *ConfigAggregate

	Name string `json:"name"`
}

func ChangedEvent(
	base *eventstore.BaseEvent,
	current *ConfigAggregate,
	changed *ConfigAggregate,
) (*ChangedEdvent, error) {
	//TODO: who to handle chanes?

	return &ChangedEdvent{
		BaseEvent: *base,
		current:   current,
		changed:   changed,
	}, nil
}

func (e *ChangedEdvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEdvent) Data() interface{} {
	if e.current.Name != e.changed.Name {
		e.Name = e.changed.Name
	}
	return e
}
