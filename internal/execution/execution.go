package execution

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/logging"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextInfo interface {
	GetHTTPRequestBody() []byte
	GetContent() interface{}
	SetHTTPResponseBody([]byte) error
}

type Target interface {
	GetTargetID() string
	IsInterruptOnError() bool
	GetEndpoint() string
	GetTargetType() domain.TargetType
	GetTimeout() time.Duration
}

// CallTargets call a list of targets in order with handling of error and responses
func CallTargets(
	ctx context.Context,
	targets []Target,
	info ContextInfo,
) (_ interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	for _, target := range targets {
		// call the type of target
		resp, err := CallTarget(ctx, target, info)
		// handle error if interrupt is set
		if err != nil && target.IsInterruptOnError() {
			return nil, err
		}
		if len(resp) > 0 {
			// error in unmarshalling
			if err := info.SetHTTPResponseBody(resp); err != nil {
				return nil, err
			}
		}
	}
	return info.GetContent(), nil
}

type ContextInfoRequest interface {
	GetHTTPRequestBody() []byte
}

// CallTarget call the desired type of target with handling of responses
func CallTarget(
	ctx context.Context,
	target Target,
	info ContextInfoRequest,
) (res []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	switch target.GetTargetType() {
	// get request, ignore response and return request and error for handling in list of targets
	case domain.TargetTypeWebhook:
		return nil, webhook(ctx, target.GetEndpoint(), target.GetTimeout(), info.GetHTTPRequestBody())
	// get request, return response and error
	case domain.TargetTypeCall:
		return call(ctx, target.GetEndpoint(), target.GetTimeout(), info.GetHTTPRequestBody())
	case domain.TargetTypeAsync:
		go func(target Target, info ContextInfoRequest) {
			if _, err := call(ctx, target.GetEndpoint(), target.GetTimeout(), info.GetHTTPRequestBody()); err != nil {
				logging.WithFields("target", target.GetTargetID()).OnError(err).Info(err)
			}
		}(target, info)
		return nil, nil
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
func call(ctx context.Context, url string, timeout time.Duration, body []byte) (_ []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		cancel()
		span.EndWithError(err)
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return handleResponse(resp)
}

func handleResponse(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Check for success between 200 and 299, redirect 300 to 399 is handled by the client, return error with statusCode >= 400
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return data, nil
	}
	return nil, zhttp.HTTPStatusCodeToZitadelError(resp.StatusCode, "EXEC-dra6yamk98", string(data))
}
