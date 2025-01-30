package gerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"

	commandErrors "github.com/zitadel/zitadel/internal/command/errors"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func TestCaosToGRPCError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no error",
			args{},
			false,
		},
		{
			"unknown error",
			args{errors.New("unknown")},
			true,
		},
		{
			"caos error",
			args{zerrors.ThrowInternal(nil, "", "message")},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ZITADELToGRPCError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("ZITADELToGRPCError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getErrorInfo(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		err    error
		result protoadapt.MessageV1
	}{
		{
			name:   "parent error nil, return message.ErrorDetail{}",
			id:     "id",
			key:    "key",
			result: &message.ErrorDetail{Id: "id", Message: "key"},
		},
		{
			name:   "parent error nil not commandErrors.WrongPasswordError{}, return message.ErrorDetail{}",
			id:     "id",
			key:    "key",
			err:    fmt.Errorf("normal error not commandErrors.WrongPasswordError{}"),
			result: &message.ErrorDetail{Id: "id", Message: "key"},
		},
		{
			name:   "parent error not nil type commandErrors.WrongPasswordError{}, return message.CredentialsCheckError{}",
			id:     "id",
			key:    "key",
			err:    &commandErrors.WrongPasswordError{FailedAttempts: 22},
			result: &message.CredentialsCheckError{Id: "id", Message: "key", FailedAttempts: 22},
		},
		{
			name: "parent error not nil wrapped commandErrors.WrongPasswordError{}, return message.CredentialsCheckError{}",
			id:   "id",
			key:  "key",
			err: func() error {
				err := fmt.Errorf("normal error")
				wpa := &commandErrors.WrongPasswordError{FailedAttempts: 26}
				return fmt.Errorf("%w: %w", err, wpa)
			}(),
			result: &message.CredentialsCheckError{Id: "id", Message: "key", FailedAttempts: 26},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorInfo := getErrorInfo(tt.id, tt.key, tt.err)
			assert.Equal(t, tt.result, errorInfo)
		})
	}
}

func Test_Extract(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantC   codes.Code
		wantMsg string
		wantID  string
		wantOk  bool
	}{
		{
			"already exists",
			args{zerrors.ThrowAlreadyExists(nil, "id", "already exists")},
			codes.AlreadyExists,
			"already exists",
			"id",
			true,
		},
		{
			"deadline exceeded",
			args{zerrors.ThrowDeadlineExceeded(nil, "id", "deadline exceeded")},
			codes.DeadlineExceeded,
			"deadline exceeded",
			"id",
			true,
		},
		{
			"internal error",
			args{zerrors.ThrowInternal(nil, "id", "internal error")},
			codes.Internal,
			"internal error",
			"id",
			true,
		},
		{
			"invalid argument",
			args{zerrors.ThrowInvalidArgument(nil, "id", "invalid argument")},
			codes.InvalidArgument,
			"invalid argument",
			"id",
			true,
		},
		{
			"not found",
			args{zerrors.ThrowNotFound(nil, "id", "not found")},
			codes.NotFound,
			"not found",
			"id",
			true,
		},
		{
			"permission denied",
			args{zerrors.ThrowPermissionDenied(nil, "id", "permission denied")},
			codes.PermissionDenied,
			"permission denied",
			"id",
			true,
		},
		{
			"precondition failed",
			args{zerrors.ThrowPreconditionFailed(nil, "id", "precondition failed")},
			codes.FailedPrecondition,
			"precondition failed",
			"id",
			true,
		},
		{
			"unauthenticated",
			args{zerrors.ThrowUnauthenticated(nil, "id", "unauthenticated")},
			codes.Unauthenticated,
			"unauthenticated",
			"id",
			true,
		},
		{
			"unavailable",
			args{zerrors.ThrowUnavailable(nil, "id", "unavailable")},
			codes.Unavailable,
			"unavailable",
			"id",
			true,
		},
		{
			"unimplemented",
			args{zerrors.ThrowUnimplemented(nil, "id", "unimplemented")},
			codes.Unimplemented,
			"unimplemented",
			"id",
			true,
		},
		{
			"exhausted",
			args{zerrors.ThrowResourceExhausted(nil, "id", "exhausted")},
			codes.ResourceExhausted,
			"exhausted",
			"id",
			true,
		},
		{
			"unknown",
			args{errors.New("unknown")},
			codes.Unknown,
			"unknown",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotMsg, gotID, gotOk := ExtractZITADELError(tt.args.err)
			if gotC != tt.wantC {
				t.Errorf("extract() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("extract() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotID != tt.wantID {
				t.Errorf("extract() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotOk != tt.wantOk {
				t.Errorf("extract() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
