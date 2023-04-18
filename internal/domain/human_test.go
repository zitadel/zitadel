package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHuman_EnsureDisplayName(t *testing.T) {
	type fields struct {
		Username string
		Profile  *Profile
		Email    *Email
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"display name set",
			fields{
				Username: "username",
				Profile: &Profile{
					FirstName:   "firstName",
					LastName:    "lastName",
					DisplayName: "displayName",
				},
				Email: &Email{
					EmailAddress: "email",
				},
			},
			"displayName",
		},
		{
			"first and lastname set",
			fields{
				Username: "username",
				Profile: &Profile{
					FirstName: "firstName",
					LastName:  "lastName",
				},
				Email: &Email{
					EmailAddress: "email",
				},
			},
			"firstName lastName",
		},
		{
			"profile nil, email set",
			fields{
				Username: "username",
				Email: &Email{
					EmailAddress: "email",
				},
			},
			"email",
		},
		{
			"email nil, username set",
			fields{
				Username: "username",
			},
			"username",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Human{
				Username: tt.fields.Username,
				Profile:  tt.fields.Profile,
				Email:    tt.fields.Email,
			}
			u.EnsureDisplayName()
			assert.Equal(t, tt.want, u.DisplayName)
		})
	}
}
