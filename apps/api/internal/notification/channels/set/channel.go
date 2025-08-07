package set

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func InitChannel(ctx context.Context, cfg Config) (channels.NotificationChannel, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logging.Debug("successfully initialized security event token json channel")
	return channels.HandleMessageFunc(func(message channels.Message) error {
		requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		msg, ok := message.(*messages.Form)
		if !ok {
			return zerrors.ThrowInternal(nil, "SET-K686U", "message is not SET")
		}
		payload, err := msg.GetContent()
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, cfg.CallURL, strings.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "calling_url", cfg.CallURL).Debug("security event token called")
		if resp.StatusCode == http.StatusOK ||
			resp.StatusCode == http.StatusAccepted ||
			resp.StatusCode == http.StatusNoContent {
			return nil
		}
		body, err := mapResponse(resp)
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "callURL", cfg.CallURL).
			OnError(err).Debug("error mapping response")
		if resp.StatusCode == http.StatusBadRequest {
			logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "callURL", cfg.CallURL, "status", resp.Status, "body", body).
				Error("security event token didn't return a success status")
			return nil
		}
		return zerrors.ThrowInternalf(err, "SET-DF3dq", "security event token to %s didn't return a success status: %s (%v)", cfg.CallURL, resp.Status, body)
	}), nil
}

func mapResponse(resp *http.Response) (map[string]any, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	requestError := make(map[string]any)
	err = json.Unmarshal(body, &requestError)
	if err != nil {
		return nil, err
	}
	return requestError, nil
}
