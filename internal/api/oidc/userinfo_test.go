package oidc

import (
	"context"
	"encoding/base64"
	"fmt"
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
		projectID    string
		scope        []string
		roleAudience []string
	}
	tests := []struct {
		name               string
		args               args
		wantRa             []string
		wantRequestedRoles []string
	}{
		{
			name: "empty scope and roleAudience",
			args: args{
				projectID:    "projID",
				scope:        nil,
				roleAudience: nil,
			},
			wantRa:             nil,
			wantRequestedRoles: nil,
		},
		{
			name: "some scope and roleAudience",
			args: args{
				projectID:    "projID",
				scope:        []string{"openid", "profile"},
				roleAudience: []string{"project2"},
			},
			wantRa:             []string{"project2", "projID"},
			wantRequestedRoles: []string{},
		},
		{
			name: "scope projects roles",
			args: args{
				projectID:    "projID",
				scope:        []string{ScopeProjectsRoles, domain.ProjectIDScope + "project2" + domain.AudSuffix},
				roleAudience: nil,
			},
			wantRa:             []string{"project2", "projID"},
			wantRequestedRoles: []string{},
		},
		{
			name: "scope project role prefix",
			args: args{
				projectID:    "projID",
				scope:        []string{"openid", "profile", ScopeProjectRolePrefix + "foo", ScopeProjectRolePrefix + "bar"},
				roleAudience: nil,
			},
			wantRa:             []string{"projID"},
			wantRequestedRoles: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRa, gotRequestedRoles := prepareRoles(context.Background(), tt.args.projectID, tt.args.scope, tt.args.roleAudience)
			assert.Equal(t, tt.wantRa, gotRa, "roleAudience")
			assert.Equal(t, tt.wantRequestedRoles, gotRequestedRoles, "requestedRoles")
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
		projectID      string
		user           *query.OIDCUserInfo
		scope          []string
		roleAudience   []string
		requestedRoles []string
	}
	tests := []struct {
		name string
		args args
		want *oidc.UserInfo
	}{
		{
			name: "human, empty",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, empty",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "human, scope openid",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{oidc.ScopeOpenID},
			},
			want: &oidc.UserInfo{
				Subject: "human1",
			},
		},
		{
			name: "machine, scope openid",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{oidc.ScopeOpenID},
			},
			want: &oidc.UserInfo{
				Subject: "machine1",
			},
		},
		{
			name: "human, scope email",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{
				UserInfoEmail: oidc.UserInfoEmail{
					Email:         "foo@bar.com",
					EmailVerified: true,
				},
			},
		},
		{
			name: "machine, scope email",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{
				UserInfoEmail: oidc.UserInfoEmail{},
			},
		},
		{
			name: "human, scope profile",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{oidc.ScopeProfile},
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
			name: "machine, scope profile",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{oidc.ScopeProfile},
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
			name: "human, scope phone",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				UserInfoPhone: oidc.UserInfoPhone{
					PhoneNumber:         "+31123456789",
					PhoneNumberVerified: true,
				},
			},
		},
		{
			name: "machine, scope phone",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				UserInfoPhone: oidc.UserInfoPhone{},
			},
		},
		{
			name: "human, scope metadata",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{ScopeUserMetaData},
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
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{ScopeUserMetaData},
			},
			want: &oidc.UserInfo{},
		},
		{
			name: "machine, scope resource owner",
			args: args{
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{ScopeResourceOwner},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					ClaimResourceOwner + "id":             "orgID",
					ClaimResourceOwner + "name":           "orgName",
					ClaimResourceOwner + "primary_domain": "orgDomain",
				},
			},
		},
		{
			name: "human, scope org primary domain prefix",
			args: args{
				projectID: "project1",
				user:      humanUserInfo,
				scope:     []string{domain.OrgDomainPrimaryScope + "foo.com"},
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
				projectID: "project1",
				user:      machineUserInfo,
				scope:     []string{domain.OrgIDScope + "orgID"},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					domain.OrgIDClaim:                     "orgID",
					ClaimResourceOwner + "id":             "orgID",
					ClaimResourceOwner + "name":           "orgName",
					ClaimResourceOwner + "primary_domain": "orgDomain",
				},
			},
		},
		{
			name: "human, roleAudience",
			args: args{
				projectID:    "project1",
				user:         humanUserInfo,
				roleAudience: []string{"project1"},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					ClaimProjectRoles: projectRoles{
						"role1": {"orgID": "orgDomain"},
						"role2": {"orgID": "orgDomain"},
					},
					fmt.Sprintf(ClaimProjectRolesFormat, "project1"): projectRoles{
						"role1": {"orgID": "orgDomain"},
						"role2": {"orgID": "orgDomain"},
					},
				},
			},
		},
		{
			name: "human, requested roles",
			args: args{
				projectID:      "project1",
				user:           humanUserInfo,
				roleAudience:   []string{"project1"},
				requestedRoles: []string{"role2"},
			},
			want: &oidc.UserInfo{
				Claims: map[string]any{
					ClaimProjectRoles: projectRoles{
						"role2": {"orgID": "orgDomain"},
					},
					fmt.Sprintf(ClaimProjectRolesFormat, "project1"): projectRoles{
						"role2": {"orgID": "orgDomain"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assetPrefix := "https://foo.com/assets"
			got := userInfoToOIDC(tt.args.projectID, tt.args.user, tt.args.scope, tt.args.roleAudience, tt.args.requestedRoles, assetPrefix)
			assert.Equal(t, tt.want, got)
		})
	}
}
