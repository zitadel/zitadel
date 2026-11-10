package command

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
)

func TestSessionWriteModel_AuthMethodTypes(t *testing.T) {
	type fields struct {
		PasswordCheckedAt     time.Time
		IntentCheckedAt       time.Time
		WebAuthNCheckedAt     time.Time
		WebAuthNUserVerified  bool
		TOTPCheckedAt         time.Time
		OTPSMSCheckedAt       time.Time
		OTPEmailCheckedAt     time.Time
		RecoveryCodeCheckedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   []domain.UserAuthMethodType
	}{
		{
			name: "password",
			fields: fields{
				PasswordCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypePassword,
			},
		},
		{
			name: "passwordless",
			fields: fields{
				WebAuthNCheckedAt:    testNow,
				WebAuthNUserVerified: true,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypePasswordless,
			},
		},
		{
			name: "u2f",
			fields: fields{
				WebAuthNCheckedAt:    testNow,
				WebAuthNUserVerified: false,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeU2F,
			},
		},
		{
			name: "intent",
			fields: fields{
				IntentCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeIDP,
			},
		},
		{
			name: "totp",
			fields: fields{
				TOTPCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeTOTP,
			},
		},
		{
			name: "otp sms",
			fields: fields{
				OTPSMSCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeOTPSMS,
			},
		},
		{
			name: "otp email",
			fields: fields{
				OTPEmailCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeOTPEmail,
			},
		},
		{
			name: "recovery code",
			fields: fields{
				RecoveryCodeCheckedAt: testNow,
			},
			want: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypeRecoveryCode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &SessionWriteModel{
				PasswordCheckedAt:     tt.fields.PasswordCheckedAt,
				IntentCheckedAt:       tt.fields.IntentCheckedAt,
				WebAuthNCheckedAt:     tt.fields.WebAuthNCheckedAt,
				WebAuthNUserVerified:  tt.fields.WebAuthNUserVerified,
				TOTPCheckedAt:         tt.fields.TOTPCheckedAt,
				OTPSMSCheckedAt:       tt.fields.OTPSMSCheckedAt,
				OTPEmailCheckedAt:     tt.fields.OTPEmailCheckedAt,
				RecoveryCodeCheckedAt: tt.fields.RecoveryCodeCheckedAt,
			}
			got := wm.AuthMethodTypes()
			assert.Equal(t, got, tt.want)
		})
	}
}
