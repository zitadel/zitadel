package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
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
		{
			name: "email UTF-8, valid", // https://github.com/zitadel/zitadel/issues/9821
			args: args{
				email: &Email{EmailAddress: "abc@sünde.com"},
			},
			result: true,
		},
		{
			name: "email with name, invalid",
			args: args{
				email: &Email{EmailAddress: "John Doe <john.doe@example.com>"},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.email.Validate() == nil
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func TestRenderConfirmURLTemplate(t *testing.T) {
	type args struct {
		tmpl   string
		userID string
		code   string
		orgID  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "invalid template",
			args: args{
				tmpl:   "{{",
				userID: "user1",
				code:   "123",
				orgID:  "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "execution error",
			args: args{
				tmpl:   "{{.Foo}}",
				userID: "user1",
				code:   "123",
				orgID:  "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "success",
			args: args{
				tmpl:   "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				userID: "user1",
				code:   "123",
				orgID:  "org1",
			},
			want: "https://example.com/email/verify?userID=user1&code=123&orgID=org1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w strings.Builder
			err := RenderConfirmURLTemplate(&w, tt.args.tmpl, tt.args.userID, tt.args.code, tt.args.orgID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, w.String())
		})
	}
}

func TestEmailAddress_Domain(t *testing.T) {
	tests := []struct {
		address string
		want    string
	}{
		{address: "user@example.com", want: "example.com"},
		{address: " user@example.com ", want: "example.com"},
		{address: "user@EXAMPLE.com", want: "EXAMPLE.com"},
		{address: "us@er@example.com", want: "example.com"},
		{address: "user", want: ""},
		{address: "user@", want: ""},
		{address: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			assert.Equal(t, tt.want, EmailAddress(tt.address).Domain())
		})
	}
}

func TestEmailAddress_IsReservedDomain(t *testing.T) {
	tests := []struct {
		address string
		want    bool
	}{
		{address: "user@example.com", want: true},
		{address: "user@example.net", want: true},
		{address: "user@example.org", want: true},
		{address: "user@EXAMPLE.COM", want: true},
		{address: "user@mail.example.com", want: true},
		{address: "user@example.com.", want: true},
		{address: "user@localhost", want: true},
		{address: "user@mail.localhost", want: true},
		{address: "user@test", want: true},
		{address: "user@invalid", want: true},
		{address: "user@user.test", want: true},
		{address: "user@user.example", want: true},
		{address: "user@user.invalid", want: true},
		{address: "user@user.test.", want: true},
		{address: "user@example.ch", want: false},
		{address: "user@examplecom", want: false},
		{address: "user@example.com.ch", want: false},
		{address: "user@test.other.example", want: true},
		{address: "user@testing.ch", want: false},
		{address: "user@other.example.ch", want: false},
		{address: "user", want: false},
		{address: "", want: false},
		{address: " user@example.com ", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			assert.Equal(t, tt.want, EmailAddress(tt.address).IsReservedDomain())
		})
	}
}
