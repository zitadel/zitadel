package connect_middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
)

func newMockContentRequest(content string) *connect.Request[structpb.Struct] {
	return connect.NewRequest(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"content": {
				Kind: &structpb.Value_StringValue{StringValue: content},
			},
		},
	})
}

func newMockContentResponse(content string) *connect.Response[structpb.Struct] {
	return connect.NewResponse(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"content": {
				Kind: &structpb.Value_StringValue{StringValue: content},
			},
		},
	})
}

func newMockContextInfoRequest(fullMethod, request string) *ContextInfoRequest {
	return &ContextInfoRequest{
		FullMethod: fullMethod,
		Request:    Message{Message: newMockContentRequest(request).Msg},
	}
}

func newMockContextInfoResponse(fullMethod, request, response string) *ContextInfoResponse {
	return &ContextInfoResponse{
		FullMethod: fullMethod,
		Request:    Message{Message: newMockContentRequest(request).Msg},
		Response:   Message{Message: newMockContentResponse(response).Msg},
	}
}

func Test_executeTargetsForGRPCFullMethod_request(t *testing.T) {
	type target struct {
		reqBody    execution.ContextInfo
		sleep      time.Duration
		statusCode int
		respBody   connect.AnyResponse
	}
	type args struct {
		ctx context.Context

		executionTargets []target_domain.Target
		targets          []target
		fullMethod       string
		req              connect.AnyRequest
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
				executionTargets: []target_domain.Target{},
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  target_domain.TargetTypeCall,
						Timeout:     time.Minute,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},

				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  target_domain.TargetTypeAsync,
						Timeout:     time.Second,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID: "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:    "target",
						TargetType:  target_domain.TargetTypeAsync,
						Timeout:     time.Minute,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeWebhook,
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeWebhook,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeWebhook,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target1",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target2",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target3",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},

				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content1"),
						respBody:   newMockContentResponse("content2"),
						sleep:      0,
						statusCode: http.StatusBadRequest,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content2"),
						respBody:   newMockContentResponse("content3"),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target1",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target2",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target3",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Second,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse("content1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content1"),
						respBody:   newMockContentResponse("content2"),
						sleep:      5 * time.Second,
						statusCode: http.StatusBadRequest,
					},
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content2"),
						respBody:   newMockContentResponse("content3"),
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

				tt.args.executionTargets[i].Endpoint = url
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForRequest(
				tt.args.ctx,
				tt.args.executionTargets,
				tt.args.fullMethod,
				tt.args.req,
				nil,
			)

			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualExportedValues(t, tt.res.want, resp)

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
	respBody connect.AnyResponse,
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
		resp, err := protojson.Marshal(respBody.Any().(proto.Message))
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(resp); err != nil {
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
		respBody   connect.AnyResponse
	}
	type args struct {
		ctx context.Context

		executionTargets []target_domain.Target
		targets          []target
		fullMethod       string
		req              connect.AnyRequest
		resp             connect.AnyResponse
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
				resp:             newMockContentResponse("response"),
			},
			res{
				want: newMockContentResponse("response"),
			},
		},
		{
			"target, executionTargets empty",
			args{
				ctx:              context.Background(),
				fullMethod:       "/service/method",
				executionTargets: []target_domain.Target{},
				req:              newMockContentRequest("request"),
				resp:             newMockContentResponse("response"),
			},
			res{
				want: newMockContentResponse("response"),
			},
		},
		{
			"target, empty response",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "request./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoRequest("/service/method", "content"),
						respBody:   newMockContentResponse(""),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req:  newMockContentRequest(""),
				resp: newMockContentResponse(""),
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
				executionTargets: []target_domain.Target{
					{
						ExecutionID:      "response./zitadel.session.v2.SessionService/SetSession",
						TargetID:         "target",
						TargetType:       target_domain.TargetTypeCall,
						Timeout:          time.Minute,
						InterruptOnError: true,
					},
				},
				targets: []target{
					{
						reqBody:    newMockContextInfoResponse("/service/method", "request", "response"),
						respBody:   newMockContentResponse("response1"),
						sleep:      0,
						statusCode: http.StatusOK,
					},
				},
				req:  newMockContentRequest("request"),
				resp: newMockContentResponse("response"),
			},
			res{
				want: newMockContentResponse("response1"),
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

				tt.args.executionTargets[i].Endpoint = url
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForResponse(
				tt.args.ctx,
				tt.args.executionTargets,
				tt.args.fullMethod,
				tt.args.req,
				tt.args.resp,
				nil,
			)

			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualExportedValues(t, tt.res.want, resp)

			for _, closeF := range closeFuncs {
				closeF()
			}
		})
	}
}

func Test_setRequestHeaders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		reqHeaders map[string][]string
		want       map[string][]string
	}{
		{
			name:       "no headers",
			reqHeaders: nil,
			want:       nil,
		},
		{
			name:       "with headers",
			reqHeaders: map[string][]string{"Authorization": {"Bearer XXX"}, "X-Random-Header": {"Random-Value"}, "X-Forwarded-For": {"1.2.3.4"}, "Host": {"localhost:8080"}},
			want:       map[string][]string{"X-Forwarded-For": {"1.2.3.4"}, "Host": {"localhost:8080"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SetRequestHeaders(tt.reqHeaders)
			assert.Equal(t, tt.want, got)
		})
	}
}
