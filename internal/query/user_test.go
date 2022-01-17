package query

import (
	"fmt"
	"regexp"
	"testing"

	errs "github.com/caos/zitadel/internal/errors"
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
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.users.id,`+
						` zitadel.projections.users.creation_date,`+
						` zitadel.projections.users.change_date,`+
						` zitadel.projections.users.resource_owner,`+
						` zitadel.projections.users.sequence,`+
						` zitadel.projections.users.state,`+
						` zitadel.projections.users.username,`+
						//` zitadel.projections.users.type,`+
						` login_names.login_names,`+
						` preferred_login_name.login_name,`+
						` zitadel.projections.users_humans.user_id,`+
						` zitadel.projections.users_humans.first_name,`+
						` zitadel.projections.users_humans.last_name,`+
						` zitadel.projections.users_humans.nick_name,`+
						` zitadel.projections.users_humans.display_name,`+
						` zitadel.projections.users_humans.preferred_language,`+
						` zitadel.projections.users_humans.gender,`+
						` zitadel.projections.users_humans.avater_key,`+
						` zitadel.projections.users_humans.email,`+
						` zitadel.projections.users_humans.is_email_verified,`+
						` zitadel.projections.users_humans.phone,`+
						` zitadel.projections.users_humans.is_phone_verified,`+
						` zitadel.projections.users_machines.user_id,`+
						` zitadel.projections.users_machines.name,`+
						` zitadel.projections.users_machines.description`+
						` FROM zitadel.projections.users`+
						` LEFT JOIN zitadel.projections.users_humans ON zitadel.projections.users.id = zitadel.projections.users_humans.user_id`+
						` LEFT JOIN zitadel.projections.users_machines ON zitadel.projections.users.id = zitadel.projections.users_machines.user_id`+
						` LEFT JOIN`+
						` (SELECT login_names.user_id, ARRAY_AGG(login_names.login_name) as login_names`+
						` FROM zitadel.projections.login_names as login_names`+
						` GROUP BY login_names.user_id) as login_names`+
						` on login_names.user_id = zitadel.projections.users.id`+
						` LEFT JOIN`+
						` (SELECT preferred_login_name.user_id, preferred_login_name.login_name FROM zitadel.projections.login_names as preferred_login_name WHERE preferred_login_name.is_primary = $1) as preferred_login_name`+
						` on preferred_login_name.user_id = zitadel.projections.users.id`),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*MessageText)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
