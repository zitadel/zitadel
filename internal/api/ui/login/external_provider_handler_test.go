package login

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
	provideridp "github.com/zitadel/zitadel/internal/repository/idp"
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

func Test_projectRolesFromIDPUser(t *testing.T) {
	tests := []struct {
		name string
		user idp.User
		want map[string]map[string]string
	}{
		{
			"not an oidc user",
			struct{ idp.User }{},
			nil,
		},
		{
			"no claims",
			openid.NewUser(&oidc.UserInfo{Claims: nil}),
			nil,
		},
		{
			"claim absent",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{"some": "other"}}),
			nil,
		},
		{
			"claim wrong type",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{zitadelProjectRolesClaim: "not-a-map"}}),
			nil,
		},
		{
			"role with non-map orgs is skipped",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": "not-a-map",
				}},
			}),
			nil,
		},
		{
			"single role and org",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": map[string]any{
						"orgID1": "org1.example.com",
					},
				}},
			}),
			map[string]map[string]string{
				"IAM_OWNER_VIEWER": {"orgID1": "org1.example.com"},
			},
		},
		{
			"multiple roles and orgs",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": map[string]any{
						"orgID1": "org1.example.com",
						"orgID2": "org2.example.com",
					},
					"OTHER_ROLE": map[string]any{
						"orgID3": "org3.example.com",
					},
				}},
			}),
			map[string]map[string]string{
				"IAM_OWNER_VIEWER": {"orgID1": "org1.example.com", "orgID2": "org2.example.com"},
				"OTHER_ROLE":       {"orgID3": "org3.example.com"},
			},
		},
		{
			"invalid org domain empty",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": map[string]any{
						"orgID1": 55,
					},
				}},
			}),
			map[string]map[string]string{
				"IAM_OWNER_VIEWER": {"orgID1": ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, projectRolesFromIDPUser(tt.user))
		})
	}
}

func Test_claimMatchesConfiguredOrg(t *testing.T) {
	tests := []struct {
		name      string
		claimOrgs map[string]string
		template  *query.ZitadelIDPTemplate
		want      bool
	}{
		{
			"id and domain match",
			map[string]string{"orgID1": "org1.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			true,
		},
		{
			"id matches, domain differs",
			map[string]string{"orgID1": "evil.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			false,
		},
		{
			"domain matches, id differs",
			map[string]string{"otherID": "org1.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			false,
		},
		{
			"neither id nor domain match",
			map[string]string{"orgIDx": "orgx.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			false,
		},
		{
			"no configured orgs",
			map[string]string{"orgID1": "org1.example.com"},
			&query.ZitadelIDPTemplate{},
			false,
		},
		{
			"no claim orgs",
			map[string]string{},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			false,
		},
		{
			"one of multiple claim orgs matches",
			map[string]string{"orgIDx": "orgx.example.com", "orgID1": "org1.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}}},
			true,
		},
		{
			"one of multiple configured orgs matches",
			map[string]string{"orgID2": "org2.example.com"},
			&query.ZitadelIDPTemplate{InstanceRolesInfo: []provideridp.RolesInfo{
				{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"},
				{OrganizationID: "orgID2", OrganizationDomain: "org2.example.com"},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, claimMatchesConfiguredOrg(tt.claimOrgs, tt.template))
		})
	}
}
