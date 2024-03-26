package execution

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func Test_Call(t *testing.T) {
	type args struct {
		ctx        context.Context
		timeout    time.Duration
		sleep      time.Duration
		method     string
		body       []byte
		respBody   []byte
		statusCode int
	}
	type res struct {
		body    []byte
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"not ok status",
			args{
				ctx:        context.Background(),
				timeout:    time.Minute,
				sleep:      time.Second,
				method:     http.MethodPost,
				body:       []byte("{\"request\": \"values\"}"),
				respBody:   []byte("{\"response\": \"values\"}"),
				statusCode: http.StatusBadRequest,
			},
			res{
				wantErr: true,
			},
		},
		{
			"timeout",
			args{
				ctx:        context.Background(),
				timeout:    time.Second,
				sleep:      time.Second,
				method:     http.MethodPost,
				body:       []byte("{\"request\": \"values\"}"),
				respBody:   []byte("{\"response\": \"values\"}"),
				statusCode: http.StatusOK,
			},
			res{
				wantErr: true,
			},
		},
		{
			"ok",
			args{
				ctx:        context.Background(),
				timeout:    time.Minute,
				sleep:      time.Second,
				method:     http.MethodPost,
				body:       []byte("{\"request\": \"values\"}"),
				respBody:   []byte("{\"response\": \"values\"}"),
				statusCode: http.StatusOK,
			},
			res{
				body: []byte("{\"response\": \"values\"}"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServerCall(t,
				tt.args.method,
				tt.args.body,
				tt.args.sleep,
				tt.args.statusCode,
				tt.args.respBody,
				testCall(tt.args.ctx, tt.args.timeout, tt.args.body),
			)
			if tt.res.wantErr {
				assert.Error(t, err)
				assert.Nil(t, respBody)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.res.body, respBody)
			}
		})
	}
}

func testCall(ctx context.Context, timeout time.Duration, body []byte) func(string) (interface{}, error) {
	return func(url string) (interface{}, error) {
		return call(ctx, url, timeout, body)
	}
}

func testTargetCall(ctx context.Context,
	target *query.Target,
	info *ContextInfo,
) func(string) (interface{}, error) {
	return func(url string) (r interface{}, err error) {
		target.URL = url
		return CallTarget(ctx, target, info)
	}
}

func testServerCall(
	t *testing.T,
	method string,
	body []byte,
	timeout time.Duration,
	statusCode int,
	respBody []byte,
	call func(string) (interface{}, error),
) (interface{}, error) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		checkRequest(t, r, method, body)

		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
			return
		}

		time.Sleep(timeout)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(respBody))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	return call(server.URL)
}

func checkRequest(t *testing.T, sent *http.Request, method string, expectedBody []byte) {
	sentBody, err := io.ReadAll(sent.Body)
	require.NoError(t, err)
	require.Equal(t, expectedBody, sentBody)
	require.Equal(t, method, sent.Method)
}

func Test_CallTarget(t *testing.T) {
	type args struct {
		ctx    context.Context
		target *query.Target
		sleep  time.Duration

		info *ContextInfo

		method string
		body   []byte

		respBody   []byte
		statusCode int
	}
	type res struct {
		body    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"marshal error",
			args{
				ctx: context.Background(),
				info: &ContextInfo{
					FullMethod: "method",
					InstanceID: "instance",
					OrgID:      "org",
					ProjectID:  "project",
					UserID:     "user",
					Request:    make(chan int),
				},
				target: &query.Target{
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"request response, unmarshall error",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info: &ContextInfo{
					FullMethod: "method",
					InstanceID: "instance",
					OrgID:      "org",
					ProjectID:  "project",
					UserID:     "user",
					Request: &execution.SetExecutionRequest{
						Targets: []string{"target"},
					},
				},
				target: &query.Target{
					TargetType: domain.TargetTypeRequestResponse,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"fullMethod\":\"method\",\"instanceID\":\"instance\",\"orgID\":\"org\",\"projectID\":\"project\",\"userID\":\"user\",\"request\":{\"targets\":[\"target\"]}}"),
				respBody:   []byte("{\"unavailable\":[\"no\"]"),
				statusCode: http.StatusOK,
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, ok",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info: &ContextInfo{
					FullMethod: "method",
					InstanceID: "instance",
					OrgID:      "org",
					ProjectID:  "project",
					UserID:     "user",
					Request: &execution.SetExecutionRequest{
						Targets: []string{"target"},
					},
				},
				target: &query.Target{
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"fullMethod\":\"method\",\"instanceID\":\"instance\",\"orgID\":\"org\",\"projectID\":\"project\",\"userID\":\"user\",\"request\":{\"targets\":[\"target\"]}}"),
				respBody:   []byte("{\"targets\":[\"target\"]}"),
				statusCode: http.StatusOK,
			},
			res{
				body: &execution.SetExecutionRequest{
					Targets: []string{"target"},
				},
			},
		},
		{
			"request response, ok",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info: &ContextInfo{
					FullMethod: "method",
					InstanceID: "instance",
					OrgID:      "org",
					ProjectID:  "project",
					UserID:     "user",
					Request: &execution.SetExecutionRequest{
						Targets: []string{"target"},
					},
				},
				target: &query.Target{
					TargetType: domain.TargetTypeRequestResponse,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"fullMethod\":\"method\",\"instanceID\":\"instance\",\"orgID\":\"org\",\"projectID\":\"project\",\"userID\":\"user\",\"request\":{\"targets\":[\"target\"]}}"),
				respBody:   []byte("{\"targets\":[\"target\"]}"),
				statusCode: http.StatusOK,
			},
			res{
				body: &execution.SetExecutionRequest{
					Targets: []string{"target"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServerCall(t,
				tt.args.method,
				tt.args.body,
				tt.args.sleep,
				tt.args.statusCode,
				tt.args.respBody,
				testTargetCall(tt.args.ctx, tt.args.target, tt.args.info),
			)
			if tt.res.wantErr {
				assert.Error(t, err)
				assert.Nil(t, respBody)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.res.body, respBody)
			}
		})
	}
}
