package execution

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func webhook(ctx context.Context, url string, timeout time.Duration, body []byte) error {
	_, err := call(ctx, url, timeout, body)
	return err
}

func call(ctx context.Context, url string, timeout time.Duration, body []byte) ([]byte, error) {
	contextWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(contextWithCancel, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for success, if redirect has to be done, or return error with statusCode >= 400
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return io.ReadAll(resp.Body)
	} else if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		redirectUrl, err := resp.Location()
		// redirectURL is empty or not parsable
		if err != nil {
			return nil, err
		}
		req.URL = redirectUrl
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, zerrors.ThrowUnknown(nil, "EXEC-dra6yamk9g", "Errors.Execution.Failed")
		}
		return io.ReadAll(resp.Body)
	}

	return nil, zerrors.ThrowUnknown(nil, "EXEC-dra6yamk9g", "Errors.Execution.Failed")
}
