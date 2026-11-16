package login

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
	provideridp "github.com/zitadel/zitadel/internal/repository/idp"
)

// Test_mapExternalNotFoundOptionFormDataToLoginUser is a regression test for
// GHSA-738m-7888-jfv8. The "external account not found" page lets a user edit
// their *profile* (name, email, language, ...) before the account is created,
// but the identity that account is bound to must come from the real, verified
// IDP callback (carried in authReq.LinkingUsers), never from the raw POST body.
//
// Before the fix, this mapping trusted the request's own external-* fields, so an
// attacker could forge an arbitrary (IDPConfigID, ExternalUserID) binding and a
// self-asserted "verified" email. After the fix, the security-critical fields are
// sourced from the trusted linkingUser, and the verified flags only survive when
// the (editable) form value still matches what the IDP actually returned.
func Test_mapExternalNotFoundOptionFormDataToLoginUser(t *testing.T) {
	type args struct {
		formData    *externalNotFoundOptionFormData
		linkingUser *domain.ExternalUser
	}
	tests := []struct {
		name string
		args args
		want *domain.ExternalUser
	}{
		{
			// The attacker sets the external identity fields in the POST body to a
			// victim's identifier / a chosen IDP. The result must ignore them and use
			// the trusted linkingUser instead.
			"forged external identity in form is ignored, trusted linkingUser wins",
			args{
				formData: &externalNotFoundOptionFormData{
					externalRegisterFormData: externalRegisterFormData{
						ExternalIDPConfigID:  "attacker-chosen-idp",
						ExternalIDPExtUserID: "victim-github-id",
						Email:                domain.EmailAddress("attacker@evil.test"),
						Username:             "attacker",
						Firstname:            "At",
						Lastname:             "Tacker",
						Nickname:             "a",
						Language:             "en",
					},
				},
				linkingUser: &domain.ExternalUser{
					IDPConfigID:    "real-idp",
					ExternalUserID: "real-callback-user-id",
					Email:          domain.EmailAddress("real@idp.test"),
				},
			},
			&domain.ExternalUser{
				IDPConfigID:       "real-idp",              // from trusted linkingUser, NOT the form
				ExternalUserID:    "real-callback-user-id", // from trusted linkingUser, NOT the form
				PreferredUsername: "attacker",
				DisplayName:       "attacker@evil.test",
				FirstName:         "At",
				LastName:          "Tacker",
				NickName:          "a",
				Email:             domain.EmailAddress("attacker@evil.test"),
				IsEmailVerified:   false, // linkingUser.IsEmailVerified is false
				PreferredLanguage: language.English,
			},
		},
		{
			// The attacker sets external-email-verified=true and email==external-email
			// in the POST body; before the fix that alone marked the email verified.
			"forged verified flag in form cannot mark email verified",
			args{
				formData: &externalNotFoundOptionFormData{
					externalRegisterFormData: externalRegisterFormData{
						ExternalEmail:         domain.EmailAddress("attacker@evil.test"),
						ExternalEmailVerified: true,
						Email:                 domain.EmailAddress("attacker@evil.test"),
						Language:              "en",
					},
				},
				linkingUser: &domain.ExternalUser{
					IDPConfigID:     "real-idp",
					ExternalUserID:  "real-callback-user-id",
					Email:           domain.EmailAddress("real@idp.test"),
					IsEmailVerified: false,
				},
			},
			&domain.ExternalUser{
				IDPConfigID:       "real-idp",
				ExternalUserID:    "real-callback-user-id",
				DisplayName:       "attacker@evil.test",
				Email:             domain.EmailAddress("attacker@evil.test"),
				IsEmailVerified:   false, // trusted linkingUser was not verified
				PreferredLanguage: language.English,
			},
		},
		{
			// Genuine flow: the IDP returned a verified email and the user did not
			// change it on the page, so it stays verified.
			"genuinely verified email preserved when user keeps it unchanged",
			args{
				formData: &externalNotFoundOptionFormData{
					externalRegisterFormData: externalRegisterFormData{
						Email:    domain.EmailAddress("real@idp.test"),
						Language: "en",
					},
				},
				linkingUser: &domain.ExternalUser{
					IDPConfigID:     "real-idp",
					ExternalUserID:  "real-callback-user-id",
					Email:           domain.EmailAddress("real@idp.test"),
					IsEmailVerified: true,
				},
			},
			&domain.ExternalUser{
				IDPConfigID:       "real-idp",
				ExternalUserID:    "real-callback-user-id",
				DisplayName:       "real@idp.test",
				Email:             domain.EmailAddress("real@idp.test"),
				IsEmailVerified:   true, // verified by IDP and unchanged by the user
				PreferredLanguage: language.English,
			},
		},
		{
			// Genuine flow but the user edits the email away from the IDP-verified one:
			// the new address is unverified until proven, so the flag must drop.
			"verified email drops when user edits it away from the IDP value",
			args{
				formData: &externalNotFoundOptionFormData{
					externalRegisterFormData: externalRegisterFormData{
						Email:    domain.EmailAddress("changed@user.test"),
						Language: "en",
					},
				},
				linkingUser: &domain.ExternalUser{
					IDPConfigID:     "real-idp",
					ExternalUserID:  "real-callback-user-id",
					Email:           domain.EmailAddress("real@idp.test"),
					IsEmailVerified: true,
				},
			},
			&domain.ExternalUser{
				IDPConfigID:       "real-idp",
				ExternalUserID:    "real-callback-user-id",
				DisplayName:       "changed@user.test",
				Email:             domain.EmailAddress("changed@user.test"),
				IsEmailVerified:   false, // edited away from the verified value
				PreferredLanguage: language.English,
			},
		},
		{
			// Same trust rule for the phone: verified only when the IDP verified it
			// and the user kept the value; a forged form phone-verified is ignored.
			"phone verified only when IDP-verified and unchanged",
			args{
				formData: &externalNotFoundOptionFormData{
					externalRegisterFormData: externalRegisterFormData{
						ExternalPhone:         domain.PhoneNumber("+41791234567"),
						ExternalPhoneVerified: true,
						Phone:                 domain.PhoneNumber("+41791234567"),
						Language:              "en",
					},
				},
				linkingUser: &domain.ExternalUser{
					IDPConfigID:     "real-idp",
					ExternalUserID:  "real-callback-user-id",
					Phone:           domain.PhoneNumber("+41791234567"),
					IsPhoneVerified: true,
				},
			},
			&domain.ExternalUser{
				IDPConfigID:       "real-idp",
				ExternalUserID:    "real-callback-user-id",
				Phone:             domain.PhoneNumber("+41791234567"),
				IsPhoneVerified:   true,
				PreferredLanguage: language.English,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapExternalNotFoundOptionFormDataToLoginUser(tt.args.formData, tt.args.linkingUser)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
			"non-string org domain is skipped",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": map[string]any{
						"orgID1": 55,
					},
				}},
			}),
			nil,
		},
		{
			"role keeps only orgs with string domains",
			openid.NewUser(&oidc.UserInfo{Claims: map[string]any{
				zitadelProjectRolesClaim: map[string]any{
					"IAM_OWNER_VIEWER": map[string]any{
						"orgID1": "org1.example.com",
						"orgID2": 55,
					},
				}},
			}),
			map[string]map[string]string{
				"IAM_OWNER_VIEWER": {"orgID1": "org1.example.com"},
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

func Test_preserveProjectRoles(t *testing.T) {
	roles := map[string]map[string]string{
		"IAM_OWNER_VIEWER": {"orgID1": "org1.example.com"},
	}
	tests := []struct {
		name         string
		user         *domain.ExternalUser
		linkingUsers []*domain.ExternalUser
		want         map[string]map[string]string
	}{
		{
			// the form-mapped user has no roles claim; it must be restored from the linking user
			name:         "roles copied from the single linking user",
			user:         &domain.ExternalUser{},
			linkingUsers: []*domain.ExternalUser{{ProjectRoles: roles}},
			want:         roles,
		},
		{
			// linking users accumulate on account-linking flows; the most recent one is mapped
			name: "roles copied from the most recent linking user",
			user: &domain.ExternalUser{},
			linkingUsers: []*domain.ExternalUser{
				{ProjectRoles: map[string]map[string]string{"OTHER_ROLE": {"orgIDx": "orgx.example.com"}}},
				{ProjectRoles: roles},
			},
			want: roles,
		},
		{
			name:         "no linking users, nil",
			user:         &domain.ExternalUser{},
			linkingUsers: nil,
			want:         nil,
		},
		{
			name:         "last linking user without roles, nil",
			user:         &domain.ExternalUser{},
			linkingUsers: []*domain.ExternalUser{{ProjectRoles: roles}, {}},
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preserveProjectRoles(tt.user, tt.linkingUsers)
			assert.Equal(t, tt.want, tt.user.ProjectRoles)
		})
	}
}

func Test_supportUserInstanceRoles(t *testing.T) {
	configured := &query.ZitadelIDPTemplate{
		InstanceRolesInfo: []provideridp.RolesInfo{{OrganizationID: "orgID1", OrganizationDomain: "org1.example.com"}},
	}
	supportClaim := map[string]map[string]string{
		domain.RoleIAMOwnerViewer: {"orgID1": "org1.example.com"},
	}
	tests := []struct {
		name         string
		provider     *query.IDPTemplate
		externalUser *domain.ExternalUser
		want         []string
	}{
		{
			name:         "instance zitadel provider, support role claim, matching org",
			provider:     &query.IDPTemplate{OwnerType: domain.IdentityProviderTypeSystem, Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: supportClaim},
			want:         []string{domain.RoleIAMOwnerViewer},
		},
		{
			name:     "multiple IAM roles claimed, all returned sorted",
			provider: &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{
				domain.RoleIAMOwnerViewer: {"orgID1": "org1.example.com"},
				domain.RoleIAMOwner:       {"orgID1": "org1.example.com"},
			}},
			want: []string{domain.RoleIAMOwner, domain.RoleIAMOwnerViewer},
		},
		{
			name:     "only the roles of a configured org are returned",
			provider: &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{
				domain.RoleIAMOwnerViewer: {"orgID1": "org1.example.com"},
				domain.RoleIAMOwner:       {"orgID2": "org2.example.com"},
			}},
			want: []string{domain.RoleIAMOwnerViewer},
		},
		{
			name:     "non-IAM roles are not returned",
			provider: &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{
				domain.RoleIAMOwnerViewer: {"orgID1": "org1.example.com"},
				domain.RoleOrgOwner:       {"orgID1": "org1.example.com"},
			}},
			want: []string{domain.RoleIAMOwnerViewer},
		},
		{
			name:         "org-scoped zitadel provider, matching org, skipped",
			provider:     &query.IDPTemplate{OwnerType: domain.IdentityProviderTypeOrg, Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: supportClaim},
			want:         nil,
		},
		{
			name:         "non-zitadel provider, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeOIDC, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: supportClaim},
			want:         nil,
		},
		{
			name:         "zitadel type but no zitadel template, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeZitadel},
			externalUser: &domain.ExternalUser{ProjectRoles: supportClaim},
			want:         nil,
		},
		{
			name:         "no roles claim at all, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{},
			want:         nil,
		},
		{
			name:         "roles claim without an IAM role, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{"OTHER_ROLE": {"orgID1": "org1.example.com"}}},
			want:         nil,
		},
		{
			name:         "support role present but empty org set, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{domain.RoleIAMOwnerViewer: {}}},
			want:         nil,
		},
		{
			name:         "support role claim for a non-configured org, skipped",
			provider:     &query.IDPTemplate{Type: domain.IDPTypeZitadel, ZitadelIDPTemplate: configured},
			externalUser: &domain.ExternalUser{ProjectRoles: map[string]map[string]string{domain.RoleIAMOwnerViewer: {"orgID1": "evil.example.com"}}},
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, supportUserInstanceRoles(tt.provider, tt.externalUser))
		})
	}
}
