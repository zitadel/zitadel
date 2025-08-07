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
			err := oidcError(tt.err)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
