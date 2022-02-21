package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	errs "github.com/caos/zitadel/internal/errors"
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
					`SELECT zitadel.projections.smtp_configs.aggregate_id,`+
						` zitadel.projections.smtp_configs.creation_date,`+
						` zitadel.projections.smtp_configs.change_date,`+
						` zitadel.projections.smtp_configs.resource_owner,`+
						` zitadel.projections.smtp_configs.sequence,`+
						` zitadel.projections.smtp_configs.tls,`+
						` zitadel.projections.smtp_configs.sender_address,`+
						` zitadel.projections.smtp_configs.sender_name,`+
						` zitadel.projections.smtp_configs.host,`+
						` zitadel.projections.smtp_configs.username,`+
						` zitadel.projections.smtp_configs.password`+
						` FROM zitadel.projections.smtp_configs`,
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
					regexp.QuoteMeta(`SELECT zitadel.projections.smtp_configs.aggregate_id,`+
						` zitadel.projections.smtp_configs.creation_date,`+
						` zitadel.projections.smtp_configs.change_date,`+
						` zitadel.projections.smtp_configs.resource_owner,`+
						` zitadel.projections.smtp_configs.sequence,`+
						` zitadel.projections.smtp_configs.tls,`+
						` zitadel.projections.smtp_configs.sender_address,`+
						` zitadel.projections.smtp_configs.sender_name,`+
						` zitadel.projections.smtp_configs.host,`+
						` zitadel.projections.smtp_configs.username,`+
						` zitadel.projections.smtp_configs.password`+
						` FROM zitadel.projections.smtp_configs`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.smtp_configs.aggregate_id,`+
						` zitadel.projections.smtp_configs.creation_date,`+
						` zitadel.projections.smtp_configs.change_date,`+
						` zitadel.projections.smtp_configs.resource_owner,`+
						` zitadel.projections.smtp_configs.sequence,`+
						` zitadel.projections.smtp_configs.tls,`+
						` zitadel.projections.smtp_configs.sender_address,`+
						` zitadel.projections.smtp_configs.sender_name,`+
						` zitadel.projections.smtp_configs.host,`+
						` zitadel.projections.smtp_configs.username,`+
						` zitadel.projections.smtp_configs.password`+
						` FROM zitadel.projections.smtp_configs`),
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
