package oidc

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_oidcError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "nil",
			err:     nil,
			wantErr: nil,
		},
		{
			name:    "status err",
			err:     op.NewStatusError(io.ErrClosedPipe, http.StatusTeapot),
			wantErr: op.NewStatusError(io.ErrClosedPipe, http.StatusTeapot),
		},
		{
			name:    "oidc err",
			err:     oidc.ErrInvalidClient().WithParent(io.ErrClosedPipe),
			wantErr: oidc.ErrInvalidClient().WithParent(io.ErrClosedPipe),
		},
		{
			name: "unknown err",
			err:  io.ErrClosedPipe,
			wantErr: op.NewStatusError(
				oidc.ErrServerError().
					WithParent(io.ErrClosedPipe).
					WithDescription("Errors.Internal"),
				http.StatusInternalServerError,
			),
		},
		{
			name: "zitadel error, invalid request",
			err:  zerrors.ThrowPreconditionFailed(io.ErrClosedPipe, "TEST-123", "oopsie"),
			wantErr: op.NewStatusError(
				oidc.ErrInvalidRequest().
					WithParent(io.ErrClosedPipe).
					WithDescription("oopsie"),
				http.StatusBadRequest,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := oidcError(t.Context(), tt.err)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_asInvalidGrant(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "nil",
			err:     nil,
			wantErr: nil,
		},
		{
			name:    "unrelated error is passed as-is",
			err:     zerrors.ThrowInternal(io.ErrClosedPipe, "TEST-123", "oopsie"),
			wantErr: zerrors.ThrowInternal(io.ErrClosedPipe, "TEST-123", "oopsie"),
		},
		{
			name:    "refresh token invalid",
			err:     zerrors.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid"),
			wantErr: oidc.ErrInvalidGrant().WithDescription("test"),
		},
		{
			name:    "v1 refresh token not found",
			err:     zerrors.ThrowNotFound(nil, "EVENT-5Bm9s", "Errors.User.RefreshToken.Invalid"),
			wantErr: oidc.ErrInvalidGrant().WithDescription("test"),
		},
		{
			name:    "op invalid refresh token",
			err:     op.ErrInvalidRefreshToken,
			wantErr: oidc.ErrInvalidGrant().WithDescription("test"),
		},
		{
			name:    "authorization code invalid",
			err:     zerrors.ThrowInvalidArgument(nil, "OIDC-ahLi2", "Errors.User.Code.Invalid"),
			wantErr: oidc.ErrInvalidGrant().WithDescription("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := asInvalidGrant(t.Context(), tt.err, "test")
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
