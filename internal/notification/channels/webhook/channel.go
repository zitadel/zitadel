package webhook

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func InitChannel(ctx context.Context, cfg Config) (channels.NotificationChannel[*messages.JSON], error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logging.Debug("successfully initialized webhook json channel")
	return channels.HandleMessageFunc[*messages.JSON](func(message *messages.JSON) error {
		requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		payload, err := message.GetContent()
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(requestCtx, cfg.Method, cfg.CallURL, strings.NewReader(payload))
		if err != nil {
			return err
		}
		if cfg.Headers != nil {
			req.Header = cfg.Headers
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
			return errors.ThrowUnknown(fmt.Errorf("calling url %s returned %s", cfg.CallURL, resp.Status), "WEBH-LBxU0", "webhook didn't return a success status")
		}
		logging.WithFields("calling_url", cfg.CallURL, "method", cfg.Method).Debug("webhook called")
		return nil
	}), nil
}
