package domain

import (
	"testing"
)

func TestEmailValid(t *testing.T) {
	type args struct {
		email *Email
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "empty email, invalid",
			args: args{
				email: &Email{},
			},
			result: false,
		},
		{
			name: "only letters email, invalid",
			args: args{
				email: &Email{EmailAddress: "testemail"},
			},
			result: false,
		},
		{
			name: "nothing after @, invalid",
			args: args{
				email: &Email{EmailAddress: "testemail@"},
			},
			result: false,
		},
		{
			name: "email, valid",
			args: args{
				email: &Email{EmailAddress: "testemail@gmail.com"},
			},
			result: true,
		},
		{
			name: "email, valid",
			args: args{
				email: &Email{EmailAddress: "test.email@gmail.com"},
			},
			result: true,
		},
		{
			name: "email, valid",
			args: args{
				email: &Email{EmailAddress: "test/email@gmail.com"},
			},
			result: true,
		},
		{
			name: "email, valid",
			args: args{
				email: &Email{EmailAddress: "test/email@gmail.com"},
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.email.IsValid()
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
