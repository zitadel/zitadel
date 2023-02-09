package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zitadel/zitadel/internal/repository/quota"
)

// ReportUsage calls notification hooks and emits the notified events
func (c *Commands) ReportUsage(ctx context.Context, dueNotifications []*quota.NotifiedEvent) error {
	for _, notification := range dueNotifications {

		if err := notify(ctx, notification); err != nil {
			if err != nil {
				return err
			}
		}

		if _, err := c.eventstore.Push(ctx, notification); err != nil {
			return err
		}
	}

	return nil
}

func notify(ctx context.Context, notification *quota.NotifiedEvent) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notification.CallURL, bytes.NewReader(payload))
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
		return fmt.Errorf("calling url %s returned %s", notification.CallURL, resp.Status)
	}

	return nil
}
