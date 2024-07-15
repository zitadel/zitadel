package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ Target = &mockTarget{}

type mockTarget struct {
	InstanceID       string
	ExecutionID      string
	TargetID         string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

func (e *mockTarget) GetTargetID() string {
	return e.TargetID
}
func (e *mockTarget) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *mockTarget) GetEndpoint() string {
	return e.Endpoint
}
func (e *mockTarget) GetTargetType() domain.TargetType {
	return e.TargetType
}
func (e *mockTarget) GetTimeout() time.Duration {
	return e.Timeout
}

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

func testCall(ctx context.Context, timeout time.Duration, body []byte) func(string) ([]byte, error) {
	return func(url string) ([]byte, error) {
		return call(ctx, url, timeout, body)
	}
}

func testCallTarget(ctx context.Context,
	target *mockTarget,
	info ContextInfoRequest,
) func(string) ([]byte, error) {
	return func(url string) (r []byte, err error) {
		target.Endpoint = url
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
	call func(string) ([]byte, error),
) ([]byte, error) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		checkRequest(t, r, method, body)

		if statusCode != http.StatusOK {
			http.Error(w, "error", statusCode)
			return
		}

		time.Sleep(timeout)

		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, string(respBody)); err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
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

var _ ContextInfoRequest = &mockContextInfoRequest{}

type request struct {
	Request string `json:"request"`
}

type mockContextInfoRequest struct {
	Request *request `json:"request"`
}

func newMockContextInfoRequest(s string) *mockContextInfoRequest {
	return &mockContextInfoRequest{&request{s}}
}

func (c *mockContextInfoRequest) GetHTTPRequestBody() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c *mockContextInfoRequest) GetContent() []byte {
	data, _ := json.Marshal(c.Request)
	return data
}

func Test_CallTarget(t *testing.T) {
	type args struct {
		ctx    context.Context
		target *mockTarget
		sleep  time.Duration

		info ContextInfoRequest

		method string
		body   []byte

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
			"unknown targettype, error",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info:   newMockContextInfoRequest("content1"),
				target: &mockTarget{
					TargetType: 4,
				},
				body:       []byte("{\"request\":{\"request\":\"content1\"}}"),
				respBody:   []byte("{\"request\":\"content2\"}"),
				statusCode: http.StatusInternalServerError,
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, error",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info:   newMockContextInfoRequest("content1"),
				target: &mockTarget{
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"request\":{\"request\":\"content1\"}}"),
				respBody:   []byte("{\"request\":\"content2\"}"),
				statusCode: http.StatusInternalServerError,
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
				info:   newMockContextInfoRequest("content1"),
				target: &mockTarget{
					TargetType: domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"request\":{\"request\":\"content1\"}}"),
				respBody:   []byte("{\"request\":\"content2\"}"),
				statusCode: http.StatusOK,
			},
			res{
				body: nil,
			},
		},
		{
			"request response, error",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info:   newMockContextInfoRequest("content1"),
				target: &mockTarget{
					TargetType: domain.TargetTypeCall,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"request\":{\"request\":\"content1\"}}"),
				respBody:   []byte("{\"request\":\"content2\"}"),
				statusCode: http.StatusInternalServerError,
			},
			res{
				wantErr: true,
			},
		},
		{
			"request response, ok",
			args{
				ctx:    context.Background(),
				sleep:  time.Second,
				method: http.MethodPost,
				info:   newMockContextInfoRequest("content1"),
				target: &mockTarget{
					TargetType: domain.TargetTypeCall,
					Timeout:    time.Minute,
				},
				body:       []byte("{\"request\":{\"request\":\"content1\"}}"),
				respBody:   []byte("{\"request\":\"content2\"}"),
				statusCode: http.StatusOK,
			},
			res{
				body: []byte("{\"request\":\"content2\"}"),
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
				testCallTarget(tt.args.ctx, tt.args.target, tt.args.info),
			)
			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.body, respBody)
		})
	}
}

func Test_handleResponse(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	type res struct {
		data    []byte
		wantErr func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"response, statuscode > 400",
			args{
				resp: &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "EXEC-dra6yamk98", ""))
				},
			},
		},
		{
			"response, statuscode > 400 and body",
			args{
				resp: &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "EXEC-dra6yamk98", "body"))
				}},
		},
		{
			"response, statuscode = 200 and body",
			args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
				},
			},
			res{
				data:    []byte("body"),
				wantErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := handleResponse(
				tt.args.resp,
			)

			if tt.res.wantErr == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.wantErr(err) {
				t.Errorf("got wrong err: %v", err)
				return
			}
			assert.Equal(t, tt.res.data, respBody)
		})
	}

}
