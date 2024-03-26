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

type ContextInfo struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

func CallTargets(ctx context.Context,
	targets []*query.Target,
	info *ContextInfo,
) (interface{}, error) {
	ret := info.Request
	for _, target := range targets {
		if target.Async {
			go func() {
				if _, err := CallTarget(ctx, target, info); err != nil {
					// TODO: log if error on async action?
				}
			}()
		} else {
			resp, err := CallTarget(ctx, target, info)
			if err != nil && target.InterruptOnError {
				return ret, err
			}
			if resp != nil {
				ret = resp
				info.Request = resp
			}
		}
	}
	return ret, nil
}

func CallTarget(ctx context.Context,
	target *query.Target,
	info *ContextInfo,
) (res interface{}, err error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	r := info.Request
	switch target.TargetType {
	case domain.TargetTypeWebhook:
		return r, webhook(ctx, target.URL, target.Timeout, data)
	case domain.TargetTypeRequestResponse:
		response, err := call(ctx, target.URL, target.Timeout, data)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(response, r); err != nil {
			return nil, err
		}
		return r, nil
	default:
		return nil, zerrors.ThrowInternal(nil, "EXEC-auqnansr2m", "Errors.Execution.Unknown")
	}
}

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
	if resp.StatusCode != http.StatusOK {
		return nil, zerrors.ThrowUnknown(nil, "EXEC-dra6yamk9g", "Errors.Execution.Failed")
	}

	return io.ReadAll(resp.Body)
}
