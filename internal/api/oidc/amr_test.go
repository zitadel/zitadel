package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
)

func TestAMR(t *testing.T) {
	type args struct {
		methodTypes []domain.UserAuthMethodType
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no checks, empty",
			args{
				nil,
			},
			nil,
		},
		{
			"pw checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
			},
			[]string{PWD},
		},
		{
			"passkey checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePasswordless},
			},
			[]string{UserPresence, MFA},
		},
		{
			"u2f checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypeU2F},
			},
			[]string{UserPresence},
		},
		{
			"totp checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypeTOTP},
			},
			[]string{OTP},
		},
		{
			"otp sms checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypeOTPSMS},
			},
			[]string{OTP},
		},
		{
			"otp email checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypeOTPEmail},
			},
			[]string{OTP},
		},
		{
			"multiple (t)otp checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypeTOTP, domain.UserAuthMethodTypeOTPEmail},
			},
			[]string{OTP, MFA},
		},
		{
			"multiple checked",
			args{
				[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword, domain.UserAuthMethodTypeU2F},
			},
			[]string{PWD, UserPresence, MFA},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AuthMethodTypesToAMR(tt.args.methodTypes)
			assert.Equal(t, tt.want, got)
		})
	}
}
