package oidc

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func Test_prepareRoles(t *testing.T) {
	type args struct {
		projectID            string
		scope                []string
		projectRoleAssertion bool
		currentProjectOnly   bool
	}
	tests := []struct {
		name               string
		args               args
		wantRoleAudience   []string
		wantRequestedRoles []string
	}{
		{
			name: "empty scope",
			args: args{
				projectID:            "projID",
				scope:                nil,
				projectRoleAssertion: false,
				currentProjectOnly:   false,
			},
			wantRoleAudience:   nil,
			wantRequestedRoles: nil,
		},
		{
			name: "project role assertion",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: true,
				scope:                nil,
				currentProjectOnly:   false,
			},
			wantRoleAudience:   []string{"projID"},
			wantRequestedRoles: nil,
		},
		{
			name: "some scope, current project only",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope:                []string{"openid", "profile"},
				currentProjectOnly:   true,
			},
			wantRoleAudience:   []string{"projID"},
			wantRequestedRoles: nil,
		},
		{
			name: "scope projects roles",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectsRoles,
					domain.ProjectIDScope + "project2" + domain.AudSuffix,
				},
				currentProjectOnly: false,
			},
			wantRoleAudience:   []string{"project2", "projID"},
			wantRequestedRoles: nil,
		},
		{
			name: "scope projects roles ignored, current project only",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectsRoles,
					domain.ProjectIDScope + "project2" + domain.AudSuffix,
				},
				currentProjectOnly: true,
			},
			wantRoleAudience:   []string{"projID"},
			wantRequestedRoles: nil,
		},
		{
			name: "scope project role prefix",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectRolePrefix + "foo",
					ScopeProjectRolePrefix + "bar",
				},
				currentProjectOnly: false,
			},
			wantRoleAudience:   []string{"projID"},
			wantRequestedRoles: []string{"foo", "bar"},
		},
		{
			name: "scope project role prefix and audience",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectRolePrefix + "foo",
					ScopeProjectRolePrefix + "bar",
					domain.ProjectIDScope + "project2" + domain.AudSuffix,
				},
				currentProjectOnly: false,
			},
			wantRoleAudience:   []string{"projID", "project2"},
			wantRequestedRoles: []string{"foo", "bar"},
		},
		{
			name: "scope project role prefix and audience ignored, current project only",
			args: args{
				projectID:            "projID",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectRolePrefix + "foo",
					ScopeProjectRolePrefix + "bar",
					domain.ProjectIDScope + "project2" + domain.AudSuffix,
				},
				currentProjectOnly: true,
			},
			wantRoleAudience:   []string{"projID"},
			wantRequestedRoles: []string{"foo", "bar"},
		},
		{
			name: "no projectID, scope project role prefix and audience",
			args: args{
				projectID:            "",
				projectRoleAssertion: false,
				scope: []string{
					"openid", "profile",
					ScopeProjectRolePrefix + "foo",
					ScopeProjectRolePrefix + "bar",
					domain.ProjectIDScope + "project2" + domain.AudSuffix,
				},
				currentProjectOnly: false,
			},
			wantRoleAudience:   []string{"project2"},
			wantRequestedRoles: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoleAudience, gotRequestedRoles := prepareRoles(context.Background(), tt.args.scope, tt.args.projectID, tt.args.projectRoleAssertion, tt.args.currentProjectOnly)
			assert.ElementsMatch(t, tt.wantRoleAudience, gotRoleAudience, "roleAudience")
			assert.ElementsMatch(t, tt.wantRequestedRoles, gotRequestedRoles, "requestedRoles")
		})
	}
}

