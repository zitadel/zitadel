package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server/connect_middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func ExecutionHandler(alg crypto.EncryptionAlgorithm) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestTargets := execution.QueryExecutionTargetsForRequest(ctx, info.FullMethod)
		// call targets otherwise return req
		handledReq, err := executeTargetsForRequest(ctx, requestTargets, info.FullMethod, req, alg)
		if err != nil {
			return nil, err
		}

		response, err := handler(ctx, handledReq)
		if err != nil {
			return nil, err
		}

		responseTargets := execution.QueryExecutionTargetsForResponse(ctx, info.FullMethod)
		return executeTargetsForResponse(ctx, responseTargets, info.FullMethod, handledReq, response, alg)
	}
}

func executeTargetsForRequest(ctx context.Context, targets []target_domain.Target, fullMethod string, req interface{}, alg crypto.EncryptionAlgorithm) (_ interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// if no targets are found, return without any calls
	if len(targets) == 0 {
		return req, nil
	}

	md, _ := metadata.FromIncomingContext(ctx)
	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoRequest{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    Message{req.(proto.Message)},
		Headers:    connect_middleware.SetRequestHeaders(md),
	}

	return execution.CallTargets(ctx, targets, info, alg)
}

func executeTargetsForResponse(ctx context.Context, targets []target_domain.Target, fullMethod string, req, resp interface{}, alg crypto.EncryptionAlgorithm) (_ interface{}, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// if no targets are found, return without any calls
	if len(targets) == 0 {
		return resp, nil
	}

	md, _ := metadata.FromIncomingContext(ctx)
	ctxData := authz.GetCtxData(ctx)
	info := &ContextInfoResponse{
		FullMethod: fullMethod,
		InstanceID: authz.GetInstance(ctx).InstanceID(),
		ProjectID:  ctxData.ProjectID,
		OrgID:      ctxData.OrgID,
		UserID:     ctxData.UserID,
		Request:    Message{req.(proto.Message)},
		Response:   Message{resp.(proto.Message)},
		Headers:    connect_middleware.SetRequestHeaders(md),
	}

	return execution.CallTargets(ctx, targets, info, alg)
}

var _ execution.ContextInfo = &ContextInfoRequest{}

type ContextInfoRequest struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    Message     `json:"request,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}

type Message struct {
	proto.Message
}

func (r *Message) MarshalJSON() ([]byte, error) {
	data, err := protojson.Marshal(r.Message)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Message) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, r.Message)
}

func (c *ContextInfoRequest) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoRequest) SetHTTPResponseBody(resp []byte) error {
	return json.Unmarshal(resp, &c.Request)
}

func (c *ContextInfoRequest) GetContent() interface{} {
	return c.Request.Message
}

var _ execution.ContextInfo = &ContextInfoResponse{}

type ContextInfoResponse struct {
	FullMethod string      `json:"fullMethod,omitempty"`
	InstanceID string      `json:"instanceID,omitempty"`
	OrgID      string      `json:"orgID,omitempty"`
	ProjectID  string      `json:"projectID,omitempty"`
	UserID     string      `json:"userID,omitempty"`
	Request    Message     `json:"request,omitempty"`
	Response   Message     `json:"response,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}

func (c *ContextInfoResponse) GetHTTPRequestBody() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return data
}

func (c *ContextInfoResponse) SetHTTPResponseBody(resp []byte) error {
	return json.Unmarshal(resp, &c.Response)
}

func (c *ContextInfoResponse) GetContent() interface{} {
	return c.Response.Message
}
