package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

func TestErrorHandler_PreservesUnauthenticatedStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/auth/v1/users/me", nil)
	resp := httptest.NewRecorder()

	errorHandler(context.Background(), runtime.NewServeMux(), jsonMarshaler, resp, req, status.Error(codes.Unauthenticated, "auth header missing"))

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), `"code":16`)
	assert.Contains(t, resp.Body.String(), `"message":"auth header missing"`)
}

// TestGatewayRecordsURIPattern guards the uri label of the HTTP metrics against the object
// IDs in the gateway's paths. A single instance serves an unbounded number of sessions,
// auth requests and users, so recording the requested path produces one metric series per
// object. Every way the gateway can answer a request has to report its route pattern.
func TestGatewayRecordsURIPattern(t *testing.T) {
	const (
		pattern    = "/v2/sessions/{session_id}"
		requestURI = "/v2/sessions/385063742120926058"
	)
	tests := []struct {
		name string
		// respond answers the request the way the gateway would in the tested case.
		respond func(ctx context.Context, mux *runtime.ServeMux, w http.ResponseWriter, r *http.Request)
		// routed tells whether the request reached a service method.
		routed bool
		want   string
	}{
		{
			name: "successful response",
			respond: func(ctx context.Context, _ *runtime.ServeMux, w http.ResponseWriter, _ *http.Request) {
				require.NoError(t, responseForwarder(ctx, w, nil))
			},
			routed: true,
			want:   pattern,
		},
		{
			name: "error response",
			respond: func(ctx context.Context, mux *runtime.ServeMux, w http.ResponseWriter, r *http.Request) {
				errorHandler(ctx, mux, jsonMarshaler, w, r, status.Error(codes.NotFound, "session not found"))
			},
			routed: true,
			want:   pattern,
		},
		{
			name: "unroutable path",
			respond: func(ctx context.Context, mux *runtime.ServeMux, w http.ResponseWriter, r *http.Request) {
				httpErrorHandler(ctx, mux, jsonMarshaler, w, r, http.StatusNotFound)
			},
			want: metrics.UnknownPath,
		},
		{
			name: "method not allowed",
			respond: func(ctx context.Context, mux *runtime.ServeMux, w http.ResponseWriter, r *http.Request) {
				httpErrorHandler(ctx, mux, jsonMarshaler, w, r, http.StatusMethodNotAllowed)
			},
			want: metrics.UnknownPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The metrics middleware seeds the request before handing it to the gateway.
			req := httptest.NewRequest(http.MethodGet, requestURI, nil)
			req = req.WithContext(metrics.WithRequestURIPattern(req.Context()))

			ctx := req.Context()
			if tt.routed {
				ctx = runtime.WithHTTPPathPattern(pattern)(ctx)
			}
			tt.respond(ctx, runtime.NewServeMux(serveMuxOptions(nil)...), httptest.NewRecorder(), req)

			assert.Equal(t, tt.want, metrics.RequestURI(req))
		})
	}
}
