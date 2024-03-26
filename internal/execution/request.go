package execution

import (
	"context"
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextInfoRequest struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
}

func CallTargetsRequest(ctx context.Context,
	targets []*query.Target,
	info *ContextInfoRequest,
) (interface{}, error) {
	ret := info.Request
	for _, target := range targets {
		if target.Async {
			go func(target query.Target) {
				if _, err := CallTargetRequest(ctx, &target, info); err != nil {
					logging.OnError(err).Error(err)
				}
			}(*target)
		} else {
			resp, err := CallTargetRequest(ctx, target, info)
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

func CallTargetRequest(ctx context.Context,
	target *query.Target,
	info *ContextInfoRequest,
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
