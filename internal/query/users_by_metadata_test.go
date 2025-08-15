package query

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	permissionmock "github.com/zitadel/zitadel/internal/domain/mock"
	"github.com/zitadel/zitadel/internal/feature"
)

var columns = []string{
	"projections.user_metadata5.key",
	"projections.user_metadata5.value",
	"projections.users14.id",
	"projections.users14.state",
	"projections.users14.username",
	"projections.users14.type",
	"projections.users14_humans.user_id",
	"projections.users14_humans.first_name",
	"projections.users14_humans.last_name",
	"projections.users14_humans.nick_name",
	"projections.users14_humans.display_name",
	"projections.users14_humans.preferred_language",
	"projections.users14_humans.gender",
	"projections.users14_humans.avatar_key",
	"projections.users14_humans.email",
	"projections.users14_humans.is_email_verified",
	"projections.users14_humans.phone",
	"projections.users14_humans.is_phone_verified",
	"projections.users14_humans.password_change_required",
	"projections.users14_humans.password_changed",
	"projections.users14_humans.mfa_init_skipped",
	"projections.users14_machines.user_id",
	"projections.users14_machines.name",
	"projections.users14_machines.description",
	"projections.users14_machines.secret",
	"projections.users14_machines.access_token_type",
	"COUNT(*) OVER ()",
}

const baseQuery = `SELECT projections.user_metadata5.key, projections.user_metadata5.value,
projections.users14.id, projections.users14.state, projections.users14.username, projections.users14.type,
projections.users14_humans.user_id, projections.users14_humans.first_name, projections.users14_humans.last_name, projections.users14_humans.nick_name, projections.users14_humans.display_name, projections.users14_humans.preferred_language, projections.users14_humans.gender, projections.users14_humans.avatar_key, projections.users14_humans.email, projections.users14_humans.is_email_verified, projections.users14_humans.phone, projections.users14_humans.is_phone_verified, projections.users14_humans.password_change_required, projections.users14_humans.password_changed, projections.users14_humans.mfa_init_skipped,
projections.users14_machines.user_id, projections.users14_machines.name, projections.users14_machines.description, projections.users14_machines.secret, projections.users14_machines.access_token_type,
COUNT(*) OVER ()
FROM projections.user_metadata5
JOIN projections.users14 ON projections.user_metadata5.user_id = projections.users14.id
AND projections.user_metadata5.instance_id = projections.users14.instance_id
LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id
AND projections.users14.instance_id = projections.users14_humans.instance_id
LEFT JOIN projections.users14_machines ON projections.users14.id = projections.users14_machines.user_id
AND projections.users14.instance_id = projections.users14_machines.instance_id
`

const permissionQuery = `INNER JOIN eventstore.permitted_orgs($1, $2, $3, $4, $5) permissions ON (permissions.instance_permitted OR projections.users14.resource_owner = ANY(permissions.org_ids) OR projections.users14.id = $6)`

