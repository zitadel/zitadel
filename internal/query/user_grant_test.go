package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
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
			", roles.role_information" +
			", projections.user_grants5.state" +
			", projections.user_grants5.user_id" +
			", projections.users14.username" +
			", projections.users14.type" +
			", user_orgs.id" +
			", user_orgs.name" +
			", user_orgs.primary_domain" +
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
			" LEFT JOIN projections.orgs1 AS user_orgs ON projections.users14.resource_owner = user_orgs.id AND projections.users14.instance_id = user_orgs.instance_id" +
			" LEFT JOIN projections.project_grants4 ON projections.user_grants5.grant_id = projections.project_grants4.grant_id AND projections.user_grants5.instance_id = projections.project_grants4.instance_id AND projections.project_grants4.project_id = projections.user_grants5.project_id" +
			" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.project_grants4.granted_org_id = granted_orgs.id AND projections.project_grants4.instance_id = granted_orgs.instance_id" +
			" LEFT JOIN projections.login_names3 ON projections.user_grants5.user_id = projections.login_names3.user_id AND projections.user_grants5.instance_id = projections.login_names3.instance_id" +
			" LEFT JOIN LATERAL (SELECT JSON_AGG( JSON_BUILD_OBJECT( 'role_key', pr.role_key, 'display_name', pr.display_name, 'group_name', pr.group_name ) ) as role_information FROM projections.project_roles4 pr WHERE pr.instance_id = projections.user_grants5.instance_id AND pr.project_id = projections.user_grants5.project_id AND pr.role_key = ANY(projections.user_grants5.roles)) as roles ON true " +
			" WHERE projections.login_names3.is_primary = $1")
	userGrantCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"grant_id",
		"roles",
		"role_information",
		"state",
		"user_id",
		"username",
		"type",
		"resource_owner", // user Organization ID
		"name",
		"primary_domain",
		"first_name",
		"last_name",
		"email",
		"display_name",
		"avatar_key",
		"login_name",
		"resource_owner", // user_grant Organization ID
		"name",           // org name
		"primary_domain",
		"project_id",
		"name",           // project name
		"resource_owner", // project_grant Organization ID
		"id",             // granted org id
		"name",           // granted org name
		"primary_domain", // granted org domain
	}
	// userGrantsHasLoginName replaces the login_names3 join of the full query in the page and count statements.
	userGrantsHasLoginName = "EXISTS (SELECT 1 FROM projections.login_names3_users" +
		" WHERE projections.login_names3_users.instance_id = projections.user_grants5.instance_id" +
		" AND projections.login_names3_users.id = projections.user_grants5.user_id)"
	userGrantsSelect = "SELECT projections.user_grants5.id" +
		", projections.user_grants5.creation_date" +
		", projections.user_grants5.change_date" +
		", projections.user_grants5.sequence" +
		", projections.user_grants5.grant_id" +
		", projections.user_grants5.roles" +
		", roles.role_information" +
		", projections.user_grants5.state" +
		", projections.user_grants5.user_id" +
		", projections.users14.username" +
		", projections.users14.type" +
		", user_orgs.id" +
		", user_orgs.name" +
		", user_orgs.primary_domain" +
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
		", (SELECT COUNT(*) FROM projections.user_grants5 WHERE " + userGrantsHasLoginName + " AND projections.user_grants5.instance_id = $1) AS total"
	userGrantsJoins = " LEFT JOIN projections.users14 ON projections.user_grants5.user_id = projections.users14.id AND projections.user_grants5.instance_id = projections.users14.instance_id" +
		" LEFT JOIN projections.users14_humans ON projections.user_grants5.user_id = projections.users14_humans.user_id AND projections.user_grants5.instance_id = projections.users14_humans.instance_id" +
		" LEFT JOIN projections.orgs1 ON projections.user_grants5.resource_owner = projections.orgs1.id AND projections.user_grants5.instance_id = projections.orgs1.instance_id" +
		" LEFT JOIN projections.projects4 ON projections.user_grants5.project_id = projections.projects4.id AND projections.user_grants5.instance_id = projections.projects4.instance_id" +
		" LEFT JOIN projections.orgs1 AS user_orgs ON projections.users14.resource_owner = user_orgs.id AND projections.users14.instance_id = user_orgs.instance_id" +
		" LEFT JOIN projections.project_grants4 ON projections.user_grants5.grant_id = projections.project_grants4.grant_id AND projections.user_grants5.instance_id = projections.project_grants4.instance_id AND projections.project_grants4.project_id = projections.user_grants5.project_id" +
		" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.project_grants4.granted_org_id = granted_orgs.id AND projections.project_grants4.instance_id = granted_orgs.instance_id" +
		" LEFT JOIN projections.login_names3 ON projections.user_grants5.user_id = projections.login_names3.user_id AND projections.user_grants5.instance_id = projections.login_names3.instance_id" +
		" LEFT JOIN LATERAL (SELECT JSON_AGG( JSON_BUILD_OBJECT( 'role_key', pr.role_key, 'display_name', pr.display_name, 'group_name', pr.group_name ) ) as role_information FROM projections.project_roles4 pr WHERE pr.instance_id = projections.user_grants5.instance_id AND pr.project_id = projections.user_grants5.project_id AND pr.role_key = ANY(projections.user_grants5.roles)) as roles ON true "
	// without limit and offset, the grants are joined in a single pass
	userGrantsStmt = regexp.QuoteMeta(userGrantsSelect +
		" FROM projections.user_grants5" + userGrantsJoins +
		" WHERE projections.login_names3.is_primary = $2 AND projections.user_grants5.instance_id = $3")
	// sorted by user id: the page is ordered and cut on the narrow ids, the grant id breaks ties
	userGrantsSortedByUserIDStmt = userGrantsPageStmt(
		"SELECT projections.user_grants5.id, projections.user_grants5.user_id AS sort_col FROM projections.user_grants5 WHERE "+userGrantsHasLoginName+" AND projections.user_grants5.instance_id = $2"+
			" ORDER BY sort_col, projections.user_grants5.id LIMIT 100 OFFSET 9900",
		" ORDER BY page.sort_col, page.id",
	)
	// sorted by grant id: the sorting column is unique, so no tie-break is added
	userGrantsSortedByIDDescStmt = userGrantsPageStmt(
		"SELECT projections.user_grants5.id, projections.user_grants5.id AS sort_col FROM projections.user_grants5 WHERE "+userGrantsHasLoginName+" AND projections.user_grants5.instance_id = $2"+
			" ORDER BY sort_col DESC LIMIT 100 OFFSET 20000",
		" ORDER BY page.sort_col DESC",
	)
	// no sorting: the page keeps the unordered LIMIT of the previous query
	userGrantsLimitStmt = userGrantsPageStmt(
		"SELECT projections.user_grants5.id FROM projections.user_grants5 WHERE "+userGrantsHasLoginName+" AND projections.user_grants5.instance_id = $2 LIMIT 2",
		"",
	)
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
						`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"user-resource-owner",
						"user-resource-owner-name",
						"user-resource-owner-domain",
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
				ID:                      "id",
				CreationDate:            testNow,
				ChangeDate:              testNow,
				Sequence:                20211111,
				Roles:                   database.TextArray[string]{"role-key"},
				RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
				GrantID:                 "grant-id",
				State:                   domain.UserGrantStateActive,
				UserID:                  "user-id",
				Username:                "username",
				UserType:                domain.UserTypeHuman,
				UserResourceOwner:       "user-resource-owner",
				UserResourceOwnerName:   "user-resource-owner-name",
				UserResourceOwnerDomain: "user-resource-owner-domain",
				FirstName:               "first-name",
				LastName:                "last-name",
				Email:                   "email",
				DisplayName:             "display-name",
				AvatarURL:               "avatar-key",
				PreferredLoginName:      "login-name",
				ResourceOwner:           "ro",
				OrgName:                 "org-name",
				OrgPrimaryDomain:        "primary-domain",
				ProjectID:               "project-id",
				ProjectName:             "project-name",
				ProjectResourceOwner:    "project-resource-owner",
				GrantedOrgID:            "granted-org-id",
				GrantedOrgName:          "granted-org-name",
				GrantedOrgDomain:        "granted-org-domain",
			},
		},
		{
			name:    "prepareUserGrantQuery service account found",
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
						`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeMachine,
						"user-resource-owner",
						"user-resource-owner-name",
						"user-resource-owner-domain",
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
				ID:                      "id",
				CreationDate:            testNow,
				ChangeDate:              testNow,
				Sequence:                20211111,
				Roles:                   database.TextArray[string]{"role-key"},
				RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
				GrantID:                 "grant-id",
				State:                   domain.UserGrantStateActive,
				UserID:                  "user-id",
				Username:                "username",
				UserType:                domain.UserTypeMachine,
				UserResourceOwner:       "user-resource-owner",
				UserResourceOwnerName:   "user-resource-owner-name",
				UserResourceOwnerDomain: "user-resource-owner-domain",
				FirstName:               "",
				LastName:                "",
				Email:                   "",
				DisplayName:             "",
				AvatarURL:               "",
				PreferredLoginName:      "login-name",
				ResourceOwner:           "ro",
				OrgName:                 "org-name",
				OrgPrimaryDomain:        "primary-domain",
				ProjectID:               "project-id",
				ProjectName:             "project-name",
				ProjectResourceOwner:    "project-resource-owner",
				GrantedOrgID:            "granted-org-id",
				GrantedOrgName:          "granted-org-name",
				GrantedOrgDomain:        "granted-org-domain",
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
						`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"user-resource-owner",
						"user-resource-owner-name",
						"user-resource-owner-domain",
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
				ID:                      "id",
				CreationDate:            testNow,
				ChangeDate:              testNow,
				Sequence:                20211111,
				Roles:                   database.TextArray[string]{"role-key"},
				RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
				GrantID:                 "grant-id",
				State:                   domain.UserGrantStateActive,
				UserID:                  "user-id",
				Username:                "username",
				UserType:                domain.UserTypeHuman,
				UserResourceOwner:       "user-resource-owner",
				UserResourceOwnerName:   "user-resource-owner-name",
				UserResourceOwnerDomain: "user-resource-owner-domain",
				FirstName:               "first-name",
				LastName:                "last-name",
				Email:                   "email",
				DisplayName:             "display-name",
				AvatarURL:               "avatar-key",
				PreferredLoginName:      "login-name",
				ResourceOwner:           "ro",
				OrgName:                 "",
				OrgPrimaryDomain:        "",
				ProjectID:               "project-id",
				ProjectName:             "project-name",
				ProjectResourceOwner:    "project-resource-owner",
				GrantedOrgID:            "granted-org-id",
				GrantedOrgName:          "granted-org-name",
				GrantedOrgDomain:        "granted-org-domain",
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
						`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"user-resource-owner",
						"user-resource-owner-name",
						"user-resource-owner-domain",
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
				ID:                      "id",
				CreationDate:            testNow,
				ChangeDate:              testNow,
				Sequence:                20211111,
				Roles:                   database.TextArray[string]{"role-key"},
				RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
				GrantID:                 "grant-id",
				State:                   domain.UserGrantStateActive,
				UserID:                  "user-id",
				Username:                "username",
				UserType:                domain.UserTypeHuman,
				UserResourceOwner:       "user-resource-owner",
				UserResourceOwnerName:   "user-resource-owner-name",
				UserResourceOwnerDomain: "user-resource-owner-domain",
				FirstName:               "first-name",
				LastName:                "last-name",
				Email:                   "email",
				DisplayName:             "display-name",
				AvatarURL:               "avatar-key",
				PreferredLoginName:      "login-name",
				ResourceOwner:           "ro",
				OrgName:                 "org-name",
				OrgPrimaryDomain:        "primary-domain",
				ProjectID:               "project-id",
				ProjectName:             "",
				ProjectResourceOwner:    "",
				GrantedOrgID:            "granted-org-id",
				GrantedOrgName:          "granted-org-name",
				GrantedOrgDomain:        "granted-org-domain",
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
						`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
						domain.UserGrantStateActive,
						"user-id",
						"username",
						domain.UserTypeHuman,
						"user-resource-owner",
						"user-resource-owner-name",
						"user-resource-owner-domain",
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
				ID:                      "id",
				CreationDate:            testNow,
				ChangeDate:              testNow,
				Sequence:                20211111,
				Roles:                   database.TextArray[string]{"role-key"},
				RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
				GrantID:                 "grant-id",
				State:                   domain.UserGrantStateActive,
				UserID:                  "user-id",
				Username:                "username",
				UserType:                domain.UserTypeHuman,
				UserResourceOwner:       "user-resource-owner",
				UserResourceOwnerName:   "user-resource-owner-name",
				UserResourceOwnerDomain: "user-resource-owner-domain",
				FirstName:               "first-name",
				LastName:                "last-name",
				Email:                   "email",
				DisplayName:             "display-name",
				AvatarURL:               "avatar-key",
				PreferredLoginName:      "",
				ResourceOwner:           "ro",
				OrgName:                 "org-name",
				OrgPrimaryDomain:        "primary-domain",
				ProjectID:               "project-id",
				ProjectName:             "project-name",
				ProjectResourceOwner:    "project-resource-owner",
				GrantedOrgID:            "granted-org-id",
				GrantedOrgName:          "granted-org-name",
				GrantedOrgDomain:        "granted-org-domain",
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
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (service account)",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeMachine,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeMachine,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "",
						LastName:                "",
						Email:                   "",
						DisplayName:             "",
						AvatarURL:               "",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no org)",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeMachine,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeMachine,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "",
						OrgPrimaryDomain:        "",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no project)",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "",
						ProjectResourceOwner:    "",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareUserGrantsQuery one grant (no loginname)",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name: "prepareUserGrantsQuery multiple grants, limit without sorting",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{
				SearchRequest: SearchRequest{Limit: 2},
			}),
			want: want{
				sqlExpectations: mockQueries(
					userGrantsLimitStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
							`[{"display_name":"displayName","group_name":"groupName","role_key":"role-key"}]`,
							domain.UserGrantStateActive,
							"user-id",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
					{
						ID:                      "id",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						RoleInformation:         []Role{{Key: "role-key", DisplayName: "displayName", GroupName: "groupName"}},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-id",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name: "prepareUserGrantsQuery sorted by user id, limit and offset",
			prepare: prepareUserGrantsQueryFor(authz.WithInstanceID(context.Background(), "instance-id"), &UserGrantsQueries{
				SearchRequest: SearchRequest{
					Offset:        9900,
					Limit:         100,
					Asc:           true,
					SortingColumn: UserGrantUserID,
				},
			}),
			want: want{
				sqlExpectations: mockQueries(
					userGrantsSortedByUserIDStmt,
					userGrantsCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							nil,
							domain.UserGrantStateActive,
							"user-1",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
							uint64(10000),
						},
						{
							"id-2",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							nil,
							domain.UserGrantStateActive,
							"user-1",
							"username",
							domain.UserTypeHuman,
							"user-resource-owner",
							"user-resource-owner-name",
							"user-resource-owner-domain",
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
							uint64(10000),
						},
					},
					"instance-id", "instance-id", "instance-id", true,
				),
			},
			// the total comes from the scalar count, not from the rows of the page
			object: &UserGrants{
				SearchResponse: SearchResponse{
					Count: 10000,
				},
				UserGrants: []*UserGrant{
					{
						ID:                      "id-1",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-1",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
					{
						ID:                      "id-2",
						CreationDate:            testNow,
						ChangeDate:              testNow,
						Sequence:                20211111,
						Roles:                   database.TextArray[string]{"role-key"},
						GrantID:                 "grant-id",
						State:                   domain.UserGrantStateActive,
						UserID:                  "user-1",
						Username:                "username",
						UserType:                domain.UserTypeHuman,
						UserResourceOwner:       "user-resource-owner",
						UserResourceOwnerName:   "user-resource-owner-name",
						UserResourceOwnerDomain: "user-resource-owner-domain",
						FirstName:               "first-name",
						LastName:                "last-name",
						Email:                   "email",
						DisplayName:             "display-name",
						AvatarURL:               "avatar-key",
						PreferredLoginName:      "login-name",
						ResourceOwner:           "ro",
						OrgName:                 "org-name",
						OrgPrimaryDomain:        "primary-domain",
						ProjectID:               "project-id",
						ProjectName:             "project-name",
						ProjectResourceOwner:    "project-resource-owner",
						GrantedOrgID:            "granted-org-id",
						GrantedOrgName:          "granted-org-name",
						GrantedOrgDomain:        "granted-org-domain",
					},
				},
			},
		},
		{
			name: "prepareUserGrantsQuery sorted by id descending, offset past the last grant",
			prepare: prepareUserGrantsQueryFor(authz.WithInstanceID(context.Background(), "instance-id"), &UserGrantsQueries{
				SearchRequest: SearchRequest{
					Offset:        20000,
					Limit:         100,
					SortingColumn: UserGrantID,
				},
			}),
			want: want{
				sqlExpectations: mockQueries(
					userGrantsSortedByIDDescStmt,
					nil,
					nil,
					"instance-id", "instance-id", "instance-id", true,
				),
			},
			object: &UserGrants{UserGrants: []*UserGrant{}},
		},
		{
			name:    "prepareUserGrantsQuery sql err",
			prepare: prepareUserGrantsQueryFor(context.Background(), &UserGrantsQueries{}),
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

func prepareUserGrantsQueryFor(ctx context.Context, queries *UserGrantsQueries) func() (sq.SelectBuilder, func(*sql.Rows) (*UserGrants, error)) {
	return func() (sq.SelectBuilder, func(*sql.Rows) (*UserGrants, error)) {
		return queries.prepareUserGrantsQuery(ctx, false)
	}
}

func TestUserGrantsQueries_filterJoins(t *testing.T) {
	displayName, err := NewUserGrantDisplayNameQuery("display", TextContains)
	require.NoError(t, err)
	withGranted, err := NewUserGrantWithGrantedQuery("org-id")
	require.NoError(t, err)
	userOwner, err := NewUserGrantUserResourceOwnerSearchQuery("org-id")
	require.NoError(t, err)
	state, err := NewUserGrantStateQuery(domain.UserGrantStateActive)
	require.NoError(t, err)
	userOrgName, err := NewTextQuery(UserOrgColumnName, "org", TextEquals)
	require.NoError(t, err)
	grantedOrgName, err := NewTextQuery(GrantedOrgColumnName, "org", TextEquals)
	require.NoError(t, err)
	loginName, err := NewTextQuery(LoginNameNameCol, "login", TextEquals)
	require.NoError(t, err)
	unknown, err := NewTextQuery(Column{name: "name", table: table{name: "projections.unknown"}}, "x", TextEquals)
	require.NoError(t, err)
	notUnknown, err := NewNotQuery(unknown)
	require.NoError(t, err)

	tests := []struct {
		name    string
		queries []SearchQuery
		want    userGrantJoins
	}{
		{
			name: "no queries",
		},
		{
			name:    "user grant columns only",
			queries: []SearchQuery{state},
		},
		{
			name:    "human and project (in or query)",
			queries: []SearchQuery{displayName, withGranted},
			want:    userGrantJoins{humans: true, projects: true},
		},
		{
			name:    "user",
			queries: []SearchQuery{userOwner},
			want:    userGrantJoins{users: true},
		},
		{
			name:    "user org needs user",
			queries: []SearchQuery{userOrgName},
			want:    userGrantJoins{users: true, userOrgs: true},
		},
		{
			name:    "granted org needs project grant",
			queries: []SearchQuery{grantedOrgName},
			want:    userGrantJoins{projectGrants: true, grantedOrgs: true},
		},
		{
			name:    "login name",
			queries: []SearchQuery{loginName},
			want:    userGrantJoins{loginNames: true},
		},
		{
			name:    "unknown table joins everything",
			queries: []SearchQuery{state, notUnknown},
			want:    allUserGrantJoins(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &UserGrantsQueries{Queries: tt.queries}
			assert.Equal(t, tt.want, q.filterJoins())
		})
	}
}

func TestUserGrantsQueries_prepareUserGrantsQuery_narrowJoins(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance-id")
	displayName, err := NewUserGrantDisplayNameQuery("display", TextEquals)
	require.NoError(t, err)
	loginName, err := NewTextQuery(LoginNameNameCol, "login", TextEquals)
	require.NoError(t, err)

	const (
		usersJoin      = "LEFT JOIN projections.users14 ON "
		humansJoin     = "LEFT JOIN projections.users14_humans ON "
		loginNamesJoin = "LEFT JOIN projections.login_names3 ON "
		rolesJoin      = "LEFT JOIN LATERAL"
		isPrimary      = "projections.login_names3.is_primary = ?"
	)

	tests := []struct {
		name          string
		queries       *UserGrantsQueries
		wantPage      []string
		wantNotInPage []string
		wantCount     []string
		wantNotCount  []string
		wantArgs      []any
	}{
		{
			name: "filter joins in page and count, sort join in page only",
			queries: &UserGrantsQueries{
				SearchRequest: SearchRequest{Limit: 10, Offset: 20, SortingColumn: UserResourceOwnerCol},
				Queries:       []SearchQuery{displayName},
			},
			wantPage:      []string{usersJoin, humansJoin, userGrantsHasLoginName, "ORDER BY sort_col DESC, projections.user_grants5.id DESC LIMIT 10 OFFSET 20"},
			wantNotInPage: []string{loginNamesJoin, rolesJoin, "OVER ()"},
			wantCount:     []string{humansJoin, userGrantsHasLoginName, "projections.users14_humans.display_name = ?"},
			wantNotCount:  []string{usersJoin, loginNamesJoin, rolesJoin},
			wantArgs:      []any{"display", "instance-id", "display", "instance-id", "instance-id", true},
		},
		{
			name: "login name filter keeps the login_names3 join and is_primary",
			queries: &UserGrantsQueries{
				SearchRequest: SearchRequest{Limit: 10},
				Queries:       []SearchQuery{loginName},
			},
			wantPage:      []string{loginNamesJoin, isPrimary},
			wantNotInPage: []string{userGrantsHasLoginName, usersJoin, rolesJoin},
			wantCount:     []string{loginNamesJoin, isPrimary},
			wantNotCount:  []string{userGrantsHasLoginName},
			wantArgs:      []any{true, "login", "instance-id", true, "login", "instance-id", "instance-id", true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterJoins := tt.queries.filterJoins()
			countStmt, _, err := tt.queries.matchingGrants(ctx, false, filterJoins, "COUNT(*)").ToSql()
			require.NoError(t, err)
			for _, want := range tt.wantCount {
				assert.Contains(t, countStmt, want)
			}
			for _, notWant := range tt.wantNotCount {
				assert.NotContains(t, countStmt, notWant)
			}

			query, _ := tt.queries.prepareUserGrantsQuery(ctx, false)
			stmt, args, err := query.PlaceholderFormat(sq.Question).ToSql()
			require.NoError(t, err)
			pageStart := strings.Index(stmt, "FROM (SELECT ")
			pageEnd := strings.Index(stmt, ") AS page")
			require.True(t, pageStart > 0 && pageEnd > pageStart, stmt)
			page := stmt[pageStart:pageEnd]
			for _, want := range tt.wantPage {
				assert.Contains(t, page, want)
			}
			for _, notWant := range tt.wantNotInPage {
				assert.NotContains(t, page, notWant)
			}
			assert.Equal(t, tt.wantArgs, args)
		})
	}
}

func userGrantsPageStmt(page, orderBy string) string {
	return regexp.QuoteMeta(userGrantsSelect +
		" FROM (" + page + ") AS page" +
		" JOIN projections.user_grants5 ON projections.user_grants5.id = page.id AND projections.user_grants5.instance_id = $3" +
		userGrantsJoins +
		" WHERE projections.login_names3.is_primary = $4" + orderBy)
}
