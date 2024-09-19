package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

func TestUser_authMethodsCheckPermission(t *testing.T) {
	type want struct {
		methods []*AuthMethod
	}
	type args struct {
		user    string
		methods *AuthMethods
	}
	tests := []struct {
		name        string
		args        args
		want        want
		permissions []string
	}{
		{
			"permissions for all users",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
				},
			},
			[]string{"first", "second", "third"},
		},
		{
			"permissions for one user, first",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "first"},
				},
			},
			[]string{"first"},
		},
		{
			"permissions for one user, second",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "second"},
				},
			},
			[]string{"second"},
		},
		{
			"permissions for one user, third",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "third"},
				},
			},
			[]string{"third"},
		},
		{
			"permissions for two users, first",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "first"}, {UserID: "third"},
				},
			},
			[]string{"first", "third"},
		},
		{
			"permissions for two users, second",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{
					{UserID: "second"}, {UserID: "third"},
				},
			},
			[]string{"second", "third"},
		},
		{
			"no permissions",
			args{
				"none",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{},
			},
			[]string{},
		},
		{
			"no permissions, self",
			args{
				"second",
				&AuthMethods{
					AuthMethods: []*AuthMethod{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				methods: []*AuthMethod{{UserID: "second"}},
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				for _, perm := range tt.permissions {
					if resourceID == perm {
						return nil
					}
				}
				return errors.New("failed")
			}
			authMethodsCheckPermission(authz.SetCtxData(context.Background(), authz.CtxData{UserID: tt.args.user}), tt.args.methods, checkPermission)
			require.Equal(t, tt.want.methods, tt.args.methods.AuthMethods)
		})
	}
}

var (
	prepareUserAuthMethodsStmt = `SELECT projections.user_auth_methods4.token_id,` +
		` projections.user_auth_methods4.creation_date,` +
		` projections.user_auth_methods4.change_date,` +
		` projections.user_auth_methods4.resource_owner,` +
		` projections.user_auth_methods4.user_id,` +
		` projections.user_auth_methods4.sequence,` +
		` projections.user_auth_methods4.name,` +
		` projections.user_auth_methods4.state,` +
		` projections.user_auth_methods4.method_type,` +
		` COUNT(*) OVER ()` +
		` FROM projections.user_auth_methods4` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareUserAuthMethodsCols = []string{
		"token_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"user_id",
		"sequence",
		"name",
		"state",
		"method_type",
		"count",
	}
	prepareActiveAuthMethodTypesStmt = `SELECT projections.users13_notifications.password_set,` +
		` auth_method_types.method_type,` +
		` user_idps_count.count` +
		` FROM projections.users13` +
		` LEFT JOIN projections.users13_notifications ON projections.users13.id = projections.users13_notifications.user_id AND projections.users13.instance_id = projections.users13_notifications.instance_id` +
		` LEFT JOIN (SELECT DISTINCT(auth_method_types.method_type), auth_method_types.user_id, auth_method_types.instance_id FROM projections.user_auth_methods4 AS auth_method_types` +
		` WHERE auth_method_types.state = $1) AS auth_method_types` +
		` ON auth_method_types.user_id = projections.users13.id AND auth_method_types.instance_id = projections.users13.instance_id` +
		` LEFT JOIN (SELECT user_idps_count.user_id, user_idps_count.instance_id, COUNT(user_idps_count.user_id) AS count FROM projections.idp_user_links3 AS user_idps_count` +
		` GROUP BY user_idps_count.user_id, user_idps_count.instance_id) AS user_idps_count` +
		` ON user_idps_count.user_id = projections.users13.id AND user_idps_count.instance_id = projections.users13.instance_id` +
		` AS OF SYSTEM TIME '-1 ms`
	prepareActiveAuthMethodTypesCols = []string{
		"password_set",
		"method_type",
		"idps_count",
	}
	prepareAuthMethodTypesRequiredStmt = `SELECT projections.users13.type,` +
		` auth_methods_force_mfa.force_mfa,` +
		` auth_methods_force_mfa.force_mfa_local_only` +
		` FROM projections.users13` +
		` LEFT JOIN (SELECT auth_methods_force_mfa.force_mfa, auth_methods_force_mfa.force_mfa_local_only, auth_methods_force_mfa.instance_id, auth_methods_force_mfa.aggregate_id, auth_methods_force_mfa.is_default FROM projections.login_policies5 AS auth_methods_force_mfa) AS auth_methods_force_mfa` +
		` ON (auth_methods_force_mfa.aggregate_id = projections.users13.instance_id OR auth_methods_force_mfa.aggregate_id = projections.users13.resource_owner) AND auth_methods_force_mfa.instance_id = projections.users13.instance_id` +
		` ORDER BY auth_methods_force_mfa.is_default LIMIT 1
`
	prepareAuthMethodTypesRequiredCols = []string{
		"type",
		"force_mfa",
		"force_mfa_local_only",
	}
)

