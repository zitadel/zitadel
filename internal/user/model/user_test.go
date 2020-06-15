package model

import (
	"testing"
)

func TestIsUserValid(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name    string
		args    args
		result  bool
		errFunc func(err error) bool
	}{
		{
			name: "user with minimal data",
			args: args{
				user: &User{
					Profile: &Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &Email{
						EmailAddress: "Email",
					},
				},
			},
			result: true,
		},
		{
			name: "user with phone data",
			args: args{
				user: &User{
					Profile: &Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &Email{
						EmailAddress: "Email",
					},
					Phone: &Phone{
						PhoneNumber: "+41711234569",
					},
				},
			},
			result: true,
		},
		{
			name: "user with address data",
			args: args{
				user: &User{
					Profile: &Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &Email{
						EmailAddress: "Email",
					},
					Address: &Address{
						StreetAddress: "Teufenerstrasse 19",
						PostalCode:    "9000",
						Locality:      "St. Gallen",
						Country:       "Switzerland",
					},
				},
			},
			result: true,
		},
		{
			name: "user with all data",
			args: args{
				user: &User{
					Profile: &Profile{
						UserName:  "UserName",
						FirstName: "FirstName",
						LastName:  "LastName",
					},
					Email: &Email{
						EmailAddress: "Email",
					},
					Phone: &Phone{
						PhoneNumber: "+41711234569",
					},
					Address: &Address{
						StreetAddress: "Teufenerstrasse 19",
						PostalCode:    "9000",
						Locality:      "St. Gallen",
						Country:       "Switzerland",
					},
				},
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.args.user.IsValid()

			if tt.result != isValid {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, isValid)
			}
		})
	}
}
