package oidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/query"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
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
	userGroups := []query.UserInfoUserGroup{
		{
			Name: "group1",
			ID:   "group1-id",
		},
		{
			Name: "group2",
			ID:   "group2-id",
		},
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
		UserGroups: userGroups,
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
			want: &oidc.UserInfo{
				Subject: "human1",
			},
		},
		{
			name: "machine, empty",
			args: args{
				user: machineUserInfo,
			},
			want: &oidc.UserInfo{
				Subject: "machine1",
			},
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
				Subject: "human1",
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
			want: &oidc.UserInfo{
				Subject: "human1",
			},
		},
		{
			name: "machine, scope email, profileInfoAssertion",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopeEmail},
			},
			want: &oidc.UserInfo{
				Subject: "machine1",
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
				Subject: "human1",
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
				Subject: "machine1",
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
			want: &oidc.UserInfo{
				Subject: "machine1",
			},
		},
		{
			name: "human, scope phone, profileInfoAssertion",
			args: args{
				user:              humanUserInfo,
				userInfoAssertion: true,
				scope:             []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				Subject: "human1",
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
			want: &oidc.UserInfo{
				Subject: "human1",
			},
		},
		{
			name: "machine, scope phone",
			args: args{
				user:  machineUserInfo,
				scope: []string{oidc.ScopePhone},
			},
			want: &oidc.UserInfo{
				Subject:       "machine1",
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
				Subject:       "human1",
				UserInfoEmail: oidc.UserInfoEmail{},
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
			want: &oidc.UserInfo{
				Subject: "machine1",
			},
		},
		{
			name: "machine, scope resource owner",
			args: args{
				user:  machineUserInfo,
				scope: []string{ScopeResourceOwner},
			},
			want: &oidc.UserInfo{
				Subject: "machine1",
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
				Subject: "human1",
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
				Subject: "machine1",
				Claims: map[string]any{
					domain.OrgIDClaim:               "orgID",
					ClaimResourceOwnerID:            "orgID",
					ClaimResourceOwnerName:          "orgName",
					ClaimResourceOwnerPrimaryDomain: "orgDomain",
				},
			},
		},
		{
			name: "human, scope custom user groups, found",
			args: args{
				user:  humanUserInfo,
				scope: []string{ScopeCustomUserGroups},
			},
			want: &oidc.UserInfo{
				Subject: "human1",
				Claims: map[string]any{
					ClaimCustomUserGroups: []query.UserInfoUserGroup{
						{
							Name: "group1",
							ID:   "group1-id",
						},
						{
							Name: "group2",
							ID:   "group2-id",
						},
					},
				},
			},
		},
		{
			name: "human, scope user groups (group names), found",
			args: args{
				user:  humanUserInfo,
				scope: []string{ScopeUserGroups},
			},
			want: &oidc.UserInfo{
				Subject: "human1",
				Claims: map[string]any{
					ClaimUserGroups: []string{"group1", "group2"},
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

func Test_functionForTriggerType(t *testing.T) {
	tests := []struct {
		name        string
		triggerType domain.TriggerType
		want        string
	}{
		{
			name:        "pre userinfo creation",
			triggerType: domain.TriggerTypePreUserinfoCreation,
			want:        exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreUserinfo.LocalizationKey()),
		},
		{
			name:        "pre access token creation",
			triggerType: domain.TriggerTypePreAccessTokenCreation,
			want:        exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreAccessToken.LocalizationKey()),
		},
		{
			name:        "unspecified",
			triggerType: domain.TriggerTypeUnspecified,
			want:        "",
		},
		{
			name:        "post authentication",
			triggerType: domain.TriggerTypePostAuthentication,
			want:        "",
		},
		{
			name:        "pre creation",
			triggerType: domain.TriggerTypePreCreation,
			want:        "",
		},
		{
			name:        "post creation",
			triggerType: domain.TriggerTypePostCreation,
			want:        "",
		},
		{
			name:        "pre saml response creation",
			triggerType: domain.TriggerTypePreSAMLResponseCreation,
			want:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, functionForTriggerType(tt.triggerType))
		})
	}
}

func Test_appendOrLogClaim(t *testing.T) {
	tests := []struct {
		name       string
		userInfo   *oidc.UserInfo
		key        string
		value      any
		wantClaims map[string]any
		wantLogs   []string
	}{
		{
			name:       "new claim is added",
			userInfo:   &oidc.UserInfo{Claims: map[string]any{}},
			key:        "foo",
			value:      "bar",
			wantClaims: map[string]any{"foo": "bar"},
		},
		{
			name:       "reserved prefix is skipped",
			userInfo:   &oidc.UserInfo{Claims: map[string]any{}},
			key:        ClaimPrefix + ":something",
			value:      "bar",
			wantClaims: map[string]any{},
		},
		{
			name:       "existing claim is not overwritten and logs a conflict",
			userInfo:   &oidc.UserInfo{Claims: map[string]any{"foo": "existing"}},
			key:        "foo",
			value:      "bar",
			wantClaims: map[string]any{"foo": "existing"},
			wantLogs:   []string{`key "foo" already exists`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logs []string
			appendOrLogClaim(tt.userInfo, tt.key, tt.value, &logs)
			assert.Equal(t, tt.wantClaims, tt.userInfo.Claims)
			assert.Equal(t, tt.wantLogs, logs)
		})
	}
}

func Test_errorFromRecover(t *testing.T) {
	tests := []struct {
		name    string
		r       any
		wantMsg string
	}{
		{
			name:    "error value is passed through",
			r:       errors.New("boom"),
			wantMsg: "boom",
		},
		{
			name:    "string value is wrapped",
			r:       "boom",
			wantMsg: "boom",
		},
		{
			name:    "other value falls back to a generic message",
			r:       42,
			wantMsg: "unknown error occurred: 42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorFromRecover(tt.r)
			require.Error(t, err)
			assert.Equal(t, tt.wantMsg, err.Error())
		})
	}
}

