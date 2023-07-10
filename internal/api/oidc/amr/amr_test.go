package amr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAMR(t *testing.T) {
	type args struct {
		model AuthenticationMethodReference
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"no checks, empty",
			args{
				new(test),
			},
			[]string{},
		},
		{
			"pw checked",
			args{
				&test{pwChecked: true},
			},
			[]string{PWD},
		},
		{
			"passkey checked",
			args{
				&test{passkeyChecked: true},
			},
			[]string{UserPresence, MFA},
		},
		{
			"u2f checked",
			args{
				&test{u2fChecked: true},
			},
			[]string{UserPresence},
		},
		{
			"otp checked",
			args{
				&test{otpChecked: true},
			},
			[]string{OTP},
		},
		{
			"multiple checked",
			args{
				&test{
					pwChecked:  true,
					u2fChecked: true,
				},
			},
			[]string{PWD, UserPresence, MFA},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := List(tt.args.model)
			assert.Equal(t, tt.want, got)
		})
	}
}

type test struct {
	pwChecked      bool
	passkeyChecked bool
	u2fChecked     bool
	otpChecked     bool
}

func (t test) IsPasswordChecked() bool {
	return t.pwChecked
}

func (t test) IsPasskeyChecked() bool {
	return t.passkeyChecked
}

func (t test) IsU2FChecked() bool {
	return t.u2fChecked
}

func (t test) IsOTPChecked() bool {
	return t.otpChecked
}
