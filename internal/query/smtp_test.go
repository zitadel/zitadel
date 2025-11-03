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
	prepareSMTPConfigStmt = `SELECT projections.smtp_configs5.creation_date,` +
		` projections.smtp_configs5.change_date,` +
		` projections.smtp_configs5.resource_owner,` +
		` projections.smtp_configs5.sequence,` +
		` projections.smtp_configs5.id,` +
		` projections.smtp_configs5.state,` +
		` projections.smtp_configs5.description,` +
		` projections.smtp_configs5_smtp.id,` +
		` projections.smtp_configs5_smtp.tls,` +
		` projections.smtp_configs5_smtp.sender_address,` +
		` projections.smtp_configs5_smtp.sender_name,` +
		` projections.smtp_configs5_smtp.reply_to_address,` +
		` projections.smtp_configs5_smtp.host,` +
		` projections.smtp_configs5_smtp.username,` +
		` projections.smtp_configs5_smtp.password,` +
		` projections.smtp_configs5_http.id,` +
		` projections.smtp_configs5_http.endpoint,` +
		` projections.smtp_configs5_http.signing_key` +
		` FROM projections.smtp_configs5` +
		` LEFT JOIN projections.smtp_configs5_smtp ON projections.smtp_configs5.id = projections.smtp_configs5_smtp.id AND projections.smtp_configs5.instance_id = projections.smtp_configs5_smtp.instance_id` +
		` LEFT JOIN projections.smtp_configs5_http ON projections.smtp_configs5.id = projections.smtp_configs5_http.id AND projections.smtp_configs5.instance_id = projections.smtp_configs5_http.instance_id`
	prepareSMTPConfigsStmt = `SELECT projections.smtp_configs5.creation_date,` +
		` projections.smtp_configs5.change_date,` +
		` projections.smtp_configs5.resource_owner,` +
		` projections.smtp_configs5.sequence,` +
		` projections.smtp_configs5.id,` +
		` projections.smtp_configs5.state,` +
		` projections.smtp_configs5.description,` +
		` projections.smtp_configs5_smtp.id,` +
		` projections.smtp_configs5_smtp.tls,` +
		` projections.smtp_configs5_smtp.sender_address,` +
		` projections.smtp_configs5_smtp.sender_name,` +
		` projections.smtp_configs5_smtp.reply_to_address,` +
		` projections.smtp_configs5_smtp.host,` +
		` projections.smtp_configs5_smtp.username,` +
		` projections.smtp_configs5_smtp.password,` +
		` projections.smtp_configs5_http.id,` +
		` projections.smtp_configs5_http.endpoint,` +
		` projections.smtp_configs5_http.signing_key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.smtp_configs5` +
		` LEFT JOIN projections.smtp_configs5_smtp ON projections.smtp_configs5.id = projections.smtp_configs5_smtp.id AND projections.smtp_configs5.instance_id = projections.smtp_configs5_smtp.instance_id` +
		` LEFT JOIN projections.smtp_configs5_http ON projections.smtp_configs5.id = projections.smtp_configs5_http.id AND projections.smtp_configs5.instance_id = projections.smtp_configs5_http.instance_id`

	prepareSMTPConfigCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"id",
		"state",
		"description",
		"id",
		"tls",
		"sender_address",
		"sender_name",
		"reply_to_address",
		"smtp_host",
		"smtp_user",
		"smtp_password",
		"id",
		"endpoint",
		"signing_key",
	}
	prepareSMTPConfigsCols = append(prepareSMTPConfigCols, "count")
)

