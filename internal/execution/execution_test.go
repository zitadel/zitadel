package execution_test

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
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/execution"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/actions"
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
		signingKey string
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
				sleep:      2 * time.Second,
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
		{
			"ok, signed",
			args{
				ctx:        context.Background(),
				timeout:    time.Minute,
				sleep:      time.Second,
				method:     http.MethodPost,
				body:       []byte("{\"request\": \"values\"}"),
				respBody:   []byte("{\"response\": \"values\"}"),
				statusCode: http.StatusOK,
				signingKey: "signingkey",
			},
			res{
				body: []byte("{\"response\": \"values\"}"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServer(t,
				&callTestServer{
					method:      tt.args.method,
					expectBody:  tt.args.body,
					timeout:     tt.args.sleep,
					statusCode:  tt.args.statusCode,
					respondBody: tt.args.respBody,
				},
				testCall(tt.args.ctx, tt.args.timeout, tt.args.body, tt.args.signingKey),
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

func Test_CallTarget(t *testing.T) {
	type args struct {
		ctx    context.Context
		info   *middleware.ContextInfoRequest
		server *callTestServer
		target target_domain.Target
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
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					timeout:     time.Second,
					statusCode:  http.StatusInternalServerError,
				},
				target: target_domain.Target{
					TargetType: 4,
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, error",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusInternalServerError,
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"webhook, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusOK,
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeWebhook,
					Timeout:    time.Minute,
				},
			},
			res{
				body: nil,
			},
		},
		{
			"webhook, signed, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusOK,
					signingKey:  "signingkey",
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeWebhook,
					Timeout:    time.Minute,
					SigningKey: &crypto.CryptoValue{
						Algorithm: "enc",
						KeyID:     "id",
						Crypted:   []byte("signingkey"),
					},
				},
			},
			res{
				body: nil,
			},
		},
		{
			"request response, error",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusInternalServerError,
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeCall,
					Timeout:    time.Minute,
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"request response, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusOK,
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeCall,
					Timeout:    time.Minute,
				},
			},
			res{
				body: []byte("{\"content\":\"request2\"}"),
			},
		},
		{
			"request response, signed, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  []byte("{\"request\":{\"content\":\"request1\"}}"),
					respondBody: []byte("{\"content\":\"request2\"}"),
					statusCode:  http.StatusOK,
					signingKey:  "signingkey",
				},
				target: target_domain.Target{
					TargetType: target_domain.TargetTypeCall,
					Timeout:    time.Minute,
					SigningKey: &crypto.CryptoValue{
						Algorithm: "enc",
						KeyID:     "id",
						Crypted:   []byte("signingkey"),
					},
				},
			},
			res{
				body: []byte("{\"content\":\"request2\"}"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServer(t, tt.args.server, testCallTarget(tt.args.ctx, tt.args.info, tt.args.target, crypto.CreateMockEncryptionAlg(gomock.NewController(t))))
			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.body, respBody)
		})
	}
}

func Test_CallTargets(t *testing.T) {
	type args struct {
		ctx     context.Context
		info    *middleware.ContextInfoRequest
		servers []*callTestServer
		targets []target_domain.Target
	}
	type res struct {
		ret     interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"interrupt on status",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: true},
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"continue on status",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: false},
				},
			},
			res{
				ret: requestContextInfo1.GetContent(),
			},
		},
		{
			"interrupt on json error",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: []byte("just a string, not json"),
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: true},
				},
			},
			res{
				wantErr: true,
			},
		},
		{
			"continue on json error",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  requestContextInfoBody1,
					respondBody: []byte("just a string, not json"),
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: false},
				}},
			res{
				ret: requestContextInfo1.GetContent(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServers(t,
				tt.args.servers,
				testCallTargets(tt.args.ctx, tt.args.info, tt.args.targets, crypto.CreateMockEncryptionAlg(gomock.NewController(t))),
			)
			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.ret, respBody)
		})
	}
}

type callTestServer struct {
	method      string
	expectBody  []byte
	timeout     time.Duration
	statusCode  int
	respondBody []byte
	signingKey  string
}

