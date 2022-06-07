package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
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
					regexp.QuoteMeta(`SELECT zitadel.projections.user_auth_methods.token_id,`+
						` zitadel.projections.user_auth_methods.creation_date,`+
						` zitadel.projections.user_auth_methods.change_date,`+
						` zitadel.projections.user_auth_methods.resource_owner,`+
						` zitadel.projections.user_auth_methods.user_id,`+
						` zitadel.projections.user_auth_methods.sequence,`+
						` zitadel.projections.user_auth_methods.name,`+
						` zitadel.projections.user_auth_methods.state,`+
						` zitadel.projections.user_auth_methods.method_type,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.user_auth_methods`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.user_auth_methods.token_id,`+
						` zitadel.projections.user_auth_methods.creation_date,`+
						` zitadel.projections.user_auth_methods.change_date,`+
						` zitadel.projections.user_auth_methods.resource_owner,`+
						` zitadel.projections.user_auth_methods.user_id,`+
						` zitadel.projections.user_auth_methods.sequence,`+
						` zitadel.projections.user_auth_methods.name,`+
						` zitadel.projections.user_auth_methods.state,`+
						` zitadel.projections.user_auth_methods.method_type,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.user_auth_methods`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.user_auth_methods.token_id,`+
						` zitadel.projections.user_auth_methods.creation_date,`+
						` zitadel.projections.user_auth_methods.change_date,`+
						` zitadel.projections.user_auth_methods.resource_owner,`+
						` zitadel.projections.user_auth_methods.user_id,`+
						` zitadel.projections.user_auth_methods.sequence,`+
						` zitadel.projections.user_auth_methods.name,`+
						` zitadel.projections.user_auth_methods.state,`+
						` zitadel.projections.user_auth_methods.method_type,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.user_auth_methods`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.user_auth_methods.token_id,`+
						` zitadel.projections.user_auth_methods.creation_date,`+
						` zitadel.projections.user_auth_methods.change_date,`+
						` zitadel.projections.user_auth_methods.resource_owner,`+
						` zitadel.projections.user_auth_methods.user_id,`+
						` zitadel.projections.user_auth_methods.sequence,`+
						` zitadel.projections.user_auth_methods.name,`+
						` zitadel.projections.user_auth_methods.state,`+
						` zitadel.projections.user_auth_methods.method_type,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.user_auth_methods`),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
