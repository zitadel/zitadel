package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
)

var _ execution.Target = &mockExecutionTarget{}

type mockExecutionTarget struct {
	InstanceID       string
	ExecutionID      string
	TargetID         string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

func (e *mockExecutionTarget) SetEndpoint(endpoint string) {
	e.Endpoint = endpoint
}
func (e *mockExecutionTarget) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *mockExecutionTarget) GetEndpoint() string {
	return e.Endpoint
}
func (e *mockExecutionTarget) GetTargetType() domain.TargetType {
	return e.TargetType
}
func (e *mockExecutionTarget) GetTimeout() time.Duration {
	return e.Timeout
}
func (e *mockExecutionTarget) GetTargetID() string {
	return e.TargetID
}
func (e *mockExecutionTarget) GetExecutionID() string {
	return e.ExecutionID
}

type mockContentRequest struct {
	Content string
}

func newMockContentRequest(content string) *mockContentRequest {
	return &mockContentRequest{
		Content: content,
	}
}

func newMockContextInfoRequest(fullMethod, request string) *ContextInfoRequest {
	return &ContextInfoRequest{
		FullMethod: fullMethod,
		Request:    newMockContentRequest(request),
	}
}

func newMockContextInfoResponse(fullMethod, request, response string) *ContextInfoResponse {
	return &ContextInfoResponse{
		FullMethod: fullMethod,
		Request:    newMockContentRequest(request),
		Response:   newMockContentRequest(response),
	}
}

func Test_executeTargetsForGRPCFullMethod_request(t *testing.T) {
	type target struct {
		reqBody    execution.ContextInfo
		sleep      time.Duration
		statusCode int
		respBody   interface{}
	}
	type args struct {
		ctx context.Context

		executionTargets []execution.Target
		targets          []target
		fullMethod       string
		req              interface{}
	}
	type res struct {
		want    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"target, executionTargets nil",
			args{
				ctx:              context.Background(),
				fullMethod:       "/service/method",
				executionTargets: nil,
				req:              newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, executionTargets empty",
			args{
				ctx:              context.Background(),
				fullMethod:       "/service/method",
				executionTargets: []execution.Target{},
				req:              newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, not reachable",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{},
				req:     newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"target, error without interrupt",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:  "instance",
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  domain.TargetTypeCall,
						Timeout:     time.Minute,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusBadRequest,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want: newMockContentRequest("content"),
			},
		},
		{
			"target, interruptOnError",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},

				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusBadRequest,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"target, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      5 * time.Second,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"target, wrong request",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{reqBody: newMockContextInfoRequest("/service/method", "wrong")},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"target, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want: newMockContentRequest("content1"),
			},
		},
		{
			"target async, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:  "instance",
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  domain.TargetTypeAsync,
						Timeout:     time.Second,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      5 * time.Second,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want: newMockContentRequest("content"),
			},
		},
		{
			"target async, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:  "instance",
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  domain.TargetTypeAsync,
						Timeout:     time.Minute,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want: newMockContentRequest("content"),
			},
		},
		{
			"webhook, error",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						sleep:      0,
						statusCode: http.StatusInternalServerError,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      5 * time.Second,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want: newMockContentRequest("content"),
			},
		},
		{
			"with includes, interruptOnError",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target1",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target2",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target3",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},

				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content1"),
						respBody:   newMockContentRequest("content2"),
						sleep:      0,
						statusCode: http.StatusBadRequest,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content2"),
						respBody:   newMockContentRequest("content3"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
		{
			"with includes, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target1",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target2",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target3",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content1"),
						respBody:   newMockContentRequest("content2"),
						sleep:      5 * time.Second,
						statusCode: http.StatusBadRequest,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content2"),
						respBody:   newMockContentRequest("content3"),
						sleep:      5 * time.Second,
						statusCode: http.StatusOK,
					},
				},
				req: newMockContentRequest("content"),
			},
			res{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closeFuncs := make([]func(), len(tt.args.targets))
			for i, target := range tt.args.targets {
				url, closeF := testServerCall(
					target.reqBody,
					target.sleep,
					target.statusCode,
					target.respBody,
				)

				et := tt.args.executionTargets[i].(*mockExecutionTarget)
				et.SetEndpoint(url)
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForRequest(
				tt.args.ctx,
				tt.args.executionTargets,
				tt.args.fullMethod,
				tt.args.req,
			)

			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.want, resp)

			for _, closeF := range closeFuncs {
				closeF()
			}
		})
	}
}

func testServerCall(
	reqBody interface{},
	sleep time.Duration,
	statusCode int,
	respBody interface{},
) (string, func()) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		sentBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		if !reflect.DeepEqual(data, sentBody) {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
			return
		}

		time.Sleep(sleep)

		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		if _, err := io.WriteString(w, string(resp)); err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	return server.URL, server.Close
}

func Test_executeTargetsForGRPCFullMethod_response(t *testing.T) {
	type target struct {
		reqBody    execution.ContextInfo
		sleep      time.Duration
		statusCode int
		respBody   interface{}
	}
	type args struct {
		ctx context.Context

		executionTargets []execution.Target
		targets          []target
		fullMethod       string
		req              interface{}
		resp             interface{}
	}
	type res struct {
		want    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"target, executionTargets nil",
			args{
				ctx:              context.Background(),
				fullMethod:       "/service/method",
				executionTargets: nil,
				req:              newMockContentRequest("request"),
				resp:             newMockContentRequest("response"),
			},
			res{
				want: newMockContentRequest("response"),
			},
		},
		{
			"target, executionTargets empty",
			args{
				ctx:              context.Background(),
				fullMethod:       "/service/method",
				executionTargets: []execution.Target{},
				req:              newMockContentRequest("request"),
				resp:             newMockContentRequest("response"),
			},
			res{
				want: newMockContentRequest("response"),
			},
		},
		{
			"target, empty response",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentRequest(""),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req: []byte{},
			},
			res{
				wantErr: true,
			},
		},
		{
			"target, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []execution.Target{
					&mockExecutionTarget{
						InstanceID:       "instance",
						ExecutionID:      "response./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoResponse("/service/method", "request", "response"),
						respBody:   newMockContentRequest("response1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req:  newMockContentRequest("request"),
				resp: newMockContentRequest("response"),
			},
			res{
				want: newMockContentRequest("response1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closeFuncs := make([]func(), len(tt.args.targets))
			for i, target := range tt.args.targets {
				url, closeF := testServerCall(
					target.reqBody,
					target.sleep,
					target.statusCode,
					target.respBody,
				)

				et := tt.args.executionTargets[i].(*mockExecutionTarget)
				et.SetEndpoint(url)
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForResponse(
				tt.args.ctx,
				tt.args.executionTargets,
				tt.args.fullMethod,
				tt.args.req,
				tt.args.resp,
			)

			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.want, resp)

			for _, closeF := range closeFuncs {
				closeF()
			}
		})
	}
}
