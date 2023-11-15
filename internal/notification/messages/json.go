package messages

import (
	"encoding/json"

	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/notification/channels"
)

var _ channels.Message = (*JSON)(nil)

type JSON struct {
	Serializable    interface{}
	TriggeringEvent eventstore.Event
}

func (msg *JSON) GetContent() (string, error) {
	bytes, err := json.Marshal(msg.Serializable)
	return string(bytes), err
}

func (msg *JSON) GetTriggeringEvent() eventstore.Event {
	return msg.TriggeringEvent
}
