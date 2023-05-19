package webauthn

import (
	"context"
	"testing"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestConfig_serverFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *webauthn.WebAuthn
		wantErr error
	}{
		{
			name:    "webauthn error",
			args:    args{context.Background()},
			wantErr: caos_errs.ThrowInternal(nil, "WEBAU-UX9ta", "Errors.User.WebAuthN.ServerConfig"),
		},
		{
			name: "success",
			args: args{authz.WithRequestedDomain(context.Background(), "example.com")},
			want: &webauthn.WebAuthn{
				Config: &webauthn.Config{
					RPDisplayName: "DisplayName",
					RPID:          "example.com",
					RPOrigin:      "https://example.com",
					Timeout:       60000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Config{
				DisplayName:    "DisplayName",
				ExternalSecure: true,
			}
			got, err := w.serverFromContext(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
