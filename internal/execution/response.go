package execution

import (
	"context"
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ContextInfoResponse struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

func CallTargetsResponse(ctx context.Context,
	targets []*query.Target,
	info *ContextInfoResponse,
) (interface{}, error) {
	ret := info.Response
	for _, target := range targets {
		if target.Async {
			go func(target query.Target) {
				if _, err := CallTargetResponse(ctx, &target, info); err != nil {
					logging.OnError(err).Error(err)
				}
			}(*target)
		} else {
			resp, err := CallTargetResponse(ctx, target, info)
			if err != nil && target.InterruptOnError {
				return ret, err
			}
			if resp != nil {
				ret = resp
				info.Response = resp
			}
		}
	}
	return ret, nil
}

func CallTargetResponse(ctx context.Context,
	target *query.Target,
	info *ContextInfoResponse,
) (res interface{}, err error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	r := info.Response
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