func TestUsersByMetadata(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName string

		inputCtx                 context.Context
		inputQueriesFunc         func(*testing.T) *UsersByMetadataSearchQueries
		inputPermissionCheckMock domain.PermissionCheck

		mockMatcher      sqlmock.QueryMatcher
		mockExpectations func(*sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery

		expectedQuery       string
		expectedUsersByMeta *UsersByMetadata
		expectedError       error
	}{
		{
			testName: "sql should return error",
			inputCtx: authz.NewMockContext("instance-1", "org-1", "user-1"),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ, err := NewUserMetadataKeySearchQuery("test key", TextEquals)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{keyQ},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				return q.WithArgs("test key", "instance-1").WillReturnError(errors.New("mocked error"))
			},
			mockMatcher:   &baseMatcher{},
			expectedQuery: baseQuery,
			expectedError: errors.New("mocked error"),
		},
		{
			testName:                 "v2 perms disabled/matching by key should return 2 records/permission granted",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1"),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ, err := NewUserMetadataKeySearchQuery("test key", TextEquals)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{keyQ},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.WithArgs("test key", "instance-1").WillReturnRows(results)
			},
			mockMatcher:   &whereMatcher{matcher: &baseMatcher{}, toMatch: "user_metadata5.key"},
			expectedQuery: baseQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms disabled/matching by key should return 2 records/permission denied/empty result",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1"),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckErr(errors.New("permission denied")),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ, err := NewUserMetadataKeySearchQuery("test key", TextEquals)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{keyQ},
				}
			},
			mockMatcher: &whereMatcher{matcher: &baseMatcher{}, toMatch: "user_metadata5.key"},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.WithArgs("test key", "instance-1").WillReturnRows(results)
			},
			expectedQuery: baseQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta:    []*UserByMetadata{},
			},
		},
		{
			testName:                 "v2 perms disabled/matching by key ignore case should return 2 records/permission granted",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1"),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ, err := NewUserMetadataKeySearchQuery("Test key", TextEqualsIgnoreCase)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{keyQ},
				}
			},
			mockMatcher: &whereMatcher{matcher: &baseMatcher{}, toMatch: "LOWER(projections.user_metadata5.key)"},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"Test key", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.WithArgs("test key", "instance-1").WillReturnRows(results)
			},
			expectedQuery: baseQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "Test key",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by key ignore case should return 2 records",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ, err := NewUserMetadataKeySearchQuery("Test key", TextEqualsIgnoreCase)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{keyQ},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"Test key", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "test key", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &whereMatcher{matcher: &baseMatcher{}, toMatch: "LOWER(projections.user_metadata5.key)"},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "Test key",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by OR-ed keys should return 2 records",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ1, err := NewUserMetadataKeySearchQuery("test key 1", TextEqualsIgnoreCase)
				require.NoError(t, err)
				keyQ2, err := NewUserMetadataKeySearchQuery("test key 2", TextEqualsIgnoreCase)
				require.NoError(t, err)
				orQ, err := NewOrQuery(keyQ1, keyQ2)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{orQ},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 1", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 2", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "test key 1", "test key 2", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &whereMatcher{matcher: &baseMatcher{}, toMatch: "LOWER(projections.user_metadata5.key),OR"},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key 1",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key 2",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by NOT-OR-ed keys should return 2 records",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ1, err := NewUserMetadataKeySearchQuery("unmatch key 1", TextEquals)
				require.NoError(t, err)
				keyQ2, err := NewUserMetadataKeySearchQuery("unmatch key 2", TextEquals)
				require.NoError(t, err)
				orQ, err := NewOrQuery(keyQ1, keyQ2)
				require.NoError(t, err)

				notQ, err := NewNotQuery(orQ)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{notQ},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 1", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 2", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "unmatch key 1", "unmatch key 2", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &whereMatcher{matcher: &baseMatcher{}, toMatch: "projections.user_metadata5.key,OR,NOT"},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key 1",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key 2",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by NOT-AND-ed keys should return 2 records with limit",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ1, err := NewUserMetadataKeySearchQuery("unmatch key 1", TextEquals)
				require.NoError(t, err)
				keyQ2, err := NewUserMetadataKeySearchQuery("unmatch key 2", TextEquals)
				require.NoError(t, err)
				orQ, err := NewAndQuery(keyQ1, keyQ2)
				require.NoError(t, err)

				notQ, err := NewNotQuery(orQ)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries: []SearchQuery{notQ},
					SearchRequest: SearchRequest{
						Limit: 2,
					},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 1", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 2", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "unmatch key 1", "unmatch key 2", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &limitMatcher{matcher: &whereMatcher{matcher: &baseMatcher{}, toMatch: "projections.user_metadata5.key,AND,NOT"}, expectedLimit: "2"},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key 1",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key 2",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by NOT-AND-ed keys should return 2 records with limit and offset",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ1, err := NewUserMetadataKeySearchQuery("unmatch key 1", TextEquals)
				require.NoError(t, err)
				keyQ2, err := NewUserMetadataKeySearchQuery("unmatch key 2", TextEquals)
				require.NoError(t, err)
				orQ, err := NewAndQuery(keyQ1, keyQ2)
				require.NoError(t, err)

				notQ, err := NewNotQuery(orQ)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries:       []SearchQuery{notQ},
					SearchRequest: SearchRequest{Limit: 2, Offset: 1},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 1", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 2", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "unmatch key 1", "unmatch key 2", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &offsetMatcher{matcher: &limitMatcher{matcher: &whereMatcher{matcher: &baseMatcher{}, toMatch: "projections.user_metadata5.key,AND,NOT"}, expectedLimit: "2"}, expectedOffset: "1"},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key 1",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key 2",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
				},
			},
		},
		{
			testName:                 "v2 perms enabled/matching by NOT-AND-ed keys should return 2 records with limit, offset and sorting",
			inputCtx:                 authz.NewMockContext("instance-1", "org-1", "user-1", authz.WithMockFeatures(feature.Features{PermissionCheckV2: true})),
			inputPermissionCheckMock: permissionmock.MockPermissionCheckOK(),
			inputQueriesFunc: func(t *testing.T) *UsersByMetadataSearchQueries {
				keyQ1, err := NewUserMetadataKeySearchQuery("unmatch key 1", TextEquals)
				require.NoError(t, err)
				keyQ2, err := NewUserMetadataKeySearchQuery("unmatch key 2", TextEquals)
				require.NoError(t, err)
				orQ, err := NewAndQuery(keyQ1, keyQ2)
				require.NoError(t, err)

				notQ, err := NewNotQuery(orQ)
				require.NoError(t, err)

				return &UsersByMetadataSearchQueries{
					Queries:       []SearchQuery{notQ},
					SearchRequest: SearchRequest{Limit: 2, Offset: 1, SortingColumn: UserMetadataKeyCol},
				}
			},
			mockExpectations: func(q *sqlmock.ExpectedQuery) *sqlmock.ExpectedQuery {
				results := sqlmock.NewRows(columns)
				results.AddRows(
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 2", "test value 2", "id-2", domain.UserStateActive, "username 2", domain.UserTypeMachine,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						"id-2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT,
						// COUNT(*) OVER ()
						2,
					},
					[]driver.Value{
						// user_metadata5.key, user_metadata5.value, users14.id, users14.state, users14.username, users14.type
						"test key 1", "test value 1", "id-1", domain.UserStateActive, "username 1", domain.UserTypeHuman,
						// users14_humans.user_id, users14_humans.first_name, users14_humans.last_name, users14_humans.nick_name, users14_humans.display_name, users14_humans.preferred_language, users14_humans.gender, users14_humans.avatar_key, users14_humans.email, users14_humans.is_email_verified, users14_humans.phone, users14_humans.is_phone_verified, users14_humans.password_change_required, users14_humans.password_changed, users14_humans.mfa_init_skipped
						"id-1", "User", "Name", "Nickyname", "User 'Nickyname' Name", "it", domain.GenderMale, "avatar_key", "nicky@username.com", true, "12345678", true, false, nil, nil,
						// users14_machines.user_id, users14_machines.name, users14_machines.description, users14_machines.secret, users14_machines.access_token_type
						nil, nil, nil, nil, nil,
						// COUNT(*) OVER ()
						2,
					},
				)
				return q.
					WithArgs("instance-1", "user-1", sqlmock.AnyArg(), domain.PermissionUserRead, sqlmock.AnyArg(), "user-1", "unmatch key 1", "unmatch key 2", "instance-1").
					WillReturnRows(results)
			},
			mockMatcher:   &orderMatcher{matcher: &offsetMatcher{matcher: &limitMatcher{matcher: &whereMatcher{matcher: &baseMatcher{}, toMatch: "projections.user_metadata5.key,AND,NOT"}, expectedLimit: "2"}, expectedOffset: "1"}, expectedSorting: []string{"projections.user_metadata5.key DESC"}},
			expectedQuery: baseQuery + permissionQuery,
			expectedUsersByMeta: &UsersByMetadata{
				SearchResponse: SearchResponse{Count: 2},
				UsersByMeta: []*UserByMetadata{
					{
						ResourceOwner: "",
						Key:           "test key 2",
						Value:         []byte("test value 2"),
						User:          machineUser("id-2", "username 2", "robot", "this is a robot", "giga secret 1234!", domain.OIDCTokenTypeJWT, domain.UserStateActive),
					},
					{
						ResourceOwner: "",
						Key:           "test key 1",
						Value:         []byte("test value 1"),
						User:          humanUser("id-1", "username 1", "User", "Name", "Nickyname", "nicky@username.com", domain.UserStateActive),
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// t.Parallel()

			// Given
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(tc.mockMatcher))
			require.NoError(t, err)
			q := Queries{
				client: &database.DB{DB: db},
			}
			defer db.Close()

			tc.mockExpectations(mock.ExpectQuery(tc.expectedQuery))

			// Test
			res, err := q.SearchUsersByMetadata(tc.inputCtx, tc.inputQueriesFunc(t), tc.inputPermissionCheckMock)

			// Verify
			require.NoError(t, mock.ExpectationsWereMet())
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUsersByMeta, res)
		})
	}
}