func Test_SMTPConfigPrepares(t *testing.T) {
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
						"2232323",
						domain.SMTPConfigStateActive,
						"test",
						"2232323",
						true,
						"sender",
						"name",
						"reply-to",
						"host",
						"user",
						&crypto.CryptoValue{},
						nil,
						nil,
						nil,
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211108,
				SMTPConfig: &SMTP{
					TLS:            true,
					SenderAddress:  "sender",
					SenderName:     "name",
					ReplyToAddress: "reply-to",
					Host:           "host",
					User:           "user",
					Password:       &crypto.CryptoValue{},
				},
				ID:          "2232323",
				State:       domain.SMTPConfigStateActive,
				Description: "test",
			},
		},
		{
			name:    "prepareSMTPConfigQuery found, http",
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
						"2232323",
						domain.SMTPConfigStateActive,
						"test",
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						"2232323",
						"endpoint",
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211108,
				HTTPConfig: &HTTP{
					Endpoint: "endpoint",
					signingKey: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "alg",
						KeyID:      "encKey",
						Crypted:    []byte("crypted"),
					},
				},
				ID:          "2232323",
				State:       domain.SMTPConfigStateActive,
				Description: "test",
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
						"44442323",
						domain.SMTPConfigStateInactive,
						"test2",
						"44442323",
						true,
						"sender2",
						"name2",
						"reply-to2",
						"host2",
						"user2",
						&crypto.CryptoValue{},
						nil,
						nil,
						nil,
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211109,
				SMTPConfig: &SMTP{
					TLS:            true,
					SenderAddress:  "sender2",
					SenderName:     "name2",
					ReplyToAddress: "reply-to2",
					Host:           "host2",
					User:           "user2",
					Password:       &crypto.CryptoValue{},
				},
				ID:          "44442323",
				State:       domain.SMTPConfigStateInactive,
				Description: "test2",
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
						"23234444",
						domain.SMTPConfigStateInactive,
						"test3",
						"23234444",
						true,
						"sender3",
						"name3",
						"reply-to3",
						"host3",
						"user3",
						&crypto.CryptoValue{},
						nil,
						nil,
						nil,
					},
				),
			},
			object: &SMTPConfig{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211109,
				SMTPConfig: &SMTP{
					TLS:            true,
					SenderAddress:  "sender3",
					SenderName:     "name3",
					ReplyToAddress: "reply-to3",
					Host:           "host3",
					User:           "user3",
					Password:       &crypto.CryptoValue{},
				},
				ID:          "23234444",
				State:       domain.SMTPConfigStateInactive,
				Description: "test3",
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

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
			name:    "prepareSMTPConfigsQuery no result",
			prepare: prepareSMTPConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSMTPConfigsStmt),
					nil,
					nil,
				),
			},
			object: &SMTPConfigs{Configs: []*SMTPConfig{}},
		},
		{
			name:    "prepareSMTPConfigsQuery found",
			prepare: prepareSMTPConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSMTPConfigsStmt),
					prepareSMTPConfigsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							"2232323",
							domain.SMTPConfigStateActive,
							"test",
							"2232323",
							true,
							"sender",
							"name",
							"reply-to",
							"host",
							"user",
							&crypto.CryptoValue{},
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &SMTPConfigs{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Configs: []*SMTPConfig{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211108,
						SMTPConfig: &SMTP{
							TLS:            true,
							SenderAddress:  "sender",
							SenderName:     "name",
							ReplyToAddress: "reply-to",
							Host:           "host",
							User:           "user",
							Password:       &crypto.CryptoValue{},
						},
						ID:          "2232323",
						State:       domain.SMTPConfigStateActive,
						Description: "test",
					},
				},
			},
		},
		{
			name:    "prepareSMTPConfigsQuery found, multi",
			prepare: prepareSMTPConfigsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSMTPConfigsStmt),
					prepareSMTPConfigsCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							"2232323",
							domain.SMTPConfigStateActive,
							"test",
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							"2232323",
							"endpoint",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
						{
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							"44442323",
							domain.SMTPConfigStateInactive,
							"test2",
							"44442323",
							true,
							"sender2",
							"name2",
							"reply-to2",
							"host2",
							"user2",
							&crypto.CryptoValue{},
							nil,
							nil,
							nil,
						},
						{
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							"23234444",
							domain.SMTPConfigStateInactive,
							"test3",
							"23234444",
							true,
							"sender3",
							"name3",
							"reply-to3",
							"host3",
							"user3",
							&crypto.CryptoValue{},
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &SMTPConfigs{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				Configs: []*SMTPConfig{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211108,
						HTTPConfig: &HTTP{
							Endpoint: "endpoint",
							signingKey: &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
						ID:          "2232323",
						State:       domain.SMTPConfigStateActive,
						Description: "test",
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						SMTPConfig: &SMTP{
							TLS:            true,
							SenderAddress:  "sender2",
							SenderName:     "name2",
							ReplyToAddress: "reply-to2",
							Host:           "host2",
							User:           "user2",
							Password:       &crypto.CryptoValue{},
						},
						ID:          "44442323",
						State:       domain.SMTPConfigStateInactive,
						Description: "test2",
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						SMTPConfig: &SMTP{
							TLS:            true,
							SenderAddress:  "sender3",
							SenderName:     "name3",
							ReplyToAddress: "reply-to3",
							Host:           "host3",
							User:           "user3",
							Password:       &crypto.CryptoValue{},
						},
						ID:          "23234444",
						State:       domain.SMTPConfigStateInactive,
						Description: "test3",
					},
				},
			},
		},
		{
			name:    "prepareSMTPConfigsQuery sql err",
			prepare: prepareSMTPConfigsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareSMTPConfigsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*SMTPConfigs)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
