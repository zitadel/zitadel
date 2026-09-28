package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

// TestMetricsHandler_reportedPattern covers routers that report their route pattern
// themselves, such as the grpc-gateway. Without the context prepared by the handler their
// report is lost and the requested path is recorded, which adds a metric series per object.
func TestMetricsHandler_reportedPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name: "router reports no pattern, path is recorded",
			want: "/v2/sessions/385063742120926058",
		},
		{
			name:    "router reports a pattern, pattern is recorded",
			pattern: "/v2/sessions/{session_id}",
			want:    "/v2/sessions/{session_id}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var served *http.Request
			handler := MetricsHandler([]metrics.MetricType{metrics.MetricTypeStatusCode})(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					served = r
					if tt.pattern != "" {
						metrics.SetRequestURIPattern(r.Context(), tt.pattern)
					}
					w.WriteHeader(http.StatusNoContent)
				}),
			)

			handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/v2/sessions/385063742120926058", nil))

			require.NotNil(t, served)
			assert.Equal(t, tt.want, metrics.RequestURI(served))
		})
	}
}

// TestMetricsHandler_chiRoutePattern covers the routers that leave the pattern to chi, such
// as the OIDC endpoints. Their client configuration endpoints are templated on the client
// ID and any unknown path below their prefixes ends up here as well.
func TestMetricsHandler_chiRoutePattern(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		requestURI string
		want       string
	}{
		{
			name:       "static endpoint",
			method:     http.MethodPost,
			requestURI: "/oauth/v2/token",
			want:       "/oauth/v2/token",
		},
		{
			name:       "static endpoint with query",
			method:     http.MethodGet,
			requestURI: "/oauth/v2/authorize?client_id=1234&state=abcd",
			want:       "/oauth/v2/authorize",
		},
		{
			name:       "client configuration endpoint",
			method:     http.MethodGet,
			requestURI: "/oauth/v2/register/307534931347210241",
			want:       "/oauth/v2/register/{client_id}",
		},
		{
			name:       "unknown path",
			method:     http.MethodGet,
			requestURI: "/oauth/v2/junk",
			want:       metrics.UnknownPath,
		},
		{
			name:       "known path, unsupported method",
			method:     http.MethodDelete,
			requestURI: "/oauth/v2/token",
			want:       metrics.UnknownPath,
		},
	}

	// served is the request as the routed handler saw it. Reading the recorded URI from it
	// after the router returned gives the value the metrics handler recorded, because the
	// two share the context the metrics handler prepared.
	var served *http.Request
	respond := func(w http.ResponseWriter, r *http.Request) {
		served = r
		w.WriteHeader(http.StatusNoContent)
	}

	// Mirrors how the OIDC server nests its routes in the metrics handler.
	router := chi.NewRouter()
	router.Use(MetricsHandler([]metrics.MetricType{metrics.MetricTypeStatusCode}))
	router.Post("/oauth/v2/token", respond)
	router.Get("/oauth/v2/authorize", respond)
	router.Get("/oauth/v2/register/{client_id}", respond)
	router.NotFound(respond)
	router.MethodNotAllowed(respond)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			served = nil

			router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(tt.method, tt.requestURI, nil))

			require.NotNil(t, served)
			assert.Equal(t, tt.want, metrics.RequestURI(served))
		})
	}
}
