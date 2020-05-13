package chat

import (
	"bytes"
	"encoding/json"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	"net/http"
	"net/url"
)

type Chat struct {
	URL *url.URL
}

func InitChatProvider(config *ChatConfig) (*Chat, error) {
	url, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}
	return &Chat{
		URL: url,
	}, nil
}

func (chat *Chat) CanHandleMessage(message providers.Message) bool {
	chatMsg := message.(ChatMessage)
	return chatMsg.Content != ""
}

func (chat *Chat) HandleMessage(message providers.Message) error {
	req, err := json.Marshal(message.GetContent())
	if err != nil {
		return caos_errs.ThrowInternal(err, "PROVI-s8uie", "Could not unmarshal content")
	}

	_, err = http.Post(chat.URL.String(), "application/json; charset=UTF-8", bytes.NewReader(req))
	if err != nil {
		return caos_errs.ThrowInternal(err, "PROVI-si93s", "unable to send message")
	}
	return nil
}
