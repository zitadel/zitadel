package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMFAType_UserAuthMethodType(t *testing.T) {
	tests := []struct {
		name string
		m    MFAType
		want UserAuthMethodType
	}{
		{
			name: "totp",
			m:    MFATypeTOTP,
			want: UserAuthMethodTypeTOTP,
		},
		{
			name: "u2f",
			m:    MFATypeU2F,
			want: UserAuthMethodTypeU2F,
		},
		{
			name: "passwordless",
			m:    MFATypeU2FUserVerification,
			want: UserAuthMethodTypePasswordless,
		},
		{
			name: "otp sms",
			m:    MFATypeOTPSMS,
			want: UserAuthMethodTypeOTPSMS,
		},
		{
			name: "otp email",
			m:    MFATypeOTPEmail,
			want: UserAuthMethodTypeOTPEmail,
		},
		{
			name: "unspecified",
			m:    99,
			want: UserAuthMethodTypeUnspecified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.UserAuthMethodType()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthRequest_UserAuthMethodTypes(t *testing.T) {
	type fields struct {
		PasswordVerified bool
		MFAsVerified     []MFAType
	}
	tests := []struct {
		name   string
		fields fields
		want   []UserAuthMethodType
	}{
		{
			name: "no auth methods",
			fields: fields{
				PasswordVerified: false,
				MFAsVerified:     nil,
			},
			want: []UserAuthMethodType{},
		},
		{
			name: "only password",
			fields: fields{
				PasswordVerified: true,
				MFAsVerified:     nil,
			},
			want: []UserAuthMethodType{
				UserAuthMethodTypePassword,
			},
		},
		{
			name: "password, with mfa",
			fields: fields{
				PasswordVerified: true,
				MFAsVerified: []MFAType{
					MFATypeTOTP,
					MFATypeU2F,
				},
			},
			want: []UserAuthMethodType{
				UserAuthMethodTypePassword,
				UserAuthMethodTypeTOTP,
				UserAuthMethodTypeU2F,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthRequest{
				PasswordVerified: tt.fields.PasswordVerified,
				MFAsVerified:     tt.fields.MFAsVerified,
			}
			got := a.UserAuthMethodTypes()
			assert.Equal(t, tt.want, got)
		})
	}
}
