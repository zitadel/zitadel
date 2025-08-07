package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	userGrantStmt = regexp.QuoteMeta(
		"SELECT projections.user_grants5.id" +
			", projections.user_grants5.creation_date" +
			", projections.user_grants5.change_date" +
			", projections.user_grants5.sequence" +
			", projections.user_grants5.grant_id" +
			", projections.user_grants5.roles" +
			", projections.user_grants5.state" +
			", projections.user_grants5.user_id" +
			", projections.users14.username" +
			", projections.users14.type" +
			", projections.users14.resource_owner" +
			", projections.users14_humans.first_name" +
			", projections.users14_humans.last_name" +
			", projections.users14_humans.email" +
			", projections.users14_humans.display_name" +
			", projections.users14_humans.avatar_key" +
			", projections.login_names3.login_name" +
			", projections.user_grants5.resource_owner" +
			", projections.orgs1.name" +
			", projections.orgs1.primary_domain" +
			", projections.user_grants5.project_id" +
			", projections.projects4.name" +
			", projections.projects4.resource_owner" +
			", granted_orgs.id" +
			", granted_orgs.name" +
			", granted_orgs.primary_domain" +
			" FROM projections.user_grants5" +
			" LEFT JOIN projections.users14 ON projections.user_grants5.user_id = projections.users14.id AND projections.user_grants5.instance_id = projections.users14.instance_id" +
			" LEFT JOIN projections.users14_humans ON projections.user_grants5.user_id = projections.users14_humans.user_id AND projections.user_grants5.instance_id = projections.users14_humans.instance_id" +
			" LEFT JOIN projections.orgs1 ON projections.user_grants5.resource_owner = projections.orgs1.id AND projections.user_grants5.instance_id = projections.orgs1.instance_id" +
			" LEFT JOIN projections.projects4 ON projections.user_grants5.project_id = projections.projects4.id AND projections.user_grants5.instance_id = projections.projects4.instance_id" +
			" LEFT JOIN projections.project_grants4 ON projections.user_grants5.grant_id = projections.project_grants4.grant_id AND projections.user_grants5.instance_id = projections.project_grants4.instance_id AND projections.project_grants4.project_id = projections.user_grants5.project_id" +
			" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.project_grants4.granted_org_id = granted_orgs.id AND projections.project_grants4.instance_id = granted_orgs.instance_id" +
			" LEFT JOIN projections.login_names3 ON projections.user_grants5.user_id = projections.login_names3.user_id AND projections.user_grants5.instance_id = projections.login_names3.instance_id" +
			" WHERE projections.login_names3.is_primary = $1")
	userGrantCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"grant_id",
		"roles",
		"state",
		"user_id",
		"username",
		"type",
		"resource_owner", //user resource owner
		"first_name",
		"last_name",
		"email",
		"display_name",
		"avatar_key",
		"login_name",
		"resource_owner", //user_grant resource owner
		"name",           //org name
		"primary_domain",
		"project_id",
		"name",           // project name
		"resource_owner", // project_grant resource owner
		"id",             // granted org id
		"name",           // granted org name
		"primary_domain", // granted org domain
	}
	userGrantsStmt = regexp.QuoteMeta(
		"SELECT projections.user_grants5.id" +
			", projections.user_grants5.creation_date" +
			", projections.user_grants5.change_date" +
			", projections.user_grants5.sequence" +
			", projections.user_grants5.grant_id" +
			", projections.user_grants5.roles" +
			", projections.user_grants5.state" +
			", projections.user_grants5.user_id" +
			", projections.users14.username" +
			", projections.users14.type" +
			", projections.users14.resource_owner" +
			", projections.users14_humans.first_name" +
			", projections.users14_humans.last_name" +
			", projections.users14_humans.email" +
			", projections.users14_humans.display_name" +
			", projections.users14_humans.avatar_key" +
			", projections.login_names3.login_name" +
			", projections.user_grants5.resource_owner" +
			", projections.orgs1.name" +
			", projections.orgs1.primary_domain" +
			", projections.user_grants5.project_id" +
			", projections.projects4.name" +
			", projections.projects4.resource_owner" +
			", granted_orgs.id" +
			", granted_orgs.name" +
			", granted_orgs.primary_domain" +
			", COUNT(*) OVER ()" +
			" FROM projections.user_grants5" +
			" LEFT JOIN projections.users14 ON projections.user_grants5.user_id = projections.users14.id AND projections.user_grants5.instance_id = projections.users14.instance_id" +
			" LEFT JOIN projections.users14_humans ON projections.user_grants5.user_id = projections.users14_humans.user_id AND projections.user_grants5.instance_id = projections.users14_humans.instance_id" +
			" LEFT JOIN projections.orgs1 ON projections.user_grants5.resource_owner = projections.orgs1.id AND projections.user_grants5.instance_id = projections.orgs1.instance_id" +
			" LEFT JOIN projections.projects4 ON projections.user_grants5.project_id = projections.projects4.id AND projections.user_grants5.instance_id = projections.projects4.instance_id" +
			" LEFT JOIN projections.project_grants4 ON projections.user_grants5.grant_id = projections.project_grants4.grant_id AND projections.user_grants5.instance_id = projections.project_grants4.instance_id AND projections.project_grants4.project_id = projections.user_grants5.project_id" +
			" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.project_grants4.granted_org_id = granted_orgs.id AND projections.project_grants4.instance_id = granted_orgs.instance_id" +
			" LEFT JOIN projections.login_names3 ON projections.user_grants5.user_id = projections.login_names3.user_id AND projections.user_grants5.instance_id = projections.login_names3.instance_id" +
			" WHERE projections.login_names3.is_primary = $1")
	userGrantsCols = append(
		userGrantCols,
		"count",
	)
)

