package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareTargetsStmt = `SELECT projections.targets2.id,` +
		` projections.targets2.creation_date,` +
		` projections.targets2.change_date,` +
		` projections.targets2.resource_owner,` +
		` projections.targets2.name,` +
		` projections.targets2.target_type,` +
		` projections.targets2.timeout,` +
		` projections.targets2.endpoint,` +
		` projections.targets2.interrupt_on_error,` +
		` projections.targets2.signing_key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.targets2`
	prepareTargetsCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"name",
		"target_type",
		"timeout",
		"endpoint",
		"interrupt_on_error",
		"signing_key",
		"count",
	}

	prepareTargetStmt = `SELECT projections.targets2.id,` +
		` projections.targets2.creation_date,` +
		` projections.targets2.change_date,` +
		` projections.targets2.resource_owner,` +
		` projections.targets2.name,` +
		` projections.targets2.target_type,` +
		` projections.targets2.timeout,` +
		` projections.targets2.endpoint,` +
		` projections.targets2.interrupt_on_error,` +
		` projections.targets2.signing_key` +
		` FROM projections.targets2`
	prepareTargetCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"name",
		"target_type",
		"timeout",
		"endpoint",
		"interrupt_on_error",
		"signing_key",
	}
)

func Test_TargetPrepares(t *testing.T) {
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
			name:    "prepareTargetsQuery no result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					nil,
					nil,
				),
			},
			object: &Targets{Targets: []*Target{}},
		},
		{
			name:    "prepareTargetsQuery one result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					prepareTargetsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							"target-name",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							true,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
					},
				),
			},
			object: &Targets{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Targets: []*Target{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: true,
						signingKey: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
				},
			},
		},
		{
			name:    "prepareTargetsQuery multiple result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					prepareTargetsCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							"target-name1",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							true,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							"target-name2",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							false,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
						{
							"id-3",
							testNow,
							testNow,
							"ro",
							"target-name3",
							domain.TargetTypeAsync,
							1 * time.Second,
							"https://example.com",
							false,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
						},
					},
				),
			},
			object: &Targets{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				Targets: []*Target{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-1",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name1",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: true,
						signingKey: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-2",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name2",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: false,
						signingKey: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-3",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name3",
						TargetType:       domain.TargetTypeAsync,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: false,
						signingKey: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
				},
			},
		},
		{
			name:    "prepareTargetsQuery sql err",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareTargetsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Target)(nil),
		},
		{
			name:    "prepareTargetQuery no result",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareTargetStmt),
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
			object: (*Target)(nil),
		},
		{
			name:    "prepareTargetQuery found",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareTargetStmt),
					prepareTargetCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"ro",
						"target-name",
						domain.TargetTypeWebhook,
						1 * time.Second,
						"https://example.com",
						true,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						},
					},
				),
			},
			object: &Target{
				ObjectDetails: domain.ObjectDetails{
					ID:            "id",
					EventDate:     testNow,
					CreationDate:  testNow,
					ResourceOwner: "ro",
				},
				Name:             "target-name",
				TargetType:       domain.TargetTypeWebhook,
				Timeout:          1 * time.Second,
				Endpoint:         "https://example.com",
				InterruptOnError: true,
				signingKey: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "alg",
					KeyID:      "encKey",
					Crypted:    []byte("crypted"),
				},
			},
		},
		{
			name:    "prepareTargetQuery sql err",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareTargetStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Target)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
