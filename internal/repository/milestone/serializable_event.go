package milestone

import "github.com/zitadel/zitadel/internal/eventstore"

type SerializableEvent struct {
	eventstore.BaseEvent `json:",inline"`
	Data                 []byte `json:"data"`
}

func newSerializableEvent(triggeringEvent eventstore.BaseEvent) SerializableEvent {
	return SerializableEvent{
		BaseEvent: triggeringEvent,
		Data:      triggeringEvent.DataAsBytes(),
	}
}
