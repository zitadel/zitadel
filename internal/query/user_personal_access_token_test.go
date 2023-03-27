package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	personalAccessTokenStmt = regexp.QuoteMeta(
		"SELECT projections.personal_access_tokens3.id," +
			" projections.personal_access_tokens3.creation_date," +
			" projections.personal_access_tokens3.change_date," +
			" projections.personal_access_tokens3.resource_owner," +
			" projections.personal_access_tokens3.sequence," +
			" projections.personal_access_tokens3.user_id," +
			" projections.personal_access_tokens3.expiration," +
			" projections.personal_access_tokens3.scopes" +
			" FROM projections.personal_access_tokens3" +
			` AS OF SYSTEM TIME '-1 ms'`)
	personalAccessTokenCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"expiration",
		"scopes",
	}
	personalAccessTokensStmt = regexp.QuoteMeta(
		"SELECT projections.personal_access_tokens3.id," +
			" projections.personal_access_tokens3.creation_date," +
			" projections.personal_access_tokens3.change_date," +
			" projections.personal_access_tokens3.resource_owner," +
			" projections.personal_access_tokens3.sequence," +
			" projections.personal_access_tokens3.user_id," +
			" projections.personal_access_tokens3.expiration," +
			" projections.personal_access_tokens3.scopes," +
			" COUNT(*) OVER ()" +
			" FROM projections.personal_access_tokens3" +
			" AS OF SYSTEM TIME '-1 ms'")
	personalAccessTokensCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"expiration",
		"scopes",
		"count",
	}
)

func Test_PersonalAccessTokenPrepares(t *testing.T) {
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
			name:    "preparePersonalAccessTokenQuery no result",
			prepare: preparePersonalAccessTokenQuery,
			want: want{
				sqlExpectations: mockQuery(
					personalAccessTokenStmt,
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
			object: (*PersonalAccessToken)(nil),
		},
		{
			name:    "preparePersonalAccessTokenQuery found",
			prepare: preparePersonalAccessTokenQuery,
			want: want{
				sqlExpectations: mockQuery(
					personalAccessTokenStmt,
					personalAccessTokenCols,
					[]driver.Value{
						"token-id",
						testNow,
						testNow,
						"ro",
						uint64(20211202),
						"user-id",
						time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						database.StringArray{"openid"},
					},
				),
			},
			object: &PersonalAccessToken{
				ID:            "token-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211202,
				UserID:        "user-id",
				Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
				Scopes:        database.StringArray{"openid"},
			},
		},
		{
			name:    "preparePersonalAccessTokenQuery sql err",
			prepare: preparePersonalAccessTokenQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					personalAccessTokenStmt,
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
		{
			name:    "preparePersonalAccessTokensQuery no result",
			prepare: preparePersonalAccessTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					personalAccessTokensStmt,
					nil,
					nil,
				),
			},
			object: &PersonalAccessTokens{PersonalAccessTokens: []*PersonalAccessToken{}},
		},
		{
			name:    "preparePersonalAccessTokensQuery one token",
			prepare: preparePersonalAccessTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					personalAccessTokensStmt,
					personalAccessTokensCols,
					[][]driver.Value{
						{
							"token-id",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							database.StringArray{"openid"},
						},
					},
				),
			},
			object: &PersonalAccessTokens{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				PersonalAccessTokens: []*PersonalAccessToken{
					{
						ID:            "token-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        database.StringArray{"openid"},
					},
				},
			},
		},
		{
			name:    "preparePersonalAccessTokensQuery multiple tokens",
			prepare: preparePersonalAccessTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					personalAccessTokensStmt,
					personalAccessTokensCols,
					[][]driver.Value{
						{
							"token-id",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							database.StringArray{"openid"},
						},
						{
							"token-id2",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							database.StringArray{"openid"},
						},
					},
				),
			},
			object: &PersonalAccessTokens{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				PersonalAccessTokens: []*PersonalAccessToken{
					{
						ID:            "token-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        database.StringArray{"openid"},
					},
					{
						ID:            "token-id2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        database.StringArray{"openid"},
					},
				},
			},
		},
		{
			name:    "preparePersonalAccessTokensQuery sql err",
			prepare: preparePersonalAccessTokensQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					personalAccessTokensStmt,
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
