package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareSMTPConfigStmt = `SELECT projections.smtp_configs1.aggregate_id,` +
		` projections.smtp_configs1.creation_date,` +
		` projections.smtp_configs1.change_date,` +
		` projections.smtp_configs1.resource_owner,` +
		` projections.smtp_configs1.sequence,` +
		` projections.smtp_configs1.tls,` +
		` projections.smtp_configs1.sender_address,` +
		` projections.smtp_configs1.sender_name,` +
		` projections.smtp_configs1.reply_to_address,` +
		` projections.smtp_configs1.host,` +
		` projections.smtp_configs1.username,` +
		` projections.smtp_configs1.password` +
		` FROM projections.smtp_configs1` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareSMTPConfigCols = []string{
		"aggregate_id",
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
						"agg-id",
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
					},
				),
			},
			object: &SMTPConfig{
				AggregateID:    "agg-id",
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
