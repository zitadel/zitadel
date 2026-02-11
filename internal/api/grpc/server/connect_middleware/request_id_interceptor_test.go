package connect_middleware

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRequestIDHandler_Connect(t *testing.T) {
	tests := []struct {
		name         string
		setupContext bool
		hasError     bool
		wantErr      bool
	}{
		{
			name:         "sets request ID header on success",
			setupContext: true,
			hasError:     false,
			wantErr:      false,
		},
		{
			name:         "sets request ID trailer on error",
			setupContext: true,
			hasError:     true,
			wantErr:      true,
		},
		{
			name:         "sets request ID header without call context",
			setupContext: false,
			hasError:     false,
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
			
			// Create test request
			req := &mockReq[struct{}]{
				procedure: "/test",
				header:    http.Header{},
			}
			
			// Create a mock handler
			var capturedID string
			handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
				id, ok := instrumentation.GetRequestID(ctx)
				require.True(t, ok, "Request ID should be present in context")
				capturedID = id.String()
				
				if tt.hasError {
					return nil, connect.NewError(connect.CodeInternal, zerrors.ThrowInternal(nil, "TEST", "test error"))
				}
				
				resp := &connect.Response[struct{}]{}
				return resp, nil
			}
			
			interceptor := RequestIDHandler()
			unaryFunc := interceptor(handler)
			resp, err := unaryFunc(ctx, req)
			
			if tt.wantErr {
				require.Error(t, err)
				
				// Verify request ID is in error metadata (trailer)
				var connectErr *connect.Error
				if assert.ErrorAs(t, err, &connectErr) {
					requestID := connectErr.Meta().Get(http_util.XRequestID)
					require.NotEmpty(t, requestID, "Request ID should be set in error metadata")
					assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, requestID)
					assert.Equal(t, capturedID, requestID, "Request ID in context should match error metadata")
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				
				// Verify request ID is set in response header
				requestID := resp.Header().Get(http_util.XRequestID)
				require.NotEmpty(t, requestID, "Request ID should be set in response header")
				assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, requestID)
				assert.Equal(t, capturedID, requestID, "Request ID in context should match header")
			}
			
			// Verify request ID is set in context
			require.NotEmpty(t, capturedID, "Request ID should be set in context")
		})
	}
}

func TestRequestIDHandler_Connect_Stability(t *testing.T) {
	// Test that the same request gets consistent request ID
	ctx := call.WithTimestamp(context.Background())
	req := &mockReq[struct{}]{
		procedure: "/test",
		header:    http.Header{},
	}
	
	var capturedID string
	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		id, ok := instrumentation.GetRequestID(ctx)
		require.True(t, ok, "Request ID should be present in context")
		capturedID = id.String()
		resp := &connect.Response[struct{}]{}
		return resp, nil
	}
	
	interceptor := RequestIDHandler()
	unaryFunc := interceptor(handler)
	resp, err := unaryFunc(ctx, req)
	
	require.NoError(t, err)
	require.NotNil(t, resp)
	
	headerID := resp.Header().Get(http_util.XRequestID)
	
	// Verify the ID in context matches the ID in header
	assert.Equal(t, capturedID, headerID, "Request ID in context should match header")
}