func Test_UserGrantPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareUserGrantQuery no result",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					userGrantStmt,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserGrant)(nil),
		},
		{
			name:    "prepareUserGrantQuery found",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					userGrantStmt,
					userGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"resource-owner",
						"first-name",
						"last-name",
						"email",
						"display-name",
						"avatar-key",
						"login-name",
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						"project-name",
						"project-resource-owner",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &UserGrant{
				ID:                   "id",
				CreationDate:         testNow,
				ChangeDate:           testNow,
				Sequence:             20211111,
				Roles:                database.TextArray[string]{"role-key"},
				GrantID:              "grant-id",
				State:                domain.UserGrantStateActive,
				UserID:               "user-id",
				Username:             "username",
				UserType:             domain.UserTypeHuman,
				UserResourceOwner:    "resource-owner",
				FirstName:            "first-name",
				LastName:             "last-name",
				Email:                "email",
				DisplayName:          "display-name",
				AvatarURL:            "avatar-key",
				PreferredLoginName:   "login-name",
				ResourceOwner:        "ro",
				OrgName:              "org-name",
				OrgPrimaryDomain:     "primary-domain",
				ProjectID:            "project-id",
				ProjectName:          "project-name",
				ProjectResourceOwner: "project-resource-owner",
				GrantedOrgID:         "granted-org-id",
				GrantedOrgName:       "granted-org-name",
				GrantedOrgDomain:     "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery machine user found",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					userGrantStmt,
					userGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeMachine,
						"resource-owner",
						nil,
						nil,
						nil,
						nil,
						nil,
						"login-name",
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						"project-name",
						"project-resource-owner",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &UserGrant{
				ID:                   "id",
				CreationDate:         testNow,
				ChangeDate:           testNow,
				Sequence:             20211111,
				Roles:                database.TextArray[string]{"role-key"},
				GrantID:              "grant-id",
				State:                domain.UserGrantStateActive,
				UserID:               "user-id",
				Username:             "username",
				UserType:             domain.UserTypeMachine,
				UserResourceOwner:    "resource-owner",
				FirstName:            "",
				LastName:             "",
				Email:                "",
				DisplayName:          "",
				AvatarURL:            "",
				PreferredLoginName:   "login-name",
				ResourceOwner:        "ro",
				OrgName:              "org-name",
				OrgPrimaryDomain:     "primary-domain",
				ProjectID:            "project-id",
				ProjectName:          "project-name",
				ProjectResourceOwner: "project-resource-owner",
				GrantedOrgID:         "granted-org-id",
				GrantedOrgName:       "granted-org-name",
				GrantedOrgDomain:     "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery (no org) found",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					userGrantStmt,
					userGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"resource-owner",
						"first-name",
						"last-name",
						"email",
						"display-name",
						"avatar-key",
						"login-name",
						"ro",
						nil,
						nil,
						"project-id",
						"project-name",
						"project-resource-owner",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &UserGrant{
				ID:                   "id",
				CreationDate:         testNow,
				ChangeDate:           testNow,
				Sequence:             20211111,
				Roles:                database.TextArray[string]{"role-key"},
				GrantID:              "grant-id",
				State:                domain.UserGrantStateActive,
				UserID:               "user-id",
				Username:             "username",
				UserType:             domain.UserTypeHuman,
				UserResourceOwner:    "resource-owner",
				FirstName:            "first-name",
				LastName:             "last-name",
				Email:                "email",
				DisplayName:          "display-name",
				AvatarURL:            "avatar-key",
				PreferredLoginName:   "login-name",
				ResourceOwner:        "ro",
				OrgName:              "",
				OrgPrimaryDomain:     "",
				ProjectID:            "project-id",
				ProjectName:          "project-name",
				ProjectResourceOwner: "project-resource-owner",
				GrantedOrgID:         "granted-org-id",
				GrantedOrgName:       "granted-org-name",
				GrantedOrgDomain:     "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery (no project) found",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					userGrantStmt,
					userGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"resource-owner",
						"first-name",
						"last-name",
						"email",
						"display-name",
						"avatar-key",
						"login-name",
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						nil,
						nil,
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &UserGrant{
				ID:                   "id",
				CreationDate:         testNow,
				ChangeDate:           testNow,
				Sequence:             20211111,
				Roles:                database.TextArray[string]{"role-key"},
				GrantID:              "grant-id",
				State:                domain.UserGrantStateActive,
				UserID:               "user-id",
				Username:             "username",
				UserType:             domain.UserTypeHuman,
				UserResourceOwner:    "resource-owner",
				FirstName:            "first-name",
				LastName:             "last-name",
				Email:                "email",
				DisplayName:          "display-name",
				AvatarURL:            "avatar-key",
				PreferredLoginName:   "login-name",
				ResourceOwner:        "ro",
				OrgName:              "org-name",
				OrgPrimaryDomain:     "primary-domain",
				ProjectID:            "project-id",
				ProjectName:          "",
				ProjectResourceOwner: "",
				GrantedOrgID:         "granted-org-id",
				GrantedOrgName:       "granted-org-name",
				GrantedOrgDomain:     "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery (no loginname) found",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					userGrantStmt,
					userGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"resource-owner",
						"first-name",
						"last-name",
						"email",
						"display-name",
						"avatar-key",
						nil,
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						"project-name",
						"project-resource-owner",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &UserGrant{
				ID:                   "id",
				CreationDate:         testNow,
				ChangeDate:           testNow,
				Sequence:             20211111,
				Roles:                database.TextArray[string]{"role-key"},
				GrantID:              "grant-id",
				State:                domain.UserGrantStateActive,
				UserID:               "user-id",
				Username:             "username",
				UserType:             domain.UserTypeHuman,
				UserResourceOwner:    "resource-owner",
				FirstName:            "first-name",
				LastName:             "last-name",
				Email:                "email",
				DisplayName:          "display-name",
				AvatarURL:            "avatar-key",
				PreferredLoginName:   "",
				ResourceOwner:        "ro",
				OrgName:              "org-name",
				OrgPrimaryDomain:     "primary-domain",
				ProjectID:            "project-id",
				ProjectName:          "project-name",
				ProjectResourceOwner: "project-resource-owner",
				GrantedOrgID:         "granted-org-id",
				GrantedOrgName:       "granted-org-name",
				GrantedOrgDomain:     "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery sql err",
			prepare: prepareUserGrantQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					userGrantStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserGrant)(nil),
		},
		{
			name:    "prepareUserGrantsQuery no result",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					nil,
					nil,
				),
			},
			object: &UserGrants{UserGrants: []*UserGrant{}},
		},
		{
			name:    "prepareUserGrantsQuery one grant",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							"login-name",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeHuman,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (machine user)",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeMachine,
							"resource-owner",
							nil,
							nil,
							nil,
							nil,
							nil,
							"login-name",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeMachine,
						UserResourceOwner:    "resource-owner",
						FirstName:            "",
						LastName:             "",
						Email:                "",
						DisplayName:          "",
						AvatarURL:            "",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no org)",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeMachine,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							"login-name",
							"ro",
							nil,
							nil,
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeMachine,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "",
						OrgPrimaryDomain:     "",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no project)",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							"login-name",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							nil,
							nil,
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeHuman,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "",
						ProjectResourceOwner: "",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no loginname)",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							nil,
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeHuman,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery multiple grants",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					userGrantsStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							"login-name",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"resource-owner",
							"first-name",
							"last-name",
							"email",
							"display-name",
							"avatar-key",
							"login-name",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"project-resource-owner",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				UserGrants: []*UserGrant{
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeHuman,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
					{
						ID:                   "id",
						CreationDate:         testNow,
						ChangeDate:           testNow,
						Sequence:             20211111,
						Roles:                database.TextArray[string]{"role-key"},
						GrantID:              "grant-id",
						State:                domain.UserGrantStateActive,
						UserID:               "user-id",
						Username:             "username",
						UserType:             domain.UserTypeHuman,
						UserResourceOwner:    "resource-owner",
						FirstName:            "first-name",
						LastName:             "last-name",
						Email:                "email",
						DisplayName:          "display-name",
						AvatarURL:            "avatar-key",
						PreferredLoginName:   "login-name",
						ResourceOwner:        "ro",
						OrgName:              "org-name",
						OrgPrimaryDomain:     "primary-domain",
						ProjectID:            "project-id",
						ProjectName:          "project-name",
						ProjectResourceOwner: "project-resource-owner",
						GrantedOrgID:         "granted-org-id",
						GrantedOrgName:       "granted-org-name",
						GrantedOrgDomain:     "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery sql err",
			prepare: prepareUserGrantsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					userGrantsStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserGrants)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
