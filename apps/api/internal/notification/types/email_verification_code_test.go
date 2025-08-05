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

func TestNotify_SendEmailVerificationCode(t *testing.T) {
	type args struct {
		user          *query.NotifyUser
		origin        *http_utils.DomainCtx
		code          string
		urlTmpl       string
		authRequestID string
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
				origin:        &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:          "123",
				urlTmpl:       "",
				authRequestID: "authRequestID",
			},
			want: &notifyResult{
				url:                                "https://example.com/ui/login/mail/verification?authRequestID=authRequestID&code=123&orgID=org1&userID=user1",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
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
				origin:        &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:          "123",
				urlTmpl:       "{{",
				authRequestID: "authRequestID",
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
				origin:        &http_utils.DomainCtx{InstanceHost: "example.com", Protocol: "https"},
				code:          "123",
				urlTmpl:       "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				authRequestID: "authRequestID",
			},
			want: &notifyResult{
				url:                                "https://example.com/email/verify?userID=user1&code=123&orgID=org1",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, notify := mockNotify()
			err := notify.SendEmailVerificationCode(http_utils.WithDomainContext(context.Background(), tt.args.origin), tt.args.user, tt.args.code, tt.args.urlTmpl, tt.args.authRequestID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
