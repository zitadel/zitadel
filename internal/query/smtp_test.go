package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareSMTPConfigStmt = `SELECT projections.smtp_configs.aggregate_id,` +
		` projections.smtp_configs.creation_date,` +
		` projections.smtp_configs.change_date,` +
		` projections.smtp_configs.resource_owner,` +
		` projections.smtp_configs.sequence,` +
		` projections.smtp_configs.tls,` +
		` projections.smtp_configs.sender_address,` +
		` projections.smtp_configs.sender_name,` +
		` projections.smtp_configs.host,` +
		` projections.smtp_configs.username,` +
		` projections.smtp_configs.password` +
		` FROM projections.smtp_configs` +
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
				sqlExpectations: mockQueries(
					prepareSMTPConfigStmt,
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
						"host",
						"user",
						&crypto.CryptoValue{},
					},
				),
			},
			object: &SMTPConfig{
				AggregateID:   "agg-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211108,
				TLS:           true,
				SenderAddress: "sender",
				SenderName:    "name",
				Host:          "host",
				User:          "user",
				Password:      &crypto.CryptoValue{},
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
