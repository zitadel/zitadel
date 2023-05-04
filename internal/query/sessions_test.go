package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	expectedSessionQuery = regexp.QuoteMeta(`SELECT projections.sessions.id,` +
		` projections.sessions.creation_date,` +
		` projections.sessions.change_date,` +
		` projections.sessions.sequence,` +
		` projections.sessions.state,` +
		` projections.sessions.resource_owner,` +
		` projections.sessions.creator,` +
		` projections.sessions.user_id,` +
		` projections.sessions.user_checked_at,` +
		` projections.login_names2.login_name,` +
		` projections.users8_humans.display_name,` +
		` projections.sessions.password_checked_at,` +
		` projections.sessions.metadata` +
		` FROM projections.sessions` +
		` LEFT JOIN projections.login_names2 ON projections.sessions.user_id = projections.login_names2.user_id AND projections.sessions.instance_id = projections.login_names2.instance_id` +
		` LEFT JOIN projections.users8_humans ON projections.sessions.user_id = projections.users8_humans.user_id AND projections.sessions.instance_id = projections.users8_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedSessionsQuery = regexp.QuoteMeta(`SELECT projections.sessions.id,` +
		` projections.sessions.creation_date,` +
		` projections.sessions.change_date,` +
		` projections.sessions.sequence,` +
		` projections.sessions.state,` +
		` projections.sessions.resource_owner,` +
		` projections.sessions.creator,` +
		` projections.sessions.user_id,` +
		` projections.sessions.user_checked_at,` +
		` projections.login_names2.login_name,` +
		` projections.users8_humans.display_name,` +
		` projections.sessions.password_checked_at,` +
		` projections.sessions.metadata,` +
		` COUNT(*) OVER ()` +
		` FROM projections.sessions` +
		` LEFT JOIN projections.login_names2 ON projections.sessions.user_id = projections.login_names2.user_id AND projections.sessions.instance_id = projections.login_names2.instance_id` +
		` LEFT JOIN projections.users8_humans ON projections.sessions.user_id = projections.users8_humans.user_id AND projections.sessions.instance_id = projections.users8_humans.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)

	sessionCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"resource_owner",
		"sessions.creator",
		"sessions.user_id",
		"sessions.user_checked_at",
		"login_names2.login_name",
		"users8_humans.display_name",
		"sessions.password_checked_at",
		"sessions.metadata",
	}
	sessionsCols = append(sessionCols, "count")
)

func Test_SessionsPrepare(t *testing.T) {
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
			name:    "prepareSessionsQuery no result",
			prepare: prepareSessionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSessionsQuery,
					nil,
					nil,
				),
			},
			object: &Sessions{Sessions: []*Session{}},
		},
		{
			name:    "prepareSessionQuery",
			prepare: prepareSessionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSessionsQuery,
					sessionsCols,
					[][]driver.Value{
						{
							"session-id",
							testNow,
							testNow,
							uint64(20211109),
							domain.SessionStateActive,
							"ro",
							"creator",
							"user-id",
							testNow,
							"login-name",
							"display-name",
							testNow,
							[]byte(`{"key": "dmFsdWU="}`),
							//&crypto.CryptoValue{},
						},
					},
				),
			},
			object: &Sessions{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Sessions: []*Session{
					{
						ID:            "session-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						State:         domain.SessionStateActive,
						ResourceOwner: "ro",
						Creator:       "creator",
						UserFactor: SessionUserFactor{
							UserID:        "user-id",
							UserCheckedAt: testNow,
							LoginName:     "login-name",
							DisplayName:   "display-name",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						Metadata: map[string][]byte{
							"key": []byte("value"),
						},
					},
				},
			},
		},
		{
			name:    "prepareSessionsQuery multiple result",
			prepare: prepareSessionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSessionsQuery,
					sessionsCols,
					[][]driver.Value{
						{
							"session-id",
							testNow,
							testNow,
							uint64(20211109),
							domain.SessionStateActive,
							"ro",
							"creator",
							"user-id",
							testNow,
							"login-name",
							"display-name",
							testNow,
							nil,
							//map[string][]byte{
							//	"key": []byte("value"),
							//},
							//&crypto.CryptoValue{},
						},
						{
							"session-id2",
							testNow,
							testNow,
							uint64(20211109),
							domain.SessionStateActive,
							"ro",
							"creator2",
							"user-id2",
							testNow,
							"login-name2",
							"display-name2",
							testNow,
							nil,
							//map[string][]byte{
							//	"key": []byte("value"),
							//},
							//&crypto.CryptoValue{},
						},
					},
				),
			},
			object: &Sessions{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Sessions: []*Session{
					{
						ID:            "session-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						State:         domain.SessionStateActive,
						ResourceOwner: "ro",
						Creator:       "creator",
						UserFactor: SessionUserFactor{
							UserID:        "user-id",
							UserCheckedAt: testNow,
							LoginName:     "login-name",
							DisplayName:   "display-name",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						//Metadata: map[string][]byte{
						//	"key": []byte("value"),
						//},
					},
					{
						ID:            "session-id2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211109,
						State:         domain.SessionStateActive,
						ResourceOwner: "ro",
						Creator:       "creator2",
						UserFactor: SessionUserFactor{
							UserID:        "user-id2",
							UserCheckedAt: testNow,
							LoginName:     "login-name2",
							DisplayName:   "display-name2",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						//Metadata: map[string][]byte{
						//	"key": []byte("value"),
						//},
					},
				},
			},
		},
		{
			name:    "prepareSessionsQuery sql err",
			prepare: prepareSessionsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedSessionsQuery,
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

func Test_SessionPrepare(t *testing.T) {
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
			name:    "prepareSessionQuery no result",
			prepare: prepareSessionQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedSessionQuery,
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
			object: (*Session)(nil),
		},
		{
			name:    "prepareSessionQuery found",
			prepare: prepareSessionQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedSessionQuery,
					sessionCols,
					[]driver.Value{
						"session-id",
						testNow,
						testNow,
						uint64(20211109),
						domain.SessionStateActive,
						"ro",
						"creator",
						"user-id",
						testNow,
						"login-name",
						"display-name",
						testNow,
						nil,
						//map[string][]byte{
						//	"key": []byte("value"),
						//},
						//&crypto.CryptoValue{},
					},
				),
			},
			object: &Session{
				ID:            "session-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				State:         domain.SessionStateActive,
				ResourceOwner: "ro",
				Creator:       "creator",
				UserFactor: SessionUserFactor{
					UserID:        "user-id",
					UserCheckedAt: testNow,
					LoginName:     "login-name",
					DisplayName:   "display-name",
				},
				PasswordFactor: SessionPasswordFactor{
					PasswordCheckedAt: testNow,
				},
				//Metadata: map[string][]byte{
				//	"key": []byte("value"),
				//},
			},
		},
		{
			name:    "prepareSessionQuery sql err",
			prepare: prepareSessionQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedSessionQuery,
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
