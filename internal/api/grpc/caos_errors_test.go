package grpc

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
)

func TestCaosToGRPCError(t *testing.T) {
	type args struct {
		err        error
		translator *i18n.Translator
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
			args{errors.New("unknown"), nil},
			true,
		},
		{
			"caos error",
			args{caos_errs.ThrowInternal(nil, "", "message"), nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CaosToGRPCError(context.Background(), tt.args.err, tt.args.translator); (err != nil) != tt.wantErr {
				t.Errorf("CaosToGRPCError() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		wantOk  bool
	}{
		{
			"already exists",
			args{caos_errs.ThrowAlreadyExists(nil, "id", "already exists")},
			codes.AlreadyExists,
			"already exists",
			true,
		},
		{
			"deadline exceeded",
			args{caos_errs.ThrowDeadlineExceeded(nil, "id", "deadline exceeded")},
			codes.DeadlineExceeded,
			"deadline exceeded",
			true,
		},
		{
			"internal error",
			args{caos_errs.ThrowInternal(nil, "id", "internal error")},
			codes.Internal,
			"internal error",
			true,
		},
		{
			"invalid argument",
			args{caos_errs.ThrowInvalidArgument(nil, "id", "invalid argument")},
			codes.InvalidArgument,
			"invalid argument",
			true,
		},
		{
			"not found",
			args{caos_errs.ThrowNotFound(nil, "id", "not found")},
			codes.NotFound,
			"not found",
			true,
		},
		{
			"permission denied",
			args{caos_errs.ThrowPermissionDenied(nil, "id", "permission denied")},
			codes.PermissionDenied,
			"permission denied",
			true,
		},
		{
			"precondition failed",
			args{caos_errs.ThrowPreconditionFailed(nil, "id", "precondition failed")},
			codes.FailedPrecondition,
			"precondition failed",
			true,
		},
		{
			"unauthenticated",
			args{caos_errs.ThrowUnauthenticated(nil, "id", "unauthenticated")},
			codes.Unauthenticated,
			"unauthenticated",
			true,
		},
		{
			"unavailable",
			args{caos_errs.ThrowUnavailable(nil, "id", "unavailable")},
			codes.Unavailable,
			"unavailable",
			true,
		},
		{
			"unimplemented",
			args{caos_errs.ThrowUnimplemented(nil, "id", "unimplemented")},
			codes.Unimplemented,
			"unimplemented",
			true,
		},
		{
			"unknown",
			args{errors.New("unknown")},
			codes.Unknown,
			"unknown",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotMsg, gotOk := Extract(tt.args.err)
			if gotC != tt.wantC {
				t.Errorf("extract() gotC = %v, want %v", gotC, tt.wantC)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("extract() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotOk != tt.wantOk {
				t.Errorf("extract() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
