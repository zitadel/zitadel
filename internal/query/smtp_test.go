package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareSMTPConfigStmt = `SELECT projections.smtp_configs2.creation_date,` +
		` projections.smtp_configs2.change_date,` +
		` projections.smtp_configs2.resource_owner,` +
		` projections.smtp_configs2.sequence,` +
		` projections.smtp_configs2.tls,` +
		` projections.smtp_configs2.sender_address,` +
		` projections.smtp_configs2.sender_name,` +
		` projections.smtp_configs2.reply_to_address,` +
		` projections.smtp_configs2.host,` +
		` projections.smtp_configs2.username,` +
		` projections.smtp_configs2.password,` +
		` projections.smtp_configs2.id,` +
		` projections.smtp_configs2.state,` +
		` projections.smtp_configs2.description` +
		` FROM projections.smtp_configs2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareSMTPConfigCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"tls",
		"sender_address",
		"sender_name",
		"reply_to_address",
		"smtp_host",
		"smtp_user",
		"smtp_password",
		"id",
		"state",
		"description",
	}
)

func Test_SMTPConfigsPrepares(t *testing.T) {
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
			name:    "prepareSMTPConfigQuery no result",
			prepare: prepareSMTPConfigQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					prepareSMTPConfigStmt,
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
			object: (*SMTPConfig)(nil),
		},
		{
			name:    "prepareSMTPConfigQuery found",
			prepare: prepareSMTPConfigQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareSMTPConfigStmt),
					prepareSMTPConfigCols,
					[]driver.Value{
						testNow,
						testNow,
						"ro",
						uint64(20211108),
						true,
						"sender",
						"name",
						"reply-to",
						"host",
						"user",
						&crypto.CryptoValue{},
						"2232323",
						domain.SMTPConfigStateActive,
						"test",
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:   testNow,
				ChangeDate:     testNow,
				ResourceOwner:  "ro",
				Sequence:       20211108,
				TLS:            true,
				SenderAddress:  "sender",
				SenderName:     "name",
				ReplyToAddress: "reply-to",
				Host:           "host",
				User:           "user",
				Password:       &crypto.CryptoValue{},
				ID:             "2232323",
				State:          domain.SMTPConfigStateActive,
				Description:    "test",
			},
		},
		{
			name:    "prepareSMTPConfigQuery another config found",
			prepare: prepareSMTPConfigQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareSMTPConfigStmt),
					prepareSMTPConfigCols,
					[]driver.Value{
						testNow,
						testNow,
						"ro",
						uint64(20211109),
						true,
						"sender2",
						"name2",
						"reply-to2",
						"host2",
						"user2",
						&crypto.CryptoValue{},
						"44442323",
						domain.SMTPConfigStateInactive,
						"test2",
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:   testNow,
				ChangeDate:     testNow,
				ResourceOwner:  "ro",
				Sequence:       20211109,
				TLS:            true,
				SenderAddress:  "sender2",
				SenderName:     "name2",
				ReplyToAddress: "reply-to2",
				Host:           "host2",
				User:           "user2",
				Password:       &crypto.CryptoValue{},
				ID:             "44442323",
				State:          domain.SMTPConfigStateInactive,
				Description:    "test2",
			},
		},
		{
			name:    "prepareSMTPConfigQuery yet another config found",
			prepare: prepareSMTPConfigQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareSMTPConfigStmt),
					prepareSMTPConfigCols,
					[]driver.Value{
						testNow,
						testNow,
						"ro",
						uint64(20211109),
						true,
						"sender3",
						"name3",
						"reply-to3",
						"host3",
						"user3",
						&crypto.CryptoValue{},
						"23234444",
						domain.SMTPConfigStateInactive,
						"test3",
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:   testNow,
				ChangeDate:     testNow,
				ResourceOwner:  "ro",
				Sequence:       20211109,
				TLS:            true,
				SenderAddress:  "sender3",
				SenderName:     "name3",
				ReplyToAddress: "reply-to3",
				Host:           "host3",
				User:           "user3",
				Password:       &crypto.CryptoValue{},
				ID:             "23234444",
				State:          domain.SMTPConfigStateInactive,
				Description:    "test3",
			},
		},
		{
			name:    "prepareSMTPConfigQuery sql err",
			prepare: prepareSMTPConfigQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareSMTPConfigStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*SMTPConfig)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