type baseMatcher struct{}

func (m *baseMatcher) Match(expectedSQL, actualSQL string) error {
	strippedExpected := strings.ReplaceAll(expectedSQL, "\n", " ")
	strippedActual := strings.ReplaceAll(actualSQL, "\n", " ")

	beforeWhere, _, _ := strings.Cut(strippedExpected, "WHERE")
	if !strings.Contains(strippedActual, beforeWhere) {
		return errors.New("actual SQL query does not contain expected sql")
	}

	return nil
}

type whereMatcher struct {
	matcher sqlmock.QueryMatcher
	toMatch string
}

func (m *whereMatcher) Match(expectedSQL, actualSQL string) error {
	existingMatchingErrs := m.matcher.Match(expectedSQL, actualSQL)

	strippedActual := strings.ReplaceAll(actualSQL, "\n", " ")

	collector := []error{existingMatchingErrs}

	if !strings.Contains(strippedActual, "WHERE") {
		collector = append(collector, fmt.Errorf(`no WHERE clause found in actual SQL "%s"`, strippedActual))
	}

	for _, v := range strings.Split(m.toMatch, ",") {
		if !strings.Contains(strippedActual, v) {
			collector = append(collector, fmt.Errorf(`value "%s" is not contained in actual SQL "%s"`, v, strippedActual))
		}
	}

	return errors.Join(collector...)
}

