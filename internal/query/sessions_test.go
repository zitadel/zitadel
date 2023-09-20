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

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	expectedSessionQuery = regexp.QuoteMeta(`SELECT projections.sessions5.id,` +
		` projections.sessions5.creation_date,` +
		` projections.sessions5.change_date,` +
		` projections.sessions5.sequence,` +
		` projections.sessions5.state,` +
		` projections.sessions5.resource_owner,` +
		` projections.sessions5.creator,` +
		` projections.sessions5.user_id,` +
		` projections.sessions5.user_checked_at,` +
		` projections.login_names2.login_name,` +
		` projections.users8_humans.display_name,` +
		` projections.users8.resource_owner,` +
		` projections.sessions5.password_checked_at,` +
		` projections.sessions5.intent_checked_at,` +
		` projections.sessions5.webauthn_checked_at,` +
		` projections.sessions5.webauthn_user_verified,` +
		` projections.sessions5.totp_checked_at,` +
		` projections.sessions5.otp_sms_checked_at,` +
		` projections.sessions5.otp_email_checked_at,` +
		` projections.sessions5.metadata,` +
		` projections.sessions5.token_id` +
		` FROM projections.sessions5` +
		` LEFT JOIN projections.login_names2 ON projections.sessions5.user_id = projections.login_names2.user_id AND projections.sessions5.instance_id = projections.login_names2.instance_id` +
		` LEFT JOIN projections.users8_humans ON projections.sessions5.user_id = projections.users8_humans.user_id AND projections.sessions5.instance_id = projections.users8_humans.instance_id` +
		` LEFT JOIN projections.users8 ON projections.sessions5.user_id = projections.users8.id AND projections.sessions5.instance_id = projections.users8.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedSessionsQuery = regexp.QuoteMeta(`SELECT projections.sessions5.id,` +
		` projections.sessions5.creation_date,` +
		` projections.sessions5.change_date,` +
		` projections.sessions5.sequence,` +
		` projections.sessions5.state,` +
		` projections.sessions5.resource_owner,` +
		` projections.sessions5.creator,` +
		` projections.sessions5.user_id,` +
		` projections.sessions5.user_checked_at,` +
		` projections.login_names2.login_name,` +
		` projections.users8_humans.display_name,` +
		` projections.users8.resource_owner,` +
		` projections.sessions5.password_checked_at,` +
		` projections.sessions5.intent_checked_at,` +
		` projections.sessions5.webauthn_checked_at,` +
		` projections.sessions5.webauthn_user_verified,` +
		` projections.sessions5.totp_checked_at,` +
		` projections.sessions5.otp_sms_checked_at,` +
		` projections.sessions5.otp_email_checked_at,` +
		` projections.sessions5.metadata,` +
		` COUNT(*) OVER ()` +
		` FROM projections.sessions5` +
		` LEFT JOIN projections.login_names2 ON projections.sessions5.user_id = projections.login_names2.user_id AND projections.sessions5.instance_id = projections.login_names2.instance_id` +
		` LEFT JOIN projections.users8_humans ON projections.sessions5.user_id = projections.users8_humans.user_id AND projections.sessions5.instance_id = projections.users8_humans.instance_id` +
		` LEFT JOIN projections.users8 ON projections.sessions5.user_id = projections.users8.id AND projections.sessions5.instance_id = projections.users8.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)

	sessionCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"resource_owner",
		"creator",
		"user_id",
		"user_checked_at",
		"login_name",
		"display_name",
		"user_resource_owner",
		"password_checked_at",
		"intent_checked_at",
		"webauthn_checked_at",
		"webauthn_user_verified",
		"totp_checked_at",
		"otp_sms_checked_at",
		"otp_email_checked_at",
		"metadata",
		"token",
	}

	sessionsCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"state",
		"resource_owner",
		"creator",
		"user_id",
		"user_checked_at",
		"login_name",
		"display_name",
		"user_resource_owner",
		"password_checked_at",
		"intent_checked_at",
		"webauthn_checked_at",
		"webauthn_user_verified",
		"totp_checked_at",
		"otp_sms_checked_at",
		"otp_email_checked_at",
		"metadata",
		"count",
	}
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
							"resourceOwner",
							testNow,
							testNow,
							testNow,
							true,
							testNow,
							testNow,
							testNow,
							[]byte(`{"key": "dmFsdWU="}`),
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
							ResourceOwner: "resourceOwner",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						IntentFactor: SessionIntentFactor{
							IntentCheckedAt: testNow,
						},
						WebAuthNFactor: SessionWebAuthNFactor{
							WebAuthNCheckedAt: testNow,
							UserVerified:      true,
						},
						TOTPFactor: SessionTOTPFactor{
							TOTPCheckedAt: testNow,
						},
						OTPSMSFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
						},
						OTPEmailFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
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
							"resourceOwner",
							testNow,
							testNow,
							testNow,
							true,
							testNow,
							testNow,
							testNow,
							[]byte(`{"key": "dmFsdWU="}`),
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
							"resourceOwner",
							testNow,
							testNow,
							testNow,
							false,
							testNow,
							testNow,
							testNow,
							[]byte(`{"key": "dmFsdWU="}`),
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
							ResourceOwner: "resourceOwner",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						IntentFactor: SessionIntentFactor{
							IntentCheckedAt: testNow,
						},
						WebAuthNFactor: SessionWebAuthNFactor{
							WebAuthNCheckedAt: testNow,
							UserVerified:      true,
						},
						TOTPFactor: SessionTOTPFactor{
							TOTPCheckedAt: testNow,
						},
						OTPSMSFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
						},
						OTPEmailFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
						},
						Metadata: map[string][]byte{
							"key": []byte("value"),
						},
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
							ResourceOwner: "resourceOwner",
						},
						PasswordFactor: SessionPasswordFactor{
							PasswordCheckedAt: testNow,
						},
						IntentFactor: SessionIntentFactor{
							IntentCheckedAt: testNow,
						},
						WebAuthNFactor: SessionWebAuthNFactor{
							WebAuthNCheckedAt: testNow,
							UserVerified:      false,
						},
						TOTPFactor: SessionTOTPFactor{
							TOTPCheckedAt: testNow,
						},
						OTPSMSFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
						},
						OTPEmailFactor: SessionOTPFactor{
							OTPCheckedAt: testNow,
						},
						Metadata: map[string][]byte{
							"key": []byte("value"),
						},
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
			object: (*Sessions)(nil),
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
			prepare: prepareSessionQueryTesting(t, ""),
			want: want{
				sqlExpectations: mockQueriesScanErr(
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
			prepare: prepareSessionQueryTesting(t, "tokenID"),
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
						"resourceOwner",
						testNow,
						testNow,
						testNow,
						true,
						testNow,
						testNow,
						testNow,
						[]byte(`{"key": "dmFsdWU="}`),
						"tokenID",
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
					ResourceOwner: "resourceOwner",
				},
				PasswordFactor: SessionPasswordFactor{
					PasswordCheckedAt: testNow,
				},
				IntentFactor: SessionIntentFactor{
					IntentCheckedAt: testNow,
				},
				WebAuthNFactor: SessionWebAuthNFactor{
					WebAuthNCheckedAt: testNow,
					UserVerified:      true,
				},
				TOTPFactor: SessionTOTPFactor{
					TOTPCheckedAt: testNow,
				},
				OTPSMSFactor: SessionOTPFactor{
					OTPCheckedAt: testNow,
				},
				OTPEmailFactor: SessionOTPFactor{
					OTPCheckedAt: testNow,
				},
				Metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
		},
		{
			name:    "prepareSessionQuery sql err",
			prepare: prepareSessionQueryTesting(t, ""),
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
			object: (*Session)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

func prepareSessionQueryTesting(t *testing.T, token string) func(context.Context, prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Session, error)) {
	return func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Session, error)) {
		builder, scan := prepareSessionQuery(ctx, db)
		return builder, func(row *sql.Row) (*Session, error) {
			session, tokenID, err := scan(row)
			require.Equal(t, tokenID, token)
			return session, err
		}
	}
}