func Test_UserAuthMethodPrepares(t *testing.T) {
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
			name:    "prepareUserAuthMethodsQuery no result",
			prepare: prepareUserAuthMethodsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserAuthMethodsStmt),
					nil,
					nil,
				),
			},
			object: &AuthMethods{AuthMethods: []*AuthMethod{}},
		},
		{
			name:    "prepareUserAuthMethodsQuery one result",
			prepare: prepareUserAuthMethodsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserAuthMethodsStmt),
					prepareUserAuthMethodsCols,
					[][]driver.Value{
						{
							"token_id",
							testNow,
							testNow,
							"ro",
							"user_id",
							uint64(20211108),
							"name",
							domain.MFAStateReady,
							domain.UserAuthMethodTypeU2F,
						},
					},
				),
			},
			object: &AuthMethods{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				AuthMethods: []*AuthMethod{
					{
						TokenID:       "token_id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						UserID:        "user_id",
						Sequence:      20211108,
						Name:          "name",
						State:         domain.MFAStateReady,
						Type:          domain.UserAuthMethodTypeU2F,
					},
				},
			},
		},
		{
			name:    "prepareUserAuthMethodsQuery multiple result",
			prepare: prepareUserAuthMethodsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserAuthMethodsStmt),
					prepareUserAuthMethodsCols,
					[][]driver.Value{
						{
							"token_id",
							testNow,
							testNow,
							"ro",
							"user_id",
							uint64(20211108),
							"name",
							domain.MFAStateReady,
							domain.UserAuthMethodTypeU2F,
						},
						{
							"token_id-2",
							testNow,
							testNow,
							"ro",
							"user_id",
							uint64(20211108),
							"name-2",
							domain.MFAStateReady,
							domain.UserAuthMethodTypePasswordless,
						},
					},
				),
			},
			object: &AuthMethods{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				AuthMethods: []*AuthMethod{
					{
						TokenID:       "token_id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						UserID:        "user_id",
						Sequence:      20211108,
						Name:          "name",
						State:         domain.MFAStateReady,
						Type:          domain.UserAuthMethodTypeU2F,
					},
					{
						TokenID:       "token_id-2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						UserID:        "user_id",
						Sequence:      20211108,
						Name:          "name-2",
						State:         domain.MFAStateReady,
						Type:          domain.UserAuthMethodTypePasswordless,
					},
				},
			},
		},
		{
			name:    "prepareUserAuthMethodsQuery sql err",
			prepare: prepareUserAuthMethodsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareUserAuthMethodsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*AuthMethodTypes)(nil),
		},
		{
			name: "prepareUserAuthMethodTypesQuery no result",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
				builder, scan := prepareUserAuthMethodTypesQuery(ctx, db, true)
				return builder, func(rows *sql.Rows) (*AuthMethodTypes, error) {
					return scan(rows)
				}
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareActiveAuthMethodTypesStmt),
					nil,
					nil,
				),
			},
			object: &AuthMethodTypes{AuthMethodTypes: []domain.UserAuthMethodType{}},
		},
		{
			name: "prepareUserAuthMethodTypesQuery one second factor",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
				builder, scan := prepareUserAuthMethodTypesQuery(ctx, db, true)
				return builder, func(rows *sql.Rows) (*AuthMethodTypes, error) {
					return scan(rows)
				}
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareActiveAuthMethodTypesStmt),
					prepareActiveAuthMethodTypesCols,
					[][]driver.Value{
						{
							true,
							domain.UserAuthMethodTypePasswordless,
							1,
						},
					},
				),
			},
			object: &AuthMethodTypes{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				AuthMethodTypes: []domain.UserAuthMethodType{
					domain.UserAuthMethodTypePasswordless,
					domain.UserAuthMethodTypePassword,
					domain.UserAuthMethodTypeIDP,
				},
			},
		},
		{
			name: "prepareUserAuthMethodTypesQuery multiple second factors",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
				builder, scan := prepareUserAuthMethodTypesQuery(ctx, db, true)
				return builder, func(rows *sql.Rows) (*AuthMethodTypes, error) {
					return scan(rows)
				}
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareActiveAuthMethodTypesStmt),
					prepareActiveAuthMethodTypesCols,
					[][]driver.Value{
						{
							true,
							domain.UserAuthMethodTypePasswordless,
							1,
						},
						{
							true,
							domain.UserAuthMethodTypeTOTP,
							1,
						},
					},
				),
			},
			object: &AuthMethodTypes{
				SearchResponse: SearchResponse{
					Count: 4,
				},
				AuthMethodTypes: []domain.UserAuthMethodType{
					domain.UserAuthMethodTypePasswordless,
					domain.UserAuthMethodTypeTOTP,
					domain.UserAuthMethodTypePassword,
					domain.UserAuthMethodTypeIDP,
				},
			},
		},
		{
			name: "prepareUserAuthMethodTypesQuery sql err",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*AuthMethodTypes, error)) {
				builder, scan := prepareUserAuthMethodTypesQuery(ctx, db, true)
				return builder, func(rows *sql.Rows) (*AuthMethodTypes, error) {
					return scan(rows)
				}
			},
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareActiveAuthMethodTypesStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*AuthMethodTypes)(nil),
		},
		{
			name: "prepareUserAuthMethodTypesRequiredQuery no result",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserAuthMethodRequirements, error)) {
				builder, scan := prepareUserAuthMethodTypesRequiredQuery(ctx, db)
				return builder, func(row *sql.Row) (*UserAuthMethodRequirements, error) {
					return scan(row)
				}
			},
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareAuthMethodTypesRequiredStmt),
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
			object: (*UserAuthMethodRequirements)(nil),
		},
		{
			name: "prepareUserAuthMethodTypesRequiredQuery one second factor",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserAuthMethodRequirements, error)) {
				builder, scan := prepareUserAuthMethodTypesRequiredQuery(ctx, db)
				return builder, func(row *sql.Row) (*UserAuthMethodRequirements, error) {
					return scan(row)
				}
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthMethodTypesRequiredStmt),
					prepareAuthMethodTypesRequiredCols,
					[][]driver.Value{
						{
							domain.UserTypeHuman,
							true,
							true,
						},
					},
				),
			},
			object: &UserAuthMethodRequirements{
				UserType:          domain.UserTypeHuman,
				ForceMFA:          true,
				ForceMFALocalOnly: true,
			},
		},
		{
			name: "prepareUserAuthMethodTypesRequiredQuery multiple second factors",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserAuthMethodRequirements, error)) {
				builder, scan := prepareUserAuthMethodTypesRequiredQuery(ctx, db)
				return builder, func(row *sql.Row) (*UserAuthMethodRequirements, error) {
					return scan(row)
				}
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthMethodTypesRequiredStmt),
					prepareAuthMethodTypesRequiredCols,
					[][]driver.Value{
						{
							domain.UserTypeHuman,
							true,
							true,
						},
					},
				),
			},

			object: &UserAuthMethodRequirements{
				UserType:          domain.UserTypeHuman,
				ForceMFA:          true,
				ForceMFALocalOnly: true,
			},
		},
		{
			name: "prepareUserAuthMethodTypesRequiredQuery sql err",
			prepare: func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*UserAuthMethodRequirements, error)) {
				builder, scan := prepareUserAuthMethodTypesRequiredQuery(ctx, db)
				return builder, func(row *sql.Row) (*UserAuthMethodRequirements, error) {
					return scan(row)
				}
			},
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareAuthMethodTypesRequiredStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
