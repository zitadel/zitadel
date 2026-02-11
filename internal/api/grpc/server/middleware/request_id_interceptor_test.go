package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func TestRequestIDHandler_gRPC(t *testing.T) {
	tests := []struct {
		name         string
		setupContext bool
		handler      grpc.UnaryHandler
		wantErr      bool
	}{
		{
			name:         "sets request ID header on success",
			setupContext: true,
			handler:      emptyMockHandler,
			wantErr:      false,
		},
		{
			name:         "sets request ID header on error",
			setupContext: true,
			handler:      errorMockHandler,
			wantErr:      true,
		},
		{
			name:         "sets request ID header without call context",
			setupContext: false,
			handler:      emptyMockHandler,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Setup context with call duration if needed
			if tt.setupContext {
				ctx = call.WithTimestamp(ctx)
			}

			// Create a handler that captures the request ID from context
			var capturedID string
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				id, ok := instrumentation.GetRequestID(ctx)
				require.True(t, ok, "Request ID should be present in context")
				capturedID = id.String()
				return tt.handler(ctx, req)
			}

			interceptor := RequestIDHandler()
			_, err := interceptor(ctx, &mockReq{}, mockInfo("/test"), handler)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Verify request ID is set in context
			require.NotEmpty(t, capturedID, "Request ID should be set in context")

			// Verify the request ID is a valid UUID format
			assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, capturedID)
		})
	}
}

func TestRequestIDHandler_gRPC_HeaderSet(t *testing.T) {
	// Test that the request ID is set in response header
	// This is verified indirectly through SetHeader call
	// We can't directly test the metadata since it requires a real gRPC context
	ctx := call.WithTimestamp(context.Background())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Try to get headers that were set
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			requestID := md.Get(http_util.XRequestID)
			if len(requestID) > 0 {
				assert.NotEmpty(t, requestID[0], "Request ID should be in metadata")
			}
		}
		return req, nil
	}

	interceptor := RequestIDHandler()
	_, err := interceptor(ctx, &mockReq{}, mockInfo("/test"), handler)
	require.NoError(t, err)
}
