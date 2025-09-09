package connect_middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var headersToForward = map[string]bool{
	strings.ToLower(http_utils.ContentType):   true,
	strings.ToLower(http_utils.ForwardedFor):  true,
	strings.ToLower(http_utils.ForwardedHost): true,
	strings.ToLower(http_utils.Host):          true,
	strings.ToLower(http_utils.Origin):        true,
}

func ExecutionHandler(alg crypto.EncryptionAlgorithm) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {

			requestTargets := execution.QueryExecutionTargetsForRequest(ctx, req.Spec().Procedure)
			handledReq, err := executeTargetsForRequest(ctx, requestTargets, req.Spec().Procedure, req, alg)
			if err != nil {
				return nil, err
			}

			response, err := handler(ctx, handledReq)
			if err != nil {
				return nil, err
			}

			responseTargets := execution.QueryExecutionTargetsForResponse(ctx, req.Spec().Procedure)
			return executeTargetsForResponse(ctx, responseTargets, req.Spec().Procedure, handledReq, response, alg)
		}
	}
}

func executeTargetsForRequest(ctx context.Context, targets []target_domain.Target, fullMethod string, req connect.AnyRequest, alg crypto.EncryptionAlgorithm) (_ connect.AnyRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
		Request:    Message{req.Any().(proto.Message)},
		Headers:    SetRequestHeaders(req.Header()),
	}

	_, err = execution.CallTargets(ctx, targets, info, alg)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func executeTargetsForResponse(ctx context.Context, targets []target_domain.Target, fullMethod string, req connect.AnyRequest, resp connect.AnyResponse, alg crypto.EncryptionAlgorithm) (_ connect.AnyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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
		Request:    Message{req.Any().(proto.Message)},
		Response:   Message{resp.Any().(proto.Message)},
		Headers:    SetRequestHeaders(req.Header()),
	}

	_, err = execution.CallTargets(ctx, targets, info, alg)
	if err != nil {
		return nil, err
	}
	return resp, nil
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

func SetRequestHeaders(reqHeaders map[string][]string) map[string][]string {
	if len(reqHeaders) == 0 {
		return nil
	}
	headers := make(map[string][]string)
	for k, v := range reqHeaders {
		if headersToForward[strings.ToLower(k)] {
			headers[k] = v
		}
	}
	return headers
}
