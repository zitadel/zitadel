package middleware

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
)

func ExecutionHandler(queries *query.Queries) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestTargets, responseTargets := queryTargets(ctx, queries, info.FullMethod)

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

func executeTargetsForRequest(ctx context.Context, targets []execution.Target, fullMethod string, req interface{}) (interface{}, error) {
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

	request, err := execution.CallTargets(ctx, targets, info)
	// error gets only returned if InterruptOnError is set, or internal errors occur
	if err != nil {
		// if an error is returned still return also the original request
		return req, err
	}
	return request, err
}

func executeTargetsForResponse(ctx context.Context, targets []execution.Target, fullMethod string, req, resp interface{}) (interface{}, error) {
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

	response, err := execution.CallTargets(ctx, targets, info)
	if err != nil {
		// if an error is returned still return also the original response
		return resp, err
	}
	return response, err
}

type ExecutionQueries interface {
	ExecutionTargetsCombined(ctx context.Context, ids1, ids2 []string) (execution []*query.ExecutionTarget, err error)
}

func queryTargets(
	ctx context.Context,
	queries ExecutionQueries,
	fullMethod string,
) ([]execution.Target, []execution.Target) {
	requestTargets := make([]execution.Target, 0)
	responseTargets := make([]execution.Target, 0)
	targets, err := queries.ExecutionTargetsCombined(ctx,
		idsForFullMethod(fullMethod, domain.ExecutionTypeRequest),
		idsForFullMethod(fullMethod, domain.ExecutionTypeResponse),
	)
	if err != nil {
		logging.WithFields("fullMethod", fullMethod).WithError(err).Info("unable to query targets")
		return requestTargets, responseTargets
	}

	for _, target := range targets {
		if strings.HasPrefix(target.GetExecutionID(), exec_repo.IDAll(domain.ExecutionTypeRequest)) {
			requestTargets = append(requestTargets, target)
		} else if strings.HasPrefix(target.GetExecutionID(), exec_repo.IDAll(domain.ExecutionTypeResponse)) {
			responseTargets = append(responseTargets, target)
		}
	}

	return requestTargets, responseTargets
}

func idsForFullMethod(fullMethod string, executionType domain.ExecutionType) []string {
	return []string{exec_repo.ID(executionType, fullMethod), exec_repo.ID(executionType, serviceFromFullMethod(fullMethod)), exec_repo.IDAll(executionType)}
}

func serviceFromFullMethod(s string) string {
	parts := strings.Split(s, "/")
	return parts[1]
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
