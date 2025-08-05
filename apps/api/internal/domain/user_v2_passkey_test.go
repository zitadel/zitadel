package domain

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRenderPasskeyURLTemplate(t *testing.T) {
	type args struct {
		tmpl   string
		userID string
		orgID  string
		codeID string
		code   string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr error
	}{
		{
			name: "parse error",
			args: args{
				tmpl: "{{",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "success",
			args: args{
				tmpl:   "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
				userID: "user1",
				orgID:  "org1",
				codeID: "99",
				code:   "123",
			},
			wantW: "https://example.com/passkey/register?userID=user1&orgID=org1&codeID=99&code=123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := RenderPasskeyURLTemplate(w, tt.args.tmpl, tt.args.userID, tt.args.orgID, tt.args.codeID, tt.args.code)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}
