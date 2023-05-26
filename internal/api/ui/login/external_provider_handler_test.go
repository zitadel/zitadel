package login

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func Test_hasEmailChanged(t *testing.T) {
	type args struct {
		user         *query.User
		externalUser *domain.ExternalUser
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"no external mail",
			args{
				user:         &query.User{},
				externalUser: &domain.ExternalUser{},
			},
			false,
		},
		{
			"same email unverified",
			args{
				user: &query.User{
					Human: &query.Human{
						Email: domain.EmailAddress("email@test.com"),
					},
				},
				externalUser: &domain.ExternalUser{
					Email: domain.EmailAddress("email@test.com"),
				},
			},
			false,
		},
		{
			"same email verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Email:           domain.EmailAddress("email@test.com"),
						IsEmailVerified: true,
					},
				},
				externalUser: &domain.ExternalUser{
					Email:           domain.EmailAddress("email@test.com"),
					IsEmailVerified: true,
				},
			},
			false,
		},
		{
			"email already verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Email:           domain.EmailAddress("email@test.com"),
						IsEmailVerified: true,
					},
				},
				externalUser: &domain.ExternalUser{
					Email: domain.EmailAddress("email@test.com"),
				},
			},
			false,
		},
		{
			"email changed to verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Email: domain.EmailAddress("email@test.com"),
					},
				},
				externalUser: &domain.ExternalUser{
					Email:           domain.EmailAddress("email@test.com"),
					IsEmailVerified: true,
				},
			},
			true,
		},
		{
			"email changed",
			args{
				user: &query.User{
					Human: &query.Human{
						Email: domain.EmailAddress("email@test.com"),
					},
				},
				externalUser: &domain.ExternalUser{
					Email: domain.EmailAddress("new-email@test.com"),
				},
			},
			true,
		},
		{
			"email changed and verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Email:           domain.EmailAddress("email@test.com"),
						IsEmailVerified: false,
					},
				},
				externalUser: &domain.ExternalUser{
					Email:           domain.EmailAddress("new-email@test.com"),
					IsEmailVerified: true,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasEmailChanged(tt.args.user, tt.args.externalUser); got != tt.want {
				t.Errorf("hasEmailChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasPhoneChanged(t *testing.T) {
	type args struct {
		user         *query.User
		externalUser *domain.ExternalUser
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"no external phone",
			args{
				user:         &query.User{},
				externalUser: &domain.ExternalUser{},
			},
			false,
			false,
		},
		{
			"invalid phone",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone: domain.PhoneNumber("invalid"),
				},
			},
			false,
			true,
		},
		{
			"same phone unverified",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone: domain.PhoneNumber("+41791234567"),
				},
			},
			false,
			false,
		},
		{
			"same phone verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone:           domain.PhoneNumber("+41791234567"),
						IsPhoneVerified: true,
					},
				},
				externalUser: &domain.ExternalUser{
					Phone:           domain.PhoneNumber("+41791234567"),
					IsPhoneVerified: true,
				},
			},
			false,
			false,
		},
		{
			"phone already verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone:           domain.PhoneNumber("+41791234567"),
						IsPhoneVerified: true,
					},
				},
				externalUser: &domain.ExternalUser{
					Phone: domain.PhoneNumber("+41791234567"),
				},
			},
			false,
			false,
		},
		{
			"phone changed to verified",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone:           domain.PhoneNumber("+41791234567"),
					IsPhoneVerified: true,
				},
			},
			true,
			false,
		},
		{
			"phone changed",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone: domain.PhoneNumber("+4179654321"),
				},
			},
			true,
			false,
		},
		{
			"phone changed",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone:           domain.PhoneNumber("+4179654321"),
					IsPhoneVerified: true,
				},
			},
			true,
			false,
		},
		{
			"normalized phone unchanged",
			args{
				user: &query.User{
					Human: &query.Human{
						Phone: domain.PhoneNumber("+41791234567"),
					},
				},
				externalUser: &domain.ExternalUser{
					Phone: domain.PhoneNumber("0791234567"),
				},
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hasPhoneChanged(tt.args.user, tt.args.externalUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("hasPhoneChanged() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hasPhoneChanged() got = %v, want %v", got, tt.want)
			}
		})
	}
}
