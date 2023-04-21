package types

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

func TestNotify_SendEmailVerificationCode(t *testing.T) {
	type res struct {
		url                                string
		args                               map[string]interface{}
		messageType                        string
		allowUnverifiedNotificationChannel bool
	}
	notify := func(dst *res) Notify {
		return func(
			url string,
			args map[string]interface{},
			messageType string,
			allowUnverifiedNotificationChannel bool,
		) error {
			dst.url = url
			dst.args = args
			dst.messageType = messageType
			dst.allowUnverifiedNotificationChannel = allowUnverifiedNotificationChannel
			return nil
		}
	}

	type args struct {
		user    *query.NotifyUser
		origin  string
		code    string
		urlTmpl *string
	}
	tests := []struct {
		name    string
		args    args
		want    *res
		wantErr error
	}{
		{
			name: "default URL",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  "https://example.com",
				code:    "123",
				urlTmpl: nil,
			},
			want: &res{
				url:                                "https://example.com/ui/login/mail/verification?userID=user1&code=123&orgID=org1",
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
				origin:  "https://example.com",
				code:    "123",
				urlTmpl: gu.Ptr("{{"),
			},
			want:    &res{},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "USERv2-ooD8p", "Errors.User.V2.Email.InvalidURLTemplate"),
		},
		{
			name: "template success",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  "https://example.com",
				code:    "123",
				urlTmpl: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
			},
			want: &res{
				url:                                "https://example.com/email/verify?userID=user1&code=123&orgID=org1",
				args:                               map[string]interface{}{"Code": "123"},
				messageType:                        domain.VerifyEmailMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(res)
			err := notify(got).SendEmailVerificationCode(tt.args.user, tt.args.origin, tt.args.code, tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
