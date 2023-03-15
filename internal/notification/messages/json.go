package messages

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/notification/channels"
)

var _ channels.Message = (*JSON)(nil)

type JSON struct {
	Serializable interface{}
}

func (msg *JSON) GetContent() (string, error) {
	bytes, err := json.Marshal(msg.Serializable)
	return string(bytes), err
}
