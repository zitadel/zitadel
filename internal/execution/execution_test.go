package execution_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	"github.com/zitadel/zitadel/internal/api/oidc/sign"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/denylist"
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
					expectBody:  validateJSONPayload(tt.args.body),
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
		signer func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error)
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					expectBody:  validateJSONPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
		{
			"webhook, JWT, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJWTPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					PayloadType: target_domain.PayloadTypeJWT,
				},
				signer: mockSigner,
			},
			res{
				body: nil,
			},
		},
		{
			"webhook, JWE, ok",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJWEPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					PayloadType:     target_domain.PayloadTypeJWE,
					EncryptionKey:   encryptionKey,
					EncryptionKeyID: encryptionKeyID,
				},
				signer: mockSigner,
			},
			res{
				body: nil,
			},
		},
		{
			"webhook, JWE no encryption key, error",
			args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				server: &callTestServer{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJWEPayload([]byte("{\"request\":{\"content\":\"request1\"}}")),
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
					PayloadType: target_domain.PayloadTypeJWE,
				},
				signer: mockSigner,
			},
			res{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respBody, err := testServer(
				t,
				tt.args.server,
				testCallTarget(
					tt.args.ctx, tt.args.info, tt.args.target,
					crypto.CreateMockEncryptionAlg(gomock.NewController(t)), tt.args.signer, &sync.Map{},
					[]denylist.AddressChecker{},
				),
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

func Test_CallTargets(t *testing.T) {
	deniedLocalhost, err := denylist.NewHostChecker("127.0.0.1")
	require.NoError(t, err)
	deniedIPs := []denylist.AddressChecker{deniedLocalhost}
	type args struct {
		ctx                    context.Context
		info                   *middleware.ContextInfoRequest
		servers                []*callTestServer
		targets                []target_domain.Target
		getActiveSigningWebKey func(*int32) execution.GetActiveSigningWebKey
	}
	type res struct {
		ret     any
		wantErr bool
	}
	tests := []struct {
		name     string
		denyList []denylist.AddressChecker
		args     args
		res      res
	}{
		{
			name: "interrupt on status",
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: true},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "continue on status",
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusInternalServerError,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: false},
				},
			},
			res: res{
				ret: requestContextInfo1.GetContent(),
			},
		},
		{
			name: "interrupt on json error",
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: []byte("just a string, not json"),
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: true},
				},
			},
			res: res{
				wantErr: true,
			},
		},
		{
			name: "continue on json error",
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: []byte("just a string, not json"),
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: false},
				}},
			res: res{
				ret: requestContextInfo1.GetContent(),
			},
		},
		{
			name: "multiple JWT/JWE targets, ok",
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJWTPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{
						TargetType:  target_domain.TargetTypeWebhook,
						PayloadType: target_domain.PayloadTypeJWT,
						Timeout:     time.Minute,
					},
					{
						TargetType:  target_domain.TargetTypeWebhook,
						PayloadType: target_domain.PayloadTypeJWE,
						Timeout:     time.Minute,
					},
				},
				getActiveSigningWebKey: testActiveSingingWebKey,
			},
			res: res{
				ret: requestContextInfo1.GetContent(),
			},
		},
		{
			name:     "block request when target in denylist",
			denyList: deniedIPs,
			args: args{
				ctx:  context.Background(),
				info: requestContextInfo1,
				servers: []*callTestServer{{
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: requestContextInfoBody2,
					statusCode:  http.StatusOK,
				}, {
					timeout:     time.Second,
					method:      http.MethodPost,
					expectBody:  validateJSONPayload(requestContextInfoBody1),
					respondBody: []byte("just a string, not json"),
					statusCode:  http.StatusOK,
				}},
				targets: []target_domain.Target{
					{InterruptOnError: false},
					{InterruptOnError: true},
				},
			},
			res: res{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var getWebKeyCalls int32
			var getActiveSigningWebKey execution.GetActiveSigningWebKey
			if tt.args.getActiveSigningWebKey != nil {
				getActiveSigningWebKey = tt.args.getActiveSigningWebKey(&getWebKeyCalls)
			}
			respBody, err := testServers(t,
				tt.args.servers,
				testCallTargets(
					tt.args.ctx, tt.args.info, tt.args.targets,
					crypto.CreateMockEncryptionAlg(gomock.NewController(t)), getActiveSigningWebKey, tt.denyList),
			)
			if tt.res.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.res.ret, respBody)
			var expectedCalls int32
			if tt.args.getActiveSigningWebKey != nil {
				expectedCalls = 1
			}
			assert.Equal(t, expectedCalls, atomic.LoadInt32(&getWebKeyCalls))
		})
	}
}

type callTestServer struct {
	method      string
	expectBody  func(*testing.T, []byte) bool
	timeout     time.Duration
	statusCode  int
	respondBody []byte
	signingKey  string
	called      bool
}