func Test_userInfoToOIDC(t *testing.T) {
	metadata := []query.UserMetadata{
		{
			Key:   "key1",
			Value: []byte{1, 2, 3},
		},
		{
			Key:   "key2",
			Value: []byte{4, 5, 6},
		},
	}
	organization := &query.UserInfoOrg{
		ID:            "orgID",
		Name:          "orgName",
		PrimaryDomain: "orgDomain",
	}
	humanUserInfo := &query.OIDCUserInfo{
		User: &query.User{
			ID:                 "human1",
			CreationDate:       time.Unix(123, 456),
			ChangeDate:         time.Unix(567, 890),
			ResourceOwner:      "orgID",
			Sequence:           22,
			State:              domain.UserStateActive,
			Type:               domain.UserTypeHuman,
			Username:           "username",
			LoginNames:         []string{"foo", "bar"},
			PreferredLoginName: "foo",
			Human: &query.Human{
				FirstName:         "user",
				LastName:          "name",
				NickName:          "foobar",
				DisplayName:       "xxx",
				AvatarKey:         "picture.png",
				PreferredLanguage: language.Dutch,
				Gender:            domain.GenderDiverse,
				Email:             "foo@bar.com",
				IsEmailVerified:   true,
				Phone:             "+31123456789",
				IsPhoneVerified:   true,
			},
		},
		Metadata: metadata,
		Org:      organization,
		UserGrants: []query.UserGrant{
			{
				ID:                "ug1",
				CreationDate:      time.Unix(444, 444),
				ChangeDate:        time.Unix(555, 555),
				Sequence:          55,
				Roles:             []string{"role1", "role2"},
				GrantID:           "grantID",
				State:             domain.UserGrantStateActive,
				UserID:            "human1",
				Username:          "username",
				ResourceOwner:     "orgID",
				ProjectID:         "project1",
				OrgName:           "orgName",
				OrgPrimaryDomain:  "orgDomain",
				ProjectName:       "projectName",
				UserResourceOwner: "org1",
			},
		},
	}
	machineUserInfo := &query.OIDCUserInfo{
		User: &query.User{
			ID:                 "machine1",
			CreationDate:       time.Unix(123, 456),
			ChangeDate:         time.Unix(567, 890),
			ResourceOwner:      "orgID",
			Sequence:           23,
			State:              domain.UserStateActive,
			Type:               domain.UserTypeMachine,
			Username:           "machine",
			PreferredLoginName: "meanMachine",
			Machine: &query.Machine{
				Name:        "machine",
				Description: "I'm a robot",
			},
		},
		Org: organization,
		UserGrants: []query.UserGrant{
			{
				ID:                "ug1",
				CreationDate:      time.Unix(444, 444),
				ChangeDate:        time.Unix(555, 555),
				Sequence:          55,
				Roles:             []string{"role1", "role2"},
				GrantID:           "grantID",
				State:             domain.UserGrantStateActive,
				UserID:            "human1",
				Username:          "username",
				ResourceOwner:     "orgID",
				ProjectID:         "project1",
				OrgName:           "orgName",
				OrgPrimaryDomain:  "orgDomain",
				ProjectName:       "projectName",
				UserResourceOwner: "org1",
			},
		},
	}

	type args struct {
		user              *query.OIDCUserInfo
		userInfoAssertion bool
		scope             []string
	}
	tests := []struct {
		name string
		args args
		want *oidc.UserInfo
	}{
		{
			name: "human, empty",
			args: args{
				user: humanUserInfo,
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, empty",
			args: args{
				user: machineUserInfo,
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "human, scope openid",
			args: args{
				user:  humanUserInfo,
				scope: []string{oidc.ScopeOpenID},
			},
			want: &oidc.UserInfo{
				Subject: "human1",
			},
		},
		{
			name: "machine, scope openid",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopeOpenID},
			},
			want: &oidc.UserInfo{
				Subject: "machine1",
			},
		},
		{
			name: "human, scope email, profileInfoAssertion",
			args: args{
				user:              humanUserInfo,
				userInfoAssertion: true,
				scope:             []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{
				UserInfoEmail: oidc.UserInfoEmail{
					Email:         "foo@bar.com",
					EmailVerified: true,
				},
			},
		},
		{
			name: "human, scope email",
			args: args{
				user:  humanUserInfo,
				scope: []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, scope email, profileInfoAssertion",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{
				UserInfoEmail: oidc.UserInfoEmail{},
			},
		},
		{
			name: "human, scope profile, profileInfoAssertion",
			args: args{
				user:              humanUserInfo,
				userInfoAssertion: true,
				scope:             []string{oidc.ScopeProfile},
			},
			want: &oidc.UserInfo{
				UserInfoProfile: oidc.UserInfoProfile{
					Name:              "xxx",
					GivenName:         "user",
					FamilyName:        "name",
					Nickname:          "foobar",
					Picture:           "https://foo.com/assets/orgID/picture.png",
					Gender:            "diverse",
					Locale:            oidc.NewLocale(language.Dutch),
					UpdatedAt:         oidc.FromTime(time.Unix(567, 890)),
					PreferredUsername: "foo",
				},
			},
		},
		{
			name: "machine, scope profile, profileInfoAssertion",
			args: args{
				user:              machineUserInfo,
				userInfoAssertion: true,
				scope:             []string{oidc.ScopeProfile},
			},
			want: &oidc.UserInfo{
				UserInfoProfile: oidc.UserInfoProfile{
					Name:              "machine",
					UpdatedAt:         oidc.FromTime(time.Unix(567, 890)),
					PreferredUsername: "meanMachine",
				},
			},
		},
		{
			name: "machine, scope profile",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopeProfile},
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "human, scope phone, profileInfoAssertion",
			args: args{
				user:              humanUserInfo,
				userInfoAssertion: true,
				scope:             []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				UserInfoPhone: oidc.UserInfoPhone{
					PhoneNumber:         "+31123456789",
					PhoneNumberVerified: true,
				},
			},
		},
		{
			name: "human, scope phone",
			args: args{
				user:  humanUserInfo,
				scope: []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, scope phone",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				UserInfoPhone: oidc.UserInfoPhone{},
			},
		},
		{
			name: "human, scope metadata",
			args: args{
				user:  humanUserInfo,
				scope: []string{ScopeUserMetaData},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					ClaimUserMetaData: map[string]string{
						"key1": base64.RawURLEncoding.EncodeToString([]byte{1, 2, 3}),
						"key2": base64.RawURLEncoding.EncodeToString([]byte{4, 5, 6}),
					},
				},
			},
		},
		{
			name: "machine, scope metadata, none found",
			args: args{
				user:  machineUserInfo,
				scope: []string{ScopeUserMetaData},
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, scope resource owner",
			args: args{
				user:  machineUserInfo,
				scope: []string{ScopeResourceOwner},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					ClaimResourceOwnerID:            "orgID",
					ClaimResourceOwnerName:          "orgName",
					ClaimResourceOwnerPrimaryDomain: "orgDomain",
				},
			},
		},
		{
			name: "human, scope org primary domain prefix",
			args: args{
				user:  humanUserInfo,
				scope: []string{domain.OrgDomainPrimaryScope + "foo.com"},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					domain.OrgDomainPrimaryClaim: "foo.com",
				},
			},
		},
		{
			name: "machine, scope org id",
			args: args{
				user:  machineUserInfo,
				scope: []string{domain.OrgIDScope + "orgID"},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					domain.OrgIDClaim:               "orgID",
					ClaimResourceOwnerID:            "orgID",
					ClaimResourceOwnerName:          "orgName",
					ClaimResourceOwnerPrimaryDomain: "orgDomain",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assetPrefix := "https://foo.com/assets"
			got := userInfoToOIDC(tt.args.user, tt.args.userInfoAssertion, tt.args.scope, assetPrefix)
			assert.Equal(t, tt.want, got)
		})
	}
}
