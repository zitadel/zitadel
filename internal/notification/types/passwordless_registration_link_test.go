package types

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotify_SendPasswordlessRegistrationLink(t *testing.T) {
	type args struct {
		user        *query.NotifyUser
		origin      string
		code        string
		codeID      string
		urlTmpl     string
		loginPolicy *query.LoginPolicy
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
				origin:  "https://example.com",
				code:    "123",
				codeID:  "456",
				urlTmpl: "",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
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
				origin:  "https://example.com",
				code:    "123",
				codeID:  "456",
				urlTmpl: "{{",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
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
				origin:  "https://example.com",
				code:    "123",
				codeID:  "456",
				urlTmpl: "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: false,
				},
			},
			want: &notifyResult{
				url:                                "https://example.com/passkey/register?userID=user1&orgID=org1&codeID=456&code=123",
				messageType:                        domain.PasswordlessRegistrationMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
		{
			name: "use default uri for url link success",
			args: args{
				user: &query.NotifyUser{
					ID:            "user1",
					ResourceOwner: "org1",
				},
				origin:  "https://example.com",
				code:    "123",
				codeID:  "456",
				urlTmpl: "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.OrgID}}&codeID={{.CodeID}}&code={{.Code}}",
				loginPolicy: &query.LoginPolicy{
					AllowUsernamePassword:      true,
					AllowRegister:              true,
					AllowExternalIDPs:          true,
					ForceMFA:                   true,
					ForceMFALocalOnly:          true,
					PasswordlessType:           domain.PasswordlessTypeAllowed,
					HidePasswordReset:          true,
					IgnoreUnknownUsernames:     true,
					AllowDomainDiscovery:       true,
					DisableLoginWithEmail:      true,
					DisableLoginWithPhone:      true,
					DefaultRedirectURI:         "https://example.com",
					PasswordCheckLifetime:      database.Duration(time.Hour),
					ExternalLoginCheckLifetime: database.Duration(time.Minute),
					MFAInitSkipLifetime:        database.Duration(time.Millisecond),
					SecondFactorCheckLifetime:  database.Duration(time.Microsecond),
					MultiFactorCheckLifetime:   database.Duration(time.Nanosecond),
					SecondFactors: []domain.SecondFactorType{
						domain.SecondFactorTypeTOTP,
						domain.SecondFactorTypeU2F,
						domain.SecondFactorTypeOTPEmail,
						domain.SecondFactorTypeOTPSMS,
					},
					MultiFactors: []domain.MultiFactorType{
						domain.MultiFactorTypeU2FWithPIN,
					},
					IsDefault: true,
					UseDefaultRedirectUriForNotificationLinks: true,
				},
			},
			want: &notifyResult{
				url:                                "https://example.com",
				messageType:                        domain.PasswordlessRegistrationMessageType,
				allowUnverifiedNotificationChannel: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, notify := mockNotify()
			err := notify.SendPasswordlessRegistrationLink(http_utils.WithComposedOrigin(context.Background(), tt.args.origin), tt.args.user, tt.args.code, tt.args.codeID, tt.args.urlTmpl, tt.args.loginPolicy)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