func Test_runUserinfoActions(t *testing.T) {
	actions.SetLogstoreService(logstore.New[*record.ExecutionLog](nil, nil))

	t.Run("happy path sets claim from script", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
		queriedActions := []*query.Action{
			{
				Name: "testFunc",
				Script: `function testFunc(ctx, api) {
					api.v1.userinfo.setClaim("foo", "bar")
				}`,
			},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", nil, queriedActions)
		require.NoError(t, err)
		assert.Equal(t, "bar", userInfo.Claims["foo"])
	})

	t.Run("panic while building ctxFields is recovered, not left uncaught", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		// A non-JSON-marshalable claim value makes userinfoClaims' eager json.Marshal panic
		// while ctxFields are being built for this (or a later) action - a phase that runs
		// before internal/actions.executeScript's own recover is armed. Without our own
		// recover in the loop, this panic would escape uncaught and skip cancel().
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{"bad": make(chan int)}}
		queriedActions := []*query.Action{
			{Name: "testFunc", Script: `function testFunc(ctx, api) {}`},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", nil, queriedActions)
		require.Error(t, err)
	})

	t.Run("panic from setMetadata bad arg count is recovered", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
		queriedActions := []*query.Action{
			{
				Name: "testFunc",
				Script: `function testFunc(ctx, api) {
					api.v1.user.setMetadata("onlyonearg")
				}`,
			},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", nil, queriedActions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly 2")
	})

	t.Run("error from a script stops the loop and is returned", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
		queriedActions := []*query.Action{
			{
				Name:   "testFunc",
				Script: `function testFunc(ctx, api) { throw new Error("script failure") }`,
			},
			{
				Name: "shouldNotRun",
				Script: `function shouldNotRun(ctx, api) {
					api.v1.userinfo.setClaim("unreachable", "value")
				}`,
			},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", nil, queriedActions)
		require.Error(t, err)
		assert.NotContains(t, userInfo.Claims, "unreachable")
	})

	t.Run("actor delegation chain is readable from the script", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
		actor := &domain.TokenActor{
			UserID: "actor1",
			Issuer: "https://issuer1.example.com",
			Actor: &domain.TokenActor{
				UserID: "actor2",
				Issuer: "https://issuer2.example.com",
			},
		}
		queriedActions := []*query.Action{
			{
				Name: "testFunc",
				Script: `function testFunc(ctx, api) {
					api.v1.claims.setClaim("actor_user", ctx.v1.actor.userId)
					api.v1.claims.setClaim("actor_issuer", ctx.v1.actor.issuer)
					api.v1.claims.setClaim("nested_actor_user", ctx.v1.actor.actor.userId)
					api.v1.claims.setClaim("chain_end", ctx.v1.actor.actor.actor === null)
				}`,
			},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", actor, queriedActions)
		require.NoError(t, err)
		assert.Equal(t, "actor1", userInfo.Claims["actor_user"])
		assert.Equal(t, "https://issuer1.example.com", userInfo.Claims["actor_issuer"])
		assert.Equal(t, "actor2", userInfo.Claims["nested_actor_user"])
		assert.Equal(t, true, userInfo.Claims["chain_end"])
	})

	t.Run("missing actor is null in the script", func(t *testing.T) {
		s := &Server{}
		qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
		userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
		queriedActions := []*query.Action{
			{
				Name: "testFunc",
				Script: `function testFunc(ctx, api) {
					api.v1.claims.setClaim("no_actor", ctx.v1.actor === null)
				}`,
			},
		}

		err := s.runUserinfoActions(context.Background(), qu, userInfo, "clientID", nil, queriedActions)
		require.NoError(t, err)
		assert.Equal(t, true, userInfo.Claims["no_actor"])
	})
}

func Test_runUserinfoExecutionTargets(t *testing.T) {
	respBody := []byte(`{"append_claims":[{"key":"foo","value":"bar"},{"key":"foo","value":"baz"}],"append_log_claims":["extra log entry"]}`)
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	s := &Server{httpClient: server.Client()}
	qu := &query.OIDCUserInfo{User: &query.User{ID: "user1", ResourceOwner: "org1"}}
	userInfo := &oidc.UserInfo{Subject: "user1", Claims: map[string]any{}}
	actor := &domain.TokenActor{
		UserID: "actor1",
		Issuer: "https://issuer1.example.com",
		Actor:  &domain.TokenActor{UserID: "actor2"},
	}
	targets := []target_domain.Target{
		{TargetType: target_domain.TargetTypeCall, Endpoint: server.URL, Timeout: 5 * time.Second},
	}

	err := s.runUserinfoExecutionTargets(context.Background(), qu, userInfo, "clientID", actor, "function/test", targets)
	require.NoError(t, err)

	var sent ContextInfo
	require.NoError(t, json.Unmarshal(gotBody, &sent))
	assert.Equal(t, actor, sent.Actor)

	assert.Equal(t, "bar", userInfo.Claims["foo"])
	logClaim, ok := userInfo.Claims[fmt.Sprintf(ClaimActionLogFormat, "function/test")]
	require.True(t, ok)
	logs, ok := logClaim.([]string)
	require.True(t, ok)
	assert.Contains(t, logs, `key "foo" already exists`)
	assert.Contains(t, logs, "extra log entry")
}
