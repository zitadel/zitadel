package middleware

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func ExecutionHandler(queries *query.Queries) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestTargets, responseTargets := execution.QueryExecutionTargetsForRequestAndResponse(ctx, queries, info.FullMethod)

		// call targets otherwise return req
		handledReq, err := executeTargetsForRequest(ctx, requestTargets, info.FullMethod, req)
		if err != nil {
			return nil, err
		}

		response, err := handler(ctx, handledReq)
		if err != nil {
			return nil, err
		}

		return executeTargetsForResponse(ctx, responseTargets, info.FullMethod, handledReq, response)
	}
}

func executeTargetsForRequest(ctx context.Context, targets []execution.Target, fullMethod string, req interface{}) (_ interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	// if no targets are found, return without any calls
	if len(targets) == 0 {
		return req, nil
	}

	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoRequest{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    req,
	}

	return execution.CallTargets(ctx, targets, info)
}

func executeTargetsForResponse(ctx context.Context, targets []execution.Target, fullMethod string, req, resp interface{}) (_ interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer span.EndWithError(err)

	// if no targets are found, return without any calls
	if len(targets) == 0 {
		return resp, nil
	}

	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoResponse{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    req,
		Response:   resp,
	}

	return execution.CallTargets(ctx, targets, info)
}

var _ execution.ContextInfo = &ContextInfoRequest{}

type ContextInfoRequest struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
}

func (c *ContextInfoRequest) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoRequest) SetHTTPResponseBody(resp []byte) error {
	if !json.Valid(resp) {
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-4m9s2", "Errors.Execution.ResponseIsNotValidJSON")
	}
	return json.Unmarshal(resp, c.Request)
}

func (c *ContextInfoRequest) GetContent() interface{} {
	return c.Request
}

var _ execution.ContextInfo = &ContextInfoResponse{}

type ContextInfoResponse struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    interface{} `json:"request,omitempty"`
	Response   interface{} `json:"response,omitempty"`
}

func (c *ContextInfoResponse) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoResponse) SetHTTPResponseBody(resp []byte) error {
	return json.Unmarshal(resp, c.Response)
}

func (c *ContextInfoResponse) GetContent() interface{} {
	return c.Response
}