type limitMatcher struct {
	expectedLimit string
	matcher       sqlmock.QueryMatcher
}

func (m *limitMatcher) Match(expectedSQL, actualSQL string) error {
	existingMatchingErrs := m.matcher.Match(expectedSQL, actualSQL)

	strippedActual := strings.ReplaceAll(actualSQL, "\n", " ")

	re := regexp.MustCompile(`(?i)LIMIT\s+(\d+)`)
	matches := re.FindStringSubmatch(strippedActual)

	collector := []error{existingMatchingErrs}

	if len(matches) == 0 {
		return errors.Join(append(collector, fmt.Errorf(`no LIMIT clause found in actual SQL "%s"`, strippedActual))...)
	}

	if matches[1] != m.expectedLimit {
		collector = append(collector, fmt.Errorf(`limit "%s" is not contained in actual SQL "%s"`, matches[1], strippedActual))
	}

	return errors.Join(collector...)
}

type offsetMatcher struct {
	expectedOffset string
	matcher        sqlmock.QueryMatcher
}

func (m *offsetMatcher) Match(expectedSQL, actualSQL string) error {
	existingMatchingErrs := m.matcher.Match(expectedSQL, actualSQL)

	strippedActual := strings.ReplaceAll(actualSQL, "\n", " ")

	re := regexp.MustCompile(`(?i)OFFSET\s+(\d+)`)
	matches := re.FindStringSubmatch(strippedActual)

	collector := []error{existingMatchingErrs}

	if len(matches) == 0 {
		return errors.Join(append(collector, fmt.Errorf(`no OFFSET clause found in actual SQL "%s"`, strippedActual))...)
	}

	if matches[1] != m.expectedOffset {
		collector = append(collector, fmt.Errorf(`offset "%s" is not contained in actual SQL "%s"`, matches[1], strippedActual))
	}

	return errors.Join(collector...)
}

type orderMatcher struct {
	expectedSorting []string
	matcher         sqlmock.QueryMatcher
}

func (m *orderMatcher) Match(expectedSQL, actualSQL string) error {
	existingMatchingErrs := m.matcher.Match(expectedSQL, actualSQL)

	strippedActual := strings.ReplaceAll(actualSQL, "\n", " ")

	re := regexp.MustCompile(`(?i)ORDER\s+BY\s+([^\s,]+(?:\s+(?:ASC|DESC))?(?:\s*,\s*[^\s,]+(?:\s+(?:ASC|DESC))?)*)`)
	matches := re.FindStringSubmatch(strippedActual)

	collector := []error{existingMatchingErrs}

	if len(matches) == 0 {
		return errors.Join(append(collector, fmt.Errorf(`no ORDER BY clause found in actual SQL "%s"`, strippedActual))...)
	}

	for _, sortValue := range m.expectedSorting {
		if !strings.Contains(matches[1], sortValue) {
			collector = append(collector, fmt.Errorf(`ORDER BY value "%s" is not contained in actual SQL "%s"`, sortValue, strippedActual))
		}
	}

	return errors.Join(collector...)
}

func machineUser(id, username, machineName, machineDescription, encodedSecret string, tokenType domain.OIDCTokenType, state domain.UserState) *User {
	return &User{
		ID:       id,
		State:    state,
		Type:     domain.UserTypeMachine,
		Username: username,
		Human:    &Human{},
		Machine: &Machine{
			Name:            machineName,
			Description:     machineDescription,
			EncodedSecret:   encodedSecret,
			AccessTokenType: tokenType,
		},
	}
}

func humanUser(id, username, firstName, lastName, nickname string, email domain.EmailAddress, state domain.UserState) *User {
	return &User{
		ID:       id,
		State:    state,
		Type:     domain.UserTypeHuman,
		Username: username,
		Human: &Human{
			FirstName:              firstName,
			LastName:               lastName,
			NickName:               nickname,
			DisplayName:            fmt.Sprintf("%s '%s' %s", firstName, nickname, lastName),
			AvatarKey:              "avatar_key",
			PreferredLanguage:      language.MustParse("it"),
			Gender:                 domain.GenderMale,
			Email:                  email,
			IsEmailVerified:        true,
			Phone:                  "12345678",
			IsPhoneVerified:        true,
			PasswordChangeRequired: false,
		},
		Machine: &Machine{},
	}
}
