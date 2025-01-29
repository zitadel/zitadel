package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUser_usersCheckPermission(t *testing.T) {
	type want struct {
		users []*User
	}
	tests := []struct {
		name        string
		want        want
		users       *Users
		permissions []string
	}{
		{
			"permissions for all users",
			want{
				users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first", "second", "third"},
		},
		{
			"permissions for one user, first",
			want{
				users: []*User{
					{ID: "first"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first"},
		},
		{
			"permissions for one user, second",
			want{
				users: []*User{
					{ID: "second"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"second"},
		},
		{
			"permissions for one user, third",
			want{
				users: []*User{
					{ID: "third"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"third"},
		},
		{
			"permissions for two users, first",
			want{
				users: []*User{
					{ID: "first"}, {ID: "third"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"first", "third"},
		},
		{
			"permissions for two users, second",
			want{
				users: []*User{
					{ID: "second"}, {ID: "third"},
				},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
			},
			[]string{"second", "third"},
		},
		{
			"no permissions",
			want{
				users: []*User{},
			},
			&Users{
				Users: []*User{
					{ID: "first"}, {ID: "second"}, {ID: "third"},
				},
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
			usersCheckPermission(context.Background(), tt.users, checkPermission)
			require.Equal(t, tt.want.users, tt.users.Users)
		})
	}
}

func TestUser_userCheckPermission(t *testing.T) {
	type args struct {
		ctxData       string
		resourceowner string
		user          string
	}
	type perm struct {
		resourceowner string
		user          string
	}
	tests := []struct {
		name        string
		wantErr     bool
		args        args
		permissions []perm
	}{
		{
			name: "permission, self",
			args: args{
				resourceowner: "org",
				user:          "user",
				ctxData:       "user",
			},
			permissions: []perm{},
		},
		{
			name: "permission, user",
			args: args{
				resourceowner: "org1",
				user:          "user1",
				ctxData:       "user2",
			},
			permissions: []perm{{"org1", "user1"}},
			wantErr:     false,
		},
		{
			name: "permission, org",
			args: args{
				resourceowner: "org1",
				user:          "user1",
				ctxData:       "user2",
			},
			permissions: []perm{{"org1", "user3"}},
		},
		{
			name: "permission, none",
			args: args{
				resourceowner: "org1",
				user:          "user1",
				ctxData:       "user2",
			},
			permissions: []perm{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				for _, perm := range tt.permissions {
					if resourceID == perm.user {
						return nil
					}
					if orgID == perm.resourceowner {
						return nil
					}
				}
				return errors.New("failed")
			}

			granted := userCheckPermission(authz.SetCtxData(context.Background(), authz.CtxData{UserID: tt.args.ctxData}), tt.args.resourceowner, tt.args.user, checkPermission)
			if tt.wantErr {
				assert.Error(t, granted)
			} else {
				assert.NoError(t, granted)
			}
		})
	}
}

var (
	loginNamesQuery = `SELECT login_names.user_id, ARRAY_AGG(login_names.login_name)::TEXT[] AS loginnames, ARRAY_AGG(LOWER(login_names.login_name))::TEXT[] AS loginnames_lower, login_names.instance_id` +
		` FROM projections.login_names3 AS login_names` +
		` GROUP BY login_names.user_id, login_names.instance_id`
	preferredLoginNameQuery = `SELECT preferred_login_name.user_id, preferred_login_name.login_name, preferred_login_name.instance_id` +
		` FROM projections.login_names3 AS preferred_login_name` +
		` WHERE  preferred_login_name.is_primary = $1`
	userQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14.state,` +
		` projections.users14.type,` +
		` projections.users14.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.first_name,` +
		` projections.users14_humans.last_name,` +
		` projections.users14_humans.nick_name,` +
		` projections.users14_humans.display_name,` +
		` projections.users14_humans.preferred_language,` +
		` projections.users14_humans.gender,` +
		` projections.users14_humans.avatar_key,` +
		` projections.users14_humans.email,` +
		` projections.users14_humans.is_email_verified,` +
		` projections.users14_humans.phone,` +
		` projections.users14_humans.is_phone_verified,` +
		` projections.users14_humans.password_change_required,` +
		` projections.users14_humans.password_changed,` +
		` projections.users14_humans.mfa_init_skipped,` +
		` projections.users14_machines.user_id,` +
		` projections.users14_machines.name,` +
		` projections.users14_machines.description,` +
		` projections.users14_machines.secret,` +
		` projections.users14_machines.access_token_type,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` LEFT JOIN projections.users14_machines ON projections.users14.id = projections.users14_machines.user_id AND projections.users14.instance_id = projections.users14_machines.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users14.id AND login_names.instance_id = projections.users14.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users14.id AND preferred_login_name.instance_id = projections.users14.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	userCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		// human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		"email",
		"is_email_verified",
		"phone",
		"is_phone_verified",
		"password_change_required",
		"password_changed",
		"mfa_init_skipped",
		// machine
		"user_id",
		"name",
		"description",
		"secret",
		"access_token_type",
		"count",
	}
	profileQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.first_name,` +
		` projections.users14_humans.last_name,` +
		` projections.users14_humans.nick_name,` +
		` projections.users14_humans.display_name,` +
		` projections.users14_humans.preferred_language,` +
		` projections.users14_humans.gender,` +
		` projections.users14_humans.avatar_key` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	profileCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
	}
	emailQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.email,` +
		` projections.users14_humans.is_email_verified` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	emailCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"email",
		"is_email_verified",
	}
	phoneQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.phone,` +
		` projections.users14_humans.is_phone_verified` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	phoneCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"phone",
		"is_phone_verified",
	}
	userUniqueQuery = `SELECT projections.users14.id,` +
		` projections.users14.state,` +
		` projections.users14.username,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.email,` +
		` projections.users14_humans.is_email_verified` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	userUniqueCols = []string{
		"id",
		"state",
		"username",
		"user_id",
		"email",
		"is_email_verified",
	}
	notifyUserQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14.state,` +
		` projections.users14.type,` +
		` projections.users14.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.first_name,` +
		` projections.users14_humans.last_name,` +
		` projections.users14_humans.nick_name,` +
		` projections.users14_humans.display_name,` +
		` projections.users14_humans.preferred_language,` +
		` projections.users14_humans.gender,` +
		` projections.users14_humans.avatar_key,` +
		` projections.users14_notifications.user_id,` +
		` projections.users14_notifications.last_email,` +
		` projections.users14_notifications.verified_email,` +
		` projections.users14_notifications.last_phone,` +
		` projections.users14_notifications.verified_phone,` +
		` projections.users14_notifications.password_set,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` LEFT JOIN projections.users14_notifications ON projections.users14.id = projections.users14_notifications.user_id AND projections.users14.instance_id = projections.users14_notifications.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users14.id AND login_names.instance_id = projections.users14.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users14.id AND preferred_login_name.instance_id = projections.users14.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	notifyUserCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		// human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		// machine
		"user_id",
		"last_email",
		"verified_email",
		"last_phone",
		"verified_phone",
		"password_set",
		"count",
	}
	usersQuery = `SELECT projections.users14.id,` +
		` projections.users14.creation_date,` +
		` projections.users14.change_date,` +
		` projections.users14.resource_owner,` +
		` projections.users14.sequence,` +
		` projections.users14.state,` +
		` projections.users14.type,` +
		` projections.users14.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users14_humans.user_id,` +
		` projections.users14_humans.first_name,` +
		` projections.users14_humans.last_name,` +
		` projections.users14_humans.nick_name,` +
		` projections.users14_humans.display_name,` +
		` projections.users14_humans.preferred_language,` +
		` projections.users14_humans.gender,` +
		` projections.users14_humans.avatar_key,` +
		` projections.users14_humans.email,` +
		` projections.users14_humans.is_email_verified,` +
		` projections.users14_humans.phone,` +
		` projections.users14_humans.is_phone_verified,` +
		` projections.users14_humans.password_change_required,` +
		` projections.users14_humans.password_changed,` +
		` projections.users14_machines.user_id,` +
		` projections.users14_machines.name,` +
		` projections.users14_machines.description,` +
		` projections.users14_machines.secret,` +
		` projections.users14_machines.access_token_type,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users14` +
		` LEFT JOIN projections.users14_humans ON projections.users14.id = projections.users14_humans.user_id AND projections.users14.instance_id = projections.users14_humans.instance_id` +
		` LEFT JOIN projections.users14_machines ON projections.users14.id = projections.users14_machines.user_id AND projections.users14.instance_id = projections.users14_machines.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users14.id AND login_names.instance_id = projections.users14.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users14.id AND preferred_login_name.instance_id = projections.users14.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	usersCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		// human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		"email",
		"is_email_verified",
		"phone",
		"is_phone_verified",
		"password_change_required",
		"password_changed",
		// machine
		"user_id",
		"name",
		"description",
		"secret",
		"access_token_type",
		"count",
	}
	countUsersQuery = "SELECT COUNT(*) OVER () FROM projections.users14"
	countUsersCols  = []string{"count"}
)

func Test_UserPrepares(t *testing.T) {
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
			name:    "prepareUserQuery no result",
			prepare: prepareUserQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(userQuery),
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
			object: (*User)(nil),
		},
		{
			name:    "prepareUserQuery human found",
			prepare: prepareUserQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userQuery),
					userCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						domain.UserStateActive,
						domain.UserTypeHuman,
						"username",
						database.TextArray[string]{"login_name1", "login_name2"},
						"login_name1",
						// human
						"id",
						"first_name",
						"last_name",
						"nick_name",
						"display_name",
						"de",
						domain.GenderUnspecified,
						"avatar_key",
						"email",
						true,
						"phone",
						true,
						true,
						testNow,
						testNow,
						// machine
						nil,
						nil,
						nil,
						nil,
						nil,
						1,
					},
				),
			},
			object: &User{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				ResourceOwner:      "resource_owner",
				Sequence:           20211108,
				State:              domain.UserStateActive,
				Type:               domain.UserTypeHuman,
				Username:           "username",
				LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
				PreferredLoginName: "login_name1",
				Human: &Human{
					FirstName:              "first_name",
					LastName:               "last_name",
					NickName:               "nick_name",
					DisplayName:            "display_name",
					AvatarKey:              "avatar_key",
					PreferredLanguage:      language.German,
					Gender:                 domain.GenderUnspecified,
					Email:                  "email",
					IsEmailVerified:        true,
					Phone:                  "phone",
					IsPhoneVerified:        true,
					PasswordChangeRequired: true,
					PasswordChanged:        testNow,
					MFAInitSkipped:         testNow,
				},
			},
		},
		{
			name:    "prepareUserQuery machine found",
			prepare: prepareUserQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userQuery),
					userCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						domain.UserStateActive,
						domain.UserTypeMachine,
						"username",
						database.TextArray[string]{"login_name1", "login_name2"},
						"login_name1",
						// human
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						// machine
						"id",
						"name",
						"description",
						nil,
						domain.OIDCTokenTypeBearer,
						1,
					},
				),
			},
			object: &User{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				ResourceOwner:      "resource_owner",
				Sequence:           20211108,
				State:              domain.UserStateActive,
				Type:               domain.UserTypeMachine,
				Username:           "username",
				LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
				PreferredLoginName: "login_name1",
				Machine: &Machine{
					Name:            "name",
					Description:     "description",
					EncodedSecret:   "",
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
			},
		},
		{
			name:    "prepareUserQuery machine with secret found",
			prepare: prepareUserQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userQuery),
					userCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						domain.UserStateActive,
						domain.UserTypeMachine,
						"username",
						database.TextArray[string]{"login_name1", "login_name2"},
						"login_name1",
						// human
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						// machine
						"id",
						"name",
						"description",
						"secret",
						domain.OIDCTokenTypeBearer,
						1,
					},
				),
			},
			object: &User{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				ResourceOwner:      "resource_owner",
				Sequence:           20211108,
				State:              domain.UserStateActive,
				Type:               domain.UserTypeMachine,
				Username:           "username",
				LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
				PreferredLoginName: "login_name1",
				Machine: &Machine{
					Name:            "name",
					Description:     "description",
					EncodedSecret:   "secret",
					AccessTokenType: domain.OIDCTokenTypeBearer,
				},
			},
		},
		{
			name:    "prepareUserQuery sql err",
			prepare: prepareUserQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(userQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*User)(nil),
		},
		{
			name:    "prepareProfileQuery no result",
			prepare: prepareProfileQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(profileQuery),
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
			object: (*Profile)(nil),
		},
		{
			name:    "prepareProfileQuery human found",
			prepare: prepareProfileQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(profileQuery),
					profileCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						"id",
						"first_name",
						"last_name",
						"nick_name",
						"display_name",
						"de",
						domain.GenderUnspecified,
						"avatar_key",
					},
				),
			},
			object: &Profile{
				ID:                "id",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "resource_owner",
				Sequence:          20211108,
				FirstName:         "first_name",
				LastName:          "last_name",
				NickName:          "nick_name",
				DisplayName:       "display_name",
				AvatarKey:         "avatar_key",
				PreferredLanguage: language.German,
				Gender:            domain.GenderUnspecified,
			},
		},
		{
			name:    "prepareProfileQuery not human found (error)",
			prepare: prepareProfileQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(profileQuery),
					profileCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsPreconditionFailed(err) {
						return fmt.Errorf("err should be zitadel.PredconditionError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Profile)(nil),
		},
		{
			name:    "prepareProfileQuery sql err",
			prepare: prepareProfileQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(profileQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Profile)(nil),
		},
		{
			name:    "prepareEmailQuery no result",
			prepare: prepareEmailQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(emailQuery),
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
			object: (*Email)(nil),
		},
		{
			name:    "prepareEmailQuery human found",
			prepare: prepareEmailQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(emailQuery),
					emailCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						//domain.UserStateActive,
						"id",
						"email",
						true,
					},
				),
			},
			object: &Email{
				ID:            "id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				//State:              domain.UserStateActive,
				Email:      "email",
				IsVerified: true,
			},
		},
		{
			name:    "prepareEmailQuery not human found (error)",
			prepare: prepareEmailQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(emailQuery),
					emailCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						//domain.UserStateActive,
						nil,
						nil,
						nil,
					},
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsPreconditionFailed(err) {
						return fmt.Errorf("err should be zitadel.PredconditionError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Email)(nil),
		},
		{
			name:    "prepareEmailQuery sql err",
			prepare: prepareEmailQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(emailQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Email)(nil),
		},
		{
			name:    "preparePhoneQuery no result",
			prepare: preparePhoneQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(phoneQuery),
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
			object: (*Phone)(nil),
		},
		{
			name:    "preparePhoneQuery human found",
			prepare: preparePhoneQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(phoneQuery),
					phoneCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						//domain.UserStateActive,
						"id",
						"phone",
						true,
					},
				),
			},
			object: &Phone{
				ID:            "id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				//State:              domain.UserStateActive,
				Phone:      "phone",
				IsVerified: true,
			},
		},
		{
			name:    "preparePhoneQuery not human found (error)",
			prepare: preparePhoneQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(phoneQuery),
					phoneCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						//domain.UserStateActive,
						nil,
						nil,
						nil,
					},
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsPreconditionFailed(err) {
						return fmt.Errorf("err should be zitadel.PredconditionError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Phone)(nil),
		},
		{
			name:    "preparePhoneQuery sql err",
			prepare: preparePhoneQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(phoneQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Phone)(nil),
		},
		{
			name:    "prepareUserUniqueQuery no result",
			prepare: prepareUserUniqueQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userUniqueQuery),
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
			object: true,
		},
		{
			name:    "prepareUserUniqueQuery found",
			prepare: prepareUserUniqueQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userUniqueQuery),
					userUniqueCols,
					[]driver.Value{
						"id",
						domain.UserStateActive,
						"username",
						"id",
						"email",
						true,
					},
				),
			},
			object: false,
		},
		{
			name:    "prepareUserUniqueQuery sql err",
			prepare: prepareUserUniqueQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(userUniqueQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: false,
		},
		{
			name:    "prepareNotifyUserQuery no result",
			prepare: prepareNotifyUserQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(notifyUserQuery),
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
			object: (*NotifyUser)(nil),
		},
		{
			name:    "prepareNotifyUserQuery notify found",
			prepare: prepareNotifyUserQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(notifyUserQuery),
					notifyUserCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						domain.UserStateActive,
						domain.UserTypeHuman,
						"username",
						database.TextArray[string]{"login_name1", "login_name2"},
						"login_name1",
						// human
						"id",
						"first_name",
						"last_name",
						"nick_name",
						"display_name",
						"de",
						domain.GenderUnspecified,
						"avatar_key",
						//notify
						"id",
						"lastEmail",
						"verifiedEmail",
						"lastPhone",
						"verifiedPhone",
						true,
						1,
					},
				),
			},
			object: &NotifyUser{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				ResourceOwner:      "resource_owner",
				Sequence:           20211108,
				State:              domain.UserStateActive,
				Type:               domain.UserTypeHuman,
				Username:           "username",
				LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
				PreferredLoginName: "login_name1",
				FirstName:          "first_name",
				LastName:           "last_name",
				NickName:           "nick_name",
				DisplayName:        "display_name",
				AvatarKey:          "avatar_key",
				PreferredLanguage:  language.German,
				Gender:             domain.GenderUnspecified,
				LastEmail:          "lastEmail",
				VerifiedEmail:      "verifiedEmail",
				LastPhone:          "lastPhone",
				VerifiedPhone:      "verifiedPhone",
				PasswordSet:        true,
			},
		},
		{
			name:    "prepareNotifyUserQuery not notify found (error)",
			prepare: prepareNotifyUserQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(notifyUserQuery),
					notifyUserCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						domain.UserStateActive,
						domain.UserTypeHuman,
						"username",
						database.TextArray[string]{"login_name1", "login_name2"},
						"login_name1",
						// human
						"id",
						"first_name",
						"last_name",
						"nick_name",
						"display_name",
						"de",
						domain.GenderUnspecified,
						"avatar_key",
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						1,
					},
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsPreconditionFailed(err) {
						return fmt.Errorf("err should be zitadel.PredconditionError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*NotifyUser)(nil),
		},
		{
			name:    "prepareNotifyUserQuery sql err",
			prepare: prepareNotifyUserQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(notifyUserQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*NotifyUser)(nil),
		},
		{
			name:    "prepareUsersQuery no result",
			prepare: prepareUsersQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(usersQuery),
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
			object: &Users{Users: []*User{}},
		},
		{
			name:    "prepareUsersQuery one result",
			prepare: prepareUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(usersQuery),
					usersCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							domain.UserStateActive,
							domain.UserTypeHuman,
							"username",
							database.TextArray[string]{"login_name1", "login_name2"},
							"login_name1",
							// human
							"id",
							"first_name",
							"last_name",
							"nick_name",
							"display_name",
							"de",
							domain.GenderUnspecified,
							"avatar_key",
							"email",
							true,
							"phone",
							true,
							true,
							testNow,
							// machine
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &Users{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Users: []*User{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						ResourceOwner:      "resource_owner",
						Sequence:           20211108,
						State:              domain.UserStateActive,
						Type:               domain.UserTypeHuman,
						Username:           "username",
						LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
						PreferredLoginName: "login_name1",
						Human: &Human{
							FirstName:              "first_name",
							LastName:               "last_name",
							NickName:               "nick_name",
							DisplayName:            "display_name",
							AvatarKey:              "avatar_key",
							PreferredLanguage:      language.German,
							Gender:                 domain.GenderUnspecified,
							Email:                  "email",
							IsEmailVerified:        true,
							Phone:                  "phone",
							IsPhoneVerified:        true,
							PasswordChangeRequired: true,
							PasswordChanged:        testNow,
						},
					},
				},
			},
		},
		{
			name:    "prepareUsersQuery multiple results",
			prepare: prepareUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(usersQuery),
					usersCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							domain.UserStateActive,
							domain.UserTypeHuman,
							"username",
							database.TextArray[string]{"login_name1", "login_name2"},
							"login_name1",
							// human
							"id",
							"first_name",
							"last_name",
							"nick_name",
							"display_name",
							"de",
							domain.GenderUnspecified,
							"avatar_key",
							"email",
							true,
							"phone",
							true,
							true,
							testNow,
							// machine
							nil,
							nil,
							nil,
							nil,
							nil,
						},
						{
							"id",
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							domain.UserStateActive,
							domain.UserTypeMachine,
							"username",
							database.TextArray[string]{"login_name1", "login_name2"},
							"login_name1",
							// human
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							// machine
							"id",
							"name",
							"description",
							"secret",
							domain.OIDCTokenTypeBearer,
						},
					},
				),
			},
			object: &Users{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Users: []*User{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						ResourceOwner:      "resource_owner",
						Sequence:           20211108,
						State:              domain.UserStateActive,
						Type:               domain.UserTypeHuman,
						Username:           "username",
						LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
						PreferredLoginName: "login_name1",
						Human: &Human{
							FirstName:              "first_name",
							LastName:               "last_name",
							NickName:               "nick_name",
							DisplayName:            "display_name",
							AvatarKey:              "avatar_key",
							PreferredLanguage:      language.German,
							Gender:                 domain.GenderUnspecified,
							Email:                  "email",
							IsEmailVerified:        true,
							Phone:                  "phone",
							IsPhoneVerified:        true,
							PasswordChangeRequired: true,
							PasswordChanged:        testNow,
						},
					},
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						ResourceOwner:      "resource_owner",
						Sequence:           20211108,
						State:              domain.UserStateActive,
						Type:               domain.UserTypeMachine,
						Username:           "username",
						LoginNames:         database.TextArray[string]{"login_name1", "login_name2"},
						PreferredLoginName: "login_name1",
						Machine: &Machine{
							Name:            "name",
							Description:     "description",
							EncodedSecret:   "secret",
							AccessTokenType: domain.OIDCTokenTypeBearer,
						},
					},
				},
			},
		},
		{
			name:    "prepareUsersQuery sql err",
			prepare: prepareUsersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(usersQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Users)(nil),
		},
		{
			name:    "prepareCountUsersQuery no result",
			prepare: prepareCountUsersQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(countUsersQuery),
					nil,
					nil,
				),
			},
			object: uint64(0),
		},
		{
			name:    "prepareCountUsersQuery one result",
			prepare: prepareCountUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(countUsersQuery),
					countUsersCols,
					[][]driver.Value{{uint64(1)}},
				),
			},
			object: uint64(1),
		},
		{
			name:    "prepareCountUsersQuery multiple results",
			prepare: prepareCountUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(countUsersQuery),
					countUsersCols,
					[][]driver.Value{{uint64(2)}},
				),
			},
			object: uint64(2),
		},
		{
			name:    "prepareCountUsersQuery sql err",
			prepare: prepareCountUsersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(countUsersQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := defaultPrepareArgs
			if reflect.TypeOf(tt.prepare).NumIn() == 0 {
				params = []reflect.Value{}
			}

			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, params...)
		})
	}
}
