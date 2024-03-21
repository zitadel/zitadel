package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextInfo[R any] struct {
	FullMethod string `json:"fullMethod"`
	InstanceID string `json:"instanceID"`
	OrgID      string `json:"orgID"`
	ProjectID  string `json:"projectID"`
	UserID     string `json:"userID"`
	Request    R      `json:"request"`
}

func CallTargets[R any](ctx context.Context,
	targets []*query.Target,
	info *ContextInfo[R],
) (r R, err error) {
	r = info.Request
	for _, target := range targets {
		if target.Async {
			go CallTarget(ctx, target, info)
		} else {
			r, err = CallTarget(ctx, target, info)
			if err != nil && target.InterruptOnError {
				return r, err
			}
		}
	}
	return r, err
}

func CallTarget[R any](ctx context.Context,
	target *query.Target,
	info *ContextInfo[R],
) (R, error) {
	var rnil R
	data, err := json.Marshal(info)
	if err != nil {
		return rnil, err
	}

	switch target.TargetType {
	case domain.TargetTypeWebhook:
		return info.Request, webhook(ctx, target.URL, target.Timeout, data)
	case domain.TargetTypeRequestResponse:
		response, err := call(ctx, target.URL, target.Timeout, data)
		if err != nil {
			return rnil, err
		}

		if err := json.Unmarshal(response, rnil); err != nil {
			return rnil, err
		}
		return rnil, nil
	default:
		return rnil, zerrors.ThrowInternal(nil, "EXEC-auqnansr2m", "Errors.Execution.Unknown")
	}
}

func webhook(ctx context.Context, url string, timeout time.Duration, body []byte) error {
	_, err := call(ctx, url, timeout, body)
	return err
}

func call(ctx context.Context, url string, timeout time.Duration, body []byte) ([]byte, error) {
	contextWithCancel, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(contextWithCancel, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, zerrors.ThrowUnknown(nil, "EXEC-dra6yamk9g", "Errors.Execution.Failed")
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
