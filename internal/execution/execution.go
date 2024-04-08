package execution

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextInfo interface {
	GetHTTPRequestBody() []byte
	GetContent() interface{}
	SetHTTPResponseBody([]byte) error
}

type Target interface {
	GetTargetID() string
	IsAsync() bool
	IsInterruptOnError() bool
	GetURL() string
	GetTargetType() domain.TargetType
	GetTimeout() time.Duration
}

func CallTargets(ctx context.Context,
	targets []Target,
	info ContextInfo,
) (interface{}, error) {
	for _, target := range targets {
		// get function to call the type of target
		callF, err := CallTargetFunc(target)
		if err != nil {
			return nil, err
		}

		// call target async, ignore response and if error occur
		if target.IsAsync() {
			go func(target Target) {
				if _, err := callF(ctx, info); err != nil {
					logging.WithFields("target", target.GetTargetID()).OnError(err).Info(err)
				}
			}(target)
			// else call synchronous, and handle error if interrupt is set
		} else {
			resp, err := callF(ctx, info)
			if err != nil && target.IsInterruptOnError() {
				return nil, err
			}
			if resp != nil {
				// error in unmarshalling
				if err := info.SetHTTPResponseBody(resp); err != nil {
					return nil, err
				}
			}
		}
	}
	return info.GetContent(), nil
}

type ContextInfoRequest interface {
	GetHTTPRequestBody() []byte
}

// CallTargetFunc get function to call the desired type of target
func CallTargetFunc(
	target Target,
) (func(ctx context.Context, info ContextInfoRequest) (res []byte, err error), error) {
	switch target.GetTargetType() {
	// get request, ignore response and return request and error for handling in list of targets
	case domain.TargetTypeWebhook:
		return func(ctx context.Context, info ContextInfoRequest) (res []byte, err error) {
			return nil, webhook(ctx, target.GetURL(), target.GetTimeout(), info.GetHTTPRequestBody())
		}, nil
	// get request, return response and error
	case domain.TargetTypeCall:
		return func(ctx context.Context, info ContextInfoRequest) (res []byte, err error) {
			return call(ctx, target.GetURL(), target.GetTimeout(), info.GetHTTPRequestBody())
		}, nil
	default:
		return nil, zerrors.ThrowInternal(nil, "EXEC-auqnansr2m", "Errors.Execution.Unknown")
	}
}

// webhook call a webhook, ignore the response but return the errror
func webhook(ctx context.Context, url string, timeout time.Duration, body []byte) error {
	_, err := call(ctx, url, timeout, body)
	return err
}

// call function to do a post HTTP request to a desired url with timeout
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
	return nil, zerrors.ThrowUnknown(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed")
}
