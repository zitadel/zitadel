package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func TestRequestDetailsHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupContext bool
		status       int
	}{
		{
			name:         "sets request ID header on success",
			setupContext: true,
			status:       http.StatusOK,
		},
		{
			name:         "sets request ID header on error",
			setupContext: true,
			status:       http.StatusInternalServerError,
		},
		{
			name:         "sets request ID header without call context",
			setupContext: false,
			status:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			// Setup context with call duration if needed
			if tt.setupContext {
				req = req.WithContext(call.WithTimestamp(req.Context()))
			}

			rec := httptest.NewRecorder()

			// Create test handler that returns the specified status
			handler := RequestDetailsHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request ID is in context
				requestID := instrumentation.GetRequestID(r.Context())
				assert.False(t, requestID.IsNil(), "Request ID should be present in context")

				w.WriteHeader(tt.status)
			}))

			handler.ServeHTTP(rec, req)

			// Verify response
			resp := rec.Result()
			assert.Equal(t, tt.status, resp.StatusCode)

			// Verify x-request-id header is set
			requestID := resp.Header.Get(http_util.XRequestID)
			require.NotEmpty(t, requestID, "x-request-id header should be set")

			// Verify the x-request-id header is a valid xid format
			_, err := xid.FromString(requestID)
			require.NoError(t, err, "x-request-id header should be a valid xid")
		})
	}
}

func TestRequestDetailsHandler_Stability(t *testing.T) {
	// Test that the same request gets consistent request ID
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(call.WithTimestamp(req.Context()))
	rec := httptest.NewRecorder()

	var capturedID string
	handler := RequestDetailsHandler()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := instrumentation.GetRequestID(r.Context())
		require.False(t, id.IsNil(), "Request ID should be present in context")
		capturedID = id.String()
	}))

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	headerID := resp.Header.Get(http_util.XRequestID)

	// Verify the ID in context matches the ID in header
	assert.Equal(t, capturedID, headerID, "Request ID in context should match header")
}
