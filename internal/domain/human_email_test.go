package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
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
			wantErr: caos_errs.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "execution error",
			args: args{
				tmpl:   "{{.Foo}}",
				userID: "user1",
				code:   "123",
				orgID:  "org1",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate"),
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
