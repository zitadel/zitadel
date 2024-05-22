package domain

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestOIDCErrorReasonFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want OIDCErrorReason
	}{
		{
			name: "invalid request",
			err:  oidc.ErrInvalidRequest().WithDescription("foo"),
			want: OIDCErrorReasonInvalidRequest,
		},
		{
			name: "invalid grant",
			err:  oidc.ErrInvalidGrant().WithDescription("foo"),
			want: OIDCErrorReasonInvalidGrant,
		},
		{
			name: "precondition failed",
			err:  zerrors.ThrowPreconditionFailed(nil, "123", "bar"),
			want: OIDCErrorReasonAccessDenied,
		},
		{
			name: "internal",
			err:  zerrors.ThrowInternal(nil, "123", "bar"),
			want: OIDCErrorReasonServerError,
		},
		{
			name: "unspecified",
			err:  io.ErrClosedPipe,
			want: OIDCErrorReasonUnspecified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OIDCErrorReasonFromError(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