func (s *callTestServer) Called() bool {
	return s.called
}

func testServers(
	t *testing.T,
	c []*callTestServer,
	call func([]string) (interface{}, error),
) (interface{}, error) {
	urls := make([]string, len(c))
	for i := range c {
		url, closeF, _ := listen(t, c[i])
		defer closeF()
		urls[i] = url
	}
	return call(urls)
}

func testServer(
	t *testing.T,
	c *callTestServer,
	call func(string) ([]byte, error),
) ([]byte, error) {
	url, closeF, _ := listen(t, c)
	defer closeF()
	return call(url)
}

func listen(
	t *testing.T,
	c *callTestServer,
) (url string, close func(), called func() bool) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		c.called = true
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
	return server.URL, server.Close, c.Called
}

func checkRequest(t *testing.T, sent *http.Request, method string, checkExpectedBody func(*testing.T, []byte) bool, signingKey string) {
	sentBody, err := io.ReadAll(sent.Body)
	require.NoError(t, err)
	if !checkExpectedBody(t, sentBody) {
		return
	}
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
	signerOnce sign.SignerFunc,
	encrypters *sync.Map,
	actionsDenyList []denylist.AddressChecker,
) func(string) ([]byte, error) {
	return func(url string) (r []byte, err error) {
		target.Endpoint = url
		return execution.CallTarget(ctx, target, info, alg, signerOnce, encrypters, actionsDenyList)
	}
}

func testCallTargets(ctx context.Context,
	info execution.ContextInfo,
	target []target_domain.Target,
	alg crypto.EncryptionAlgorithm,
	activeSigningKey execution.GetActiveSigningWebKey,
	actionsDenyList []denylist.AddressChecker,
) func([]string) (any, error) {
	return func(urls []string) (any, error) {
		targets := make([]target_domain.Target, len(target))
		for i, t := range target {
			t.Endpoint = urls[i]
			targets[i] = t
		}
		return execution.CallTargets(ctx, targets, info, alg, activeSigningKey, actionsDenyList)
	}
}

var requestContextInfo1 = &middleware.ContextInfoRequest{
	Request: middleware.Message{Message: &structpb.Struct{
		Fields: map[string]*structpb.Value{"content": structpb.NewStringValue("request1")},
	}},
}

var requestContextInfoBody1 = []byte("{\"request\":{\"content\":\"request1\"}}")
var requestContextInfoBody2 = []byte("{\"request\":{\"content\":\"request2\"}}")

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

var (
	privateKey = func() *rsa.PrivateKey {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return privateKey
	}()
	encryptionKey = func() []byte {
		data, _ := crypto.PublicKeyToBytes(&privateKey.PublicKey)
		return data
	}()
	encryptionKeyID  = "encryption-key-id"
	signingAlgorithm = jose.RS256
)

func mockSigner(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error) {
	return sign.GetSignerOnce(mockGetActiveSigningWebKey)(ctx)
}

func mockGetActiveSigningWebKey(ctx context.Context) (*jose.JSONWebKey, error) {
	return &jose.JSONWebKey{
		Key:       privateKey,
		Algorithm: string(signingAlgorithm),
		Use:       "sig",
	}, nil
}

func testActiveSingingWebKey(getWebKeyCalls *int32) execution.GetActiveSigningWebKey {
	return func(ctx context.Context) (*jose.JSONWebKey, error) {
		atomic.AddInt32(getWebKeyCalls, 1)
		return mockGetActiveSigningWebKey(ctx)
	}
}

func validateJSONPayload(expected []byte) func(*testing.T, []byte) bool {
	return func(t *testing.T, actual []byte) bool {
		require.Equal(t, expected, actual)
		return false
	}
}

func validateJWTPayload(expected []byte) func(*testing.T, []byte) bool {
	return func(t *testing.T, actual []byte) bool {
		jws, err := jose.ParseSigned(string(actual), []jose.SignatureAlgorithm{jose.RS256})
		require.NoError(t, err)
		payload, err := jws.Verify(privateKey.Public())
		require.NoError(t, err)
		return bytes.Equal(expected, payload)
	}
}

func validateJWEPayload(expected []byte) func(*testing.T, []byte) bool {
	return func(t *testing.T, actual []byte) bool {
		parsedJWE, err := jose.ParseEncrypted(string(actual), []jose.KeyAlgorithm{jose.RSA_OAEP_256, jose.ECDH_ES_A256KW}, []jose.ContentEncryption{jose.A256GCM})
		if err != nil {
			return false
		}
		require.Equal(t, encryptionKeyID, parsedJWE.Header.KeyID)
		require.Equal(t, "JWT", parsedJWE.Header.ExtraHeaders[jose.HeaderContentType].(string))

		decryptedJWS, err := parsedJWE.Decrypt(privateKey)
		require.NoError(t, err)
		return validateJWTPayload(expected)(t, decryptedJWS)
	}
}
