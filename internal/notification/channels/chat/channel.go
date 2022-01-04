package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels"
)

var _ channels.NotificationChannel = (*Chat)(nil)

type Chat struct {
	URL        *url.URL
	SplitCount int
}

func InitChatProvider(config ChatConfig) (*Chat, error) {

	url, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}
	return &Chat{
		URL:        url,
		SplitCount: config.SplitCount,
	}, nil
}

func (chat *Chat) HandleMessage(message channels.Message) error {
	contentText := message.GetContent()
	for _, splittedMsg := range splitMessage(contentText, chat.SplitCount) {
		if err := chat.sendMessage(splittedMsg); err != nil {
			return err
		}
	}
	return nil
}

func (chat *Chat) sendMessage(message string) error {
	req, err := json.Marshal(message)
	if err != nil {
		return caos_errs.ThrowInternal(err, "PROVI-s8uie", "Could not unmarshal content")
	}

	response, err := http.Post(chat.URL.String(), "application/json; charset=UTF-8", bytes.NewReader(req))
	if err != nil {
		return caos_errs.ThrowInternal(err, "PROVI-si93s", "unable to send message")
	}
	if response.StatusCode != 200 {
		defer response.Body.Close()
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return caos_errs.ThrowInternal(err, "PROVI-PSLd3", "unable to read response message")
		}
		logging.LogWithFields("PROVI-PS0kx", "Body", string(bodyBytes)).Warn("Chat Message post didnt get 200 OK")
		return caos_errs.ThrowInternal(nil, "PROVI-LSopw", string(bodyBytes))
	}
	return nil
}

func splitMessage(message string, count int) []string {
	if count == 0 {
		return []string{message}
	}
	var splits []string
	var l, r int
	for l, r = 0, count; r < len(message); l, r = r, r+count {
		for !utf8.RuneStart(message[r]) {
			r--
		}
		splits = append(splits, message[l:r])
	}
	splits = append(splits, message[l:])
	return splits
}
