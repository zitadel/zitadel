package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/errors"
)

func TestStripIDPrefix(t *testing.T) {
	type args struct {
		authRequestID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "empty",
			args: args{
				authRequestID: "",
			},
			wantErr: errors.ThrowInvalidArgumentf(nil, "OIDC-Aumu8", "auth_request_id wrong version, missing %s prefix", IDPrefix),
		},
		{
			name: "no prefix",
			args: args{
				authRequestID: "123",
			},
			wantErr: errors.ThrowInvalidArgumentf(nil, "OIDC-Aumu8", "auth_request_id wrong version, missing %s prefix", IDPrefix),
		},
		{
			name: "success",
			args: args{
				authRequestID: "V2_123",
			},
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StripIDPrefix(tt.args.authRequestID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
