package types

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotify_SendPasswordlessRegistrationLink(t *testing.T) {
	type args struct {
		user    *query.NotifyUser
		origin  *http_utils.DomainCtx
		code    string
		codeID  string
		urlTmpl string
	}
	tests := []struct {
		name    string
		args    args
		want    *notifyResult
		wantErr error
	}{
		{
			name: "default URL",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:    "123",
				codeID:  "456",
				urlTmpl: "",
			},
			want: &notifyResult{
				url:                                "https://example.com/ui/login/login/passwordless/init?userID=user1&orgID=org1&codeID=456&code=123",
				messageType:                        domain.PasswordlessRegistrationMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
		{
			name: "template error",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:    "123",
				codeID:  "456",
				urlTmpl: "{{",
			},
			want:    &notifyResult{},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "template success",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:    "123",
				codeID:  "456",
				urlTmpl: "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
			},
			want: &notifyResult{
				url:                                "https://example.com/passkey/register?userID=user1&orgID=org1&codeID=456&code=123",
				messageType:                        domain.PasswordlessRegistrationMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, notify := mockNotify()
			err := notify.SendPasswordlessRegistrationLink(http_utils.WithDomainContext(context.Background(), tt.args.origin), tt.args.user, tt.args.code, tt.args.codeID, tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
