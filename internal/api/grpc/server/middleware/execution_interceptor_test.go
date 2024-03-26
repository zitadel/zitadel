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
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type mockExecutionQueries struct {
	executionTargets *query.ExecutionTargets
	etError          error
	targets          *query.Targets
	tError           error
}

func (e *mockExecutionQueries) ExecutionTargetsRequestResponse(_ context.Context, _, _, _ string) (execution *query.ExecutionTargets, err error) {
	return e.executionTargets, e.etError
}

func (e *mockExecutionQueries) SearchTargets(_ context.Context, _ *query.TargetSearchQueries) (targets *query.Targets, err error) {
	return e.targets, e.tError
}

type mockContentRequest struct {
	Content string
}

func newMockContentRequest(content string) *mockContentRequest {
	return &mockContentRequest{
		Content: content,
	}
}

func newMockContextInfoRequest(fullMethod, request string) *execution.ContextInfo {
	return &execution.ContextInfo{
		FullMethod: fullMethod,
		Request:    newMockContentRequest(request),
	}
}

func newMockContextInfoResponse(fullMethod, request, response string) *execution.ContextInfo {
	return &execution.ContextInfo{
		FullMethod: fullMethod,
		Request:    newMockContentRequest(request),
		Response:   newMockContentRequest(response),
	}
}

func Test_executeTargetsForGRPCFullMethod_request(t *testing.T) {
	type target struct {
		reqBody    *execution.ContextInfo
		sleep      time.Duration
		statusCode int
		respBody   interface{}
	}
	type args struct {
		ctx context.Context

		queries    *mockExecutionQueries
		targets    []target
		fullMethod string
		req        interface{}
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
			"target, not found",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					etError: zerrors.ThrowNotFound(nil, "error", "NotFound"),
				},
				req: newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, not existing",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{},
					},
				},
				req: newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, targets not found",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					tError: zerrors.ThrowNotFound(nil, "error", "NotFound"),
				},
				req: newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, targets empty",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{},
					},
				},
				req: newMockContentRequest("request"),
			},
			res{
				want: newMockContentRequest("request"),
			},
		},
		{
			"target, not found",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
					},
				},
				targets: []target{},
				req:     newMockContentRequest("content"),
			},
			res{
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"target, error without interrupt",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:         "target",
								TargetType: domain.TargetTypeRequestResponse,
								Timeout:    time.Minute,
							},
						},
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
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"target, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Second,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"target, wrong request",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Second,
								InterruptOnError: true,
							},
						},
					},
				},
				targets: []target{
					{reqBody: newMockContextInfoRequest("/service/method", "wrong")},
				},
				req: newMockContentRequest("content"),
			},
			res{
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"target, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:         "target",
								TargetType: domain.TargetTypeRequestResponse,
								Timeout:    time.Second,
								Async:      true,
							},
						},
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
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:         "target",
								TargetType: domain.TargetTypeRequestResponse,
								Timeout:    time.Minute,
								Async:      true,
							},
						},
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
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeWebhook,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"webhook, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeWebhook,
								Timeout:          time.Second,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content"),
				wantErr: true,
			},
		},
		{
			"webhook, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeWebhook,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target1", "target2", "target3"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target1",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							}, {
								ID:               "target2",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							}, {
								ID:               "target3",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content1"),
				wantErr: true,
			},
		},
		{
			"with includes, timeout",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target1", "target2", "target3"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target1",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							}, {
								ID:               "target2",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Second,
								InterruptOnError: true,
							}, {
								ID:               "target3",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Second,
								InterruptOnError: true,
							},
						},
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
				want:    newMockContentRequest("content1"),
				wantErr: true,
			},
		},
		{
			"with includes, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target1", "target2", "target3"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target1",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							}, {
								ID:               "target2",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							}, {
								ID:               "target3",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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
						statusCode: http.StatusOK,
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
				want: newMockContentRequest("content3"),
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

				tt.args.queries.targets.Targets[i].URL = url
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForRequest(
				tt.args.ctx,
				tt.args.queries,
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
		reqBody    *execution.ContextInfo
		sleep      time.Duration
		statusCode int
		respBody   interface{}
	}
	type args struct {
		ctx context.Context

		queries    *mockExecutionQueries
		targets    []target
		fullMethod string
		req        interface{}
		resp       interface{}
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
			"target, targets not found",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					tError: zerrors.ThrowNotFound(nil, "error", "NotFound"),
				},
				req:  newMockContentRequest("request"),
				resp: newMockContentRequest("response"),
			},
			res{
				want: newMockContentRequest("response"),
			},
		},
		{
			"target, not found",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					etError: zerrors.ThrowNotFound(nil, "error", "NotFound"),
				},
				req:  newMockContentRequest("request"),
				resp: newMockContentRequest("response"),
			},
			res{
				want: newMockContentRequest("response"),
			},
		},
		{
			"target, ok",
			args{
				ctx:        context.Background(),
				fullMethod: "/service/method",
				queries: &mockExecutionQueries{
					executionTargets: &query.ExecutionTargets{
						ID:      "/zitadel.session.v2beta.SessionService/SetSession",
						Targets: []string{"target"},
					},
					targets: &query.Targets{
						Targets: []*query.Target{
							{
								ID:               "target",
								TargetType:       domain.TargetTypeRequestResponse,
								Timeout:          time.Minute,
								InterruptOnError: true,
							},
						},
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

				tt.args.queries.targets.Targets[i].URL = url
				closeFuncs[i] = closeF
			}

			resp, err := executeTargetsForResponse(
				tt.args.ctx,
				tt.args.queries,
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