func testServers(
	t *testing.T,
	c []*callTestServer,
	call func([]string) (interface{}, error),
) (interface{}, error) {
	urls := make([]string, len(c))
	for i := range c {
		url, close := listen(t, c[i])
		defer close()
		urls[i] = url
	}
	return call(urls)
}

func testServer(
	t *testing.T,
	c *callTestServer,
	call func(string) ([]byte, error),
) ([]byte, error) {
	url, close := listen(t, c)
	defer close()
	return call(url)
}

func listen(
	t *testing.T,
	c *callTestServer,
) (url string, close func()) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		checkRequest(t, r, c.method, c.expectBody, c.signingKey)

		if c.statusCode != http.StatusOK {
			http.Error(w, "error", c.statusCode)
			return
		}

		time.Sleep(c.timeout)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(c.respondBody); err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	return server.URL, server.Close
}

func checkRequest(t *testing.T, sent *http.Request, method string, expectedBody []byte, signingKey string) {
	sentBody, err := io.ReadAll(sent.Body)
	require.NoError(t, err)
	require.Equal(t, expectedBody, sentBody)
	require.Equal(t, method, sent.Method)
	if signingKey != "" {
		require.NoError(t, actions.ValidatePayload(sentBody, sent.Header.Get(actions.SigningHeader), signingKey))
	}
}

func testCall(ctx context.Context, timeout time.Duration, body []byte, signingKey string) func(string) ([]byte, error) {
	return func(url string) ([]byte, error) {
		return execution.Call(ctx, url, timeout, body, signingKey)
	}
}

func testCallTarget(ctx context.Context,
	info *middleware.ContextInfoRequest,
	target target_domain.Target,
	alg crypto.EncryptionAlgorithm,
) func(string) ([]byte, error) {
	return func(url string) (r []byte, err error) {
		target.Endpoint = url
		return execution.CallTarget(ctx, target, info, alg)
	}
}

func testCallTargets(ctx context.Context,
	info *middleware.ContextInfoRequest,
	target []target_domain.Target,
	alg crypto.EncryptionAlgorithm,
) func([]string) (interface{}, error) {
	return func(urls []string) (interface{}, error) {
		targets := make([]target_domain.Target, len(target))
		for i, t := range target {
			t.Endpoint = urls[i]
			targets[i] = t
		}
		return execution.CallTargets(ctx, targets, info, alg)
	}
}

var requestContextInfo1 = &middleware.ContextInfoRequest{
	Request: middleware.Message{Message: &structpb.Struct{
		Fields: map[string]*structpb.Value{"content": structpb.NewStringValue("request1")},
	}},
}

var requestContextInfoBody1 = []byte("{\"request\":{\"content\":\"request1\"}}")
var requestContextInfoBody2 = []byte("{\"request\":{\"content\":\"request2\"}}")

type request struct {
	Request string `json:"request"`
}

func testErrorBody(code int, message string) []byte {
	body := &execution.ErrorBody{ForwardedStatusCode: code, ForwardedErrorMessage: message}
	data, _ := json.Marshal(body)
	return data
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
			"response, statuscode unknown and body",
			args{
				resp: &http.Response{
					StatusCode: 1000,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed"))
				},
			},
		},
		{
			"response, statuscode >= 400 and no body",
			args{
				resp: &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed"))
				},
			},
		},
		{
			"response, statuscode >= 400 and body",
			args{
				resp: &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed"))
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
		{
			"response, statuscode = 200 no body",
			args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				},
			},
			res{
				data:    []byte(""),
				wantErr: nil,
			},
		},
		{
			"response, statuscode = 200, error body >= 400 < 500",
			args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(testErrorBody(http.StatusForbidden, "forbidden"))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "EXEC-reUaUZCzCp", "forbidden"))
				},
			},
		},
		{
			"response, statuscode = 200, error body >= 500",
			args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(testErrorBody(http.StatusInternalServerError, "internal"))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-bmhNhpcqpF", "internal"))
				},
			},
		},
		{
			"response, statuscode = 308, no body, should not happen",
			args{
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(testErrorBody(http.StatusPermanentRedirect, "redirect"))),
				},
			},
			res{
				wantErr: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "EXEC-bmhNhpcqpF", "redirect"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := execution.HandleResponse(
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
