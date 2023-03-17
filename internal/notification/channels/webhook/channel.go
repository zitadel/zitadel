package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func InitWebhookChannel(ctx context.Context, cfg Config) (channels.NotificationChannel, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logging.Debug("successfully initialized webhook json channel")
	return channels.HandleMessageFunc(func(message channels.Message) error {

		msg, ok := message.(*messages.JSON)
		if !ok {
			return caos_errs.ThrowInternal(nil, "WEBH-K686U", "message is not JSON")
		}
		payload, err := msg.GetContent()
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.CallURL, strings.NewReader(payload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if err = resp.Body.Close(); err != nil {
			return err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return caos_errs.ThrowUnknown(fmt.Errorf("calling url %s returned %s", cfg.CallURL, resp.Status), "WEBH-LBxU0", "webhook didn't return a success status")
		}

		logging.WithFields("calling_url", cfg.CallURL, "method", cfg.Method).Debug("webhook called")
		return nil
	}), nil
}
