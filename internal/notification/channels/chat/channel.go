package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/k3a/html2text"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels"
)

func InitChatChannel(config ChatConfig) (channels.NotificationChannel, error) {

	url, err := url.Parse(config.Url)
	if err != nil {
		return nil, err
	}

	logging.Log("NOTIF-kSvPp").Debug("successfully initialized chat email and sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {
		contentText := message.GetContent()
		if config.Compact {
			contentText = html2text.HTML2Text(contentText)
		}
		for _, splittedMsg := range splitMessage(contentText, config.SplitCount) {
			if err := sendMessage(splittedMsg, url); err != nil {
				return err
			}
		}
		return nil
	}), nil
}

func sendMessage(message string, chatUrl *url.URL) error {
	req, err := json.Marshal(message)
	if err != nil {
		return caos_errs.ThrowInternal(err, "PROVI-s8uie", "Could not unmarshal content")
	}

	response, err := http.Post(chatUrl.String(), "application/json; charset=UTF-8", bytes.NewReader(req))
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
