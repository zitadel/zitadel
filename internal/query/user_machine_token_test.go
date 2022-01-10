package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/lib/pq"

	errs "github.com/caos/zitadel/internal/errors"
)

var (
	machineTokenStmt = regexp.QuoteMeta(
		"SELECT zitadel.projections.machine_tokens.id," +
			" zitadel.projections.machine_tokens.creation_date," +
			" zitadel.projections.machine_tokens.change_date," +
			" zitadel.projections.machine_tokens.resource_owner," +
			" zitadel.projections.machine_tokens.sequence," +
			" zitadel.projections.machine_tokens.user_id," +
			" zitadel.projections.machine_tokens.expiration," +
			" zitadel.projections.machine_tokens.scopes" +
			" FROM zitadel.projections.machine_tokens")
	machineTokenCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"expiration",
		"scopes",
	}
	machineTokensStmt = regexp.QuoteMeta(
		"SELECT zitadel.projections.machine_tokens.id," +
			" zitadel.projections.machine_tokens.creation_date," +
			" zitadel.projections.machine_tokens.change_date," +
			" zitadel.projections.machine_tokens.resource_owner," +
			" zitadel.projections.machine_tokens.sequence," +
			" zitadel.projections.machine_tokens.user_id," +
			" zitadel.projections.machine_tokens.expiration," +
			" zitadel.projections.machine_tokens.scopes," +
			" COUNT(*) OVER ()" +
			" FROM zitadel.projections.machine_tokens")
	machineTokensCols = []string{
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

func Test_MachineTokenPrepares(t *testing.T) {
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
			name:    "prepareMachineTokenQuery no result",
			prepare: prepareMachineTokenQuery,
			want: want{
				sqlExpectations: mockQuery(
					machineTokenStmt,
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
			object: (*MachineToken)(nil),
		},
		{
			name:    "prepareMachineTokenQuery found",
			prepare: prepareMachineTokenQuery,
			want: want{
				sqlExpectations: mockQuery(
					machineTokenStmt,
					machineTokenCols,
					[]driver.Value{
						"token-id",
						testNow,
						testNow,
						"ro",
						uint64(20211202),
						"user-id",
						time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						pq.StringArray{"openid"},
					},
				),
			},
			object: &MachineToken{
				ID:            "token-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211202,
				UserID:        "user-id",
				Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
				Scopes:        []string{"openid"},
			},
		},
		{
			name:    "prepareMachineTokenQuery sql err",
			prepare: prepareMachineTokenQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					machineTokenStmt,
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
			name:    "prepareMachineTokensQuery no result",
			prepare: prepareMachineTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					machineTokensStmt,
					nil,
					nil,
				),
			},
			object: &MachineTokens{MachineTokens: []*MachineToken{}},
		},
		{
			name:    "prepareMachineTokensQuery one token",
			prepare: prepareMachineTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					machineTokensStmt,
					machineTokensCols,
					[][]driver.Value{
						{
							"token-id",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							pq.StringArray{"openid"},
						},
					},
				),
			},
			object: &MachineTokens{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				MachineTokens: []*MachineToken{
					{
						ID:            "token-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        []string{"openid"},
					},
				},
			},
		},
		{
			name:    "prepareMachineTokensQuery multiple tokens",
			prepare: prepareMachineTokensQuery,
			want: want{
				sqlExpectations: mockQueries(
					machineTokensStmt,
					machineTokensCols,
					[][]driver.Value{
						{
							"token-id",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							pq.StringArray{"openid"},
						},
						{
							"token-id2",
							testNow,
							testNow,
							"ro",
							uint64(20211202),
							"user-id",
							time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
							pq.StringArray{"openid"},
						},
					},
				),
			},
			object: &MachineTokens{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				MachineTokens: []*MachineToken{
					{
						ID:            "token-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        []string{"openid"},
					},
					{
						ID:            "token-id2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211202,
						UserID:        "user-id",
						Expiration:    time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
						Scopes:        []string{"openid"},
					},
				},
			},
		},
		{
			name:    "prepareMachineTokensQuery sql err",
			prepare: prepareMachineTokensQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					machineTokensStmt,
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
