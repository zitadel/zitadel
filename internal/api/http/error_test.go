package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestZitadelErrorToHTTPStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantOk         bool
	}{
		{
			name: "no error",
			args: args{
				err: nil,
			},
			wantStatusCode: http.StatusOK,
			wantOk:         true,
		},
		{
			name: "wrapped already exists",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowAlreadyExists(nil, "id", "message")),
			},
			wantStatusCode: http.StatusConflict,
			wantOk:         true,
		},
		{
			name: "wrapped deadline exceeded",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowDeadlineExceeded(nil, "id", "message")),
			},
			wantStatusCode: http.StatusGatewayTimeout,
			wantOk:         true,
		},
		{
			name: "wrapped internal",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowInternal(nil, "id", "message")),
			},
			wantStatusCode: http.StatusInternalServerError,
			wantOk:         true,
		},
		{
			name: "wrapped invalid argument",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowInvalidArgument(nil, "id", "message")),
			},
			wantStatusCode: http.StatusBadRequest,
			wantOk:         true,
		},
		{
			name: "wrapped not found",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowNotFound(nil, "id", "message")),
			},
			wantStatusCode: http.StatusNotFound,
			wantOk:         true,
		},
		{
			name: "wrapped permission denied",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowPermissionDenied(nil, "id", "message")),
			},
			wantStatusCode: http.StatusForbidden,
			wantOk:         true,
		},
		{
			name: "wrapped precondition failed",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowPreconditionFailed(nil, "id", "message")),
			},
			wantStatusCode: http.StatusBadRequest,
			wantOk:         true,
		},
		{
			name: "wrapped unauthenticated",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowUnauthenticated(nil, "id", "message")),
			},
			wantStatusCode: http.StatusUnauthorized,
			wantOk:         true,
		},
		{
			name: "wrapped unavailable",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowUnavailable(nil, "id", "message")),
			},
			wantStatusCode: http.StatusServiceUnavailable,
			wantOk:         true,
		},
		{
			name: "wrapped unimplemented",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowUnimplemented(nil, "id", "message")),
			},
			wantStatusCode: http.StatusNotImplemented,
			wantOk:         true,
		},
		{
			name: "wrapped resource exhausted",
			args: args{
				err: fmt.Errorf("wrapped %w", zerrors.ThrowResourceExhausted(nil, "id", "message")),
			},
			wantStatusCode: http.StatusTooManyRequests,
			wantOk:         true,
		},
		{
			name: "no caos/zitadel error",
			args: args{
				err: errors.New("error"),
			},
			wantStatusCode: http.StatusInternalServerError,
			wantOk:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatusCode, gotOk := ZitadelErrorToHTTPStatusCode(tt.args.err)
			if gotStatusCode != tt.wantStatusCode {
				t.Errorf("ZitadelErrorToHTTPStatusCode() gotStatusCode = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ZitadelErrorToHTTPStatusCode() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestHTTPStatusCodeToZitadelError(t *testing.T) {
	type args struct {
		statusCode int
		id         string
		message    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "StatusOK",
			args: args{
				statusCode: http.StatusOK,
			},
			wantErr: nil,
		},
		{
			name: "StatusConflict",
			args: args{
				statusCode: http.StatusConflict,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowAlreadyExists(nil, "id", "message"),
		},
		{
			name: "StatusGatewayTimeout",
			args: args{
				statusCode: http.StatusGatewayTimeout,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowDeadlineExceeded(nil, "id", "message"),
		},
		{
			name: "StatusInternalServerError",
			args: args{
				statusCode: http.StatusInternalServerError,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowInternal(nil, "id", "message"),
		},
		{
			name: "StatusBadRequest",
			args: args{
				statusCode: http.StatusBadRequest,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "id", "message"),
		},
		{
			name: "StatusNotFound",
			args: args{
				statusCode: http.StatusNotFound,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowNotFound(nil, "id", "message"),
		},
		{
			name: "StatusForbidden",
			args: args{
				statusCode: http.StatusForbidden,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "id", "message"),
		},
		{
			name: "StatusUnauthorized",
			args: args{
				statusCode: http.StatusUnauthorized,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowUnauthenticated(nil, "id", "message"),
		},
		{
			name: "StatusServiceUnavailable",
			args: args{
				statusCode: http.StatusServiceUnavailable,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowUnavailable(nil, "id", "message"),
		},
		{
			name: "StatusNotImplemented",
			args: args{
				statusCode: http.StatusNotImplemented,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowUnimplemented(nil, "id", "message"),
		},
		{
			name: "StatusTooManyRequests",
			args: args{
				statusCode: http.StatusTooManyRequests,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowResourceExhausted(nil, "id", "message"),
		},
		{
			name: "Unknown",
			args: args{
				statusCode: 1000,
				id:         "id",
				message:    "message",
			},
			wantErr: zerrors.ThrowUnknown(nil, "id", "message"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HTTPStatusCodeToZitadelError(tt.args.statusCode, tt.args.id, tt.args.message)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
