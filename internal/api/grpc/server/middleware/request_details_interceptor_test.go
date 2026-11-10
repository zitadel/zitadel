package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
)

func TestRequestDetailsHandler_gRPC(t *testing.T) {
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
				id := instrumentation.GetRequestID(ctx)
				require.False(t, id.IsNil(), "Request ID should be present in context")
				capturedID = id.String()
				return tt.handler(ctx, req)
			}

			interceptor := RequestDetailsHandler()
			_, err := interceptor(ctx, &mockReq{}, mockInfo("/test"), handler)
			require.Equalf(t, tt.wantErr, err != nil, "Expected error: %v, got: %v", tt.wantErr, err)
			// Verify request ID is set in context
			require.NotEmpty(t, capturedID, "Request ID should be set in context")
		})
	}
}
