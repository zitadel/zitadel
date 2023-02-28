package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareSecretGeneratorStmt = `SELECT projections.secret_generators2.aggregate_id,` +
		` projections.secret_generators2.generator_type,` +
		` projections.secret_generators2.creation_date,` +
		` projections.secret_generators2.change_date,` +
		` projections.secret_generators2.resource_owner,` +
		` projections.secret_generators2.sequence,` +
		` projections.secret_generators2.length,` +
		` projections.secret_generators2.expiry,` +
		` projections.secret_generators2.include_lower_letters,` +
		` projections.secret_generators2.include_upper_letters,` +
		` projections.secret_generators2.include_digits,` +
		` projections.secret_generators2.include_symbols` +
		` FROM projections.secret_generators2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareSecretGeneratorCols = []string{
		"aggregate_id",
		"generator_type",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"length",
		"expiry",
		"include_lower_letters",
		"include_upper_letters",
		"include_digits",
		"include_symbols",
	}
	prepareSecretGeneratorsStmt = `SELECT projections.secret_generators2.aggregate_id,` +
		` projections.secret_generators2.generator_type,` +
		` projections.secret_generators2.creation_date,` +
		` projections.secret_generators2.change_date,` +
		` projections.secret_generators2.resource_owner,` +
		` projections.secret_generators2.sequence,` +
		` projections.secret_generators2.length,` +
		` projections.secret_generators2.expiry,` +
		` projections.secret_generators2.include_lower_letters,` +
		` projections.secret_generators2.include_upper_letters,` +
		` projections.secret_generators2.include_digits,` +
		` projections.secret_generators2.include_symbols,` +
		` COUNT(*) OVER ()` +
		` FROM projections.secret_generators2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareSecretGeneratorsCols = []string{
		"aggregate_id",
		"generator_type",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"length",
		"expiry",
		"include_lower_letters",
		"include_upper_letters",
		"include_digits",
		"include_symbols",
		"count",
	}
)

func Test_SecretGeneratorsPrepares(t *testing.T) {
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
			name:    "prepareSecretGeneratorsQuery no result",
			prepare: prepareSecretGeneratorsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSecretGeneratorsStmt),
					nil,
					nil,
				),
			},
			object: &SecretGenerators{SecretGenerators: []*SecretGenerator{}},
		},
		{
			name:    "prepareSecretGeneratorsQuery one result",
			prepare: prepareSecretGeneratorsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSecretGeneratorsStmt),
					prepareSecretGeneratorsCols,
					[][]driver.Value{
						{
							"agg-id",
							domain.SecretGeneratorTypeInitCode,
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							4,
							time.Minute * 1,
							true,
							true,
							true,
							true,
						},
					},
				),
			},
			object: &SecretGenerators{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				SecretGenerators: []*SecretGenerator{
					{
						AggregateID:         "agg-id",
						GeneratorType:       1,
						CreationDate:        testNow,
						ChangeDate:          testNow,
						ResourceOwner:       "ro",
						Sequence:            20211108,
						Length:              4,
						Expiry:              time.Minute * 1,
						IncludeLowerLetters: true,
						IncludeUpperLetters: true,
						IncludeDigits:       true,
						IncludeSymbols:      true,
					},
				},
			},
		},
		{
			name:    "prepareSecretGeneratorsQuery multiple result",
			prepare: prepareSecretGeneratorsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareSecretGeneratorsStmt),
					prepareSecretGeneratorsCols,
					[][]driver.Value{
						{
							"agg-id",
							domain.SecretGeneratorTypeInitCode,
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							4,
							time.Minute * 1,
							true,
							true,
							true,
							true,
						},
						{
							"agg-id",
							domain.SecretGeneratorTypeVerifyEmailCode,
							testNow,
							testNow,
							"ro",
							uint64(20211108),
							4,
							time.Minute * 1,
							true,
							true,
							true,
							true,
						},
					},
				),
			},
			object: &SecretGenerators{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				SecretGenerators: []*SecretGenerator{
					{
						AggregateID:         "agg-id",
						GeneratorType:       1,
						CreationDate:        testNow,
						ChangeDate:          testNow,
						ResourceOwner:       "ro",
						Sequence:            20211108,
						Length:              4,
						Expiry:              time.Minute * 1,
						IncludeLowerLetters: true,
						IncludeUpperLetters: true,
						IncludeDigits:       true,
						IncludeSymbols:      true,
					},
					{
						AggregateID:         "agg-id",
						GeneratorType:       2,
						CreationDate:        testNow,
						ChangeDate:          testNow,
						ResourceOwner:       "ro",
						Sequence:            20211108,
						Length:              4,
						Expiry:              time.Minute * 1,
						IncludeLowerLetters: true,
						IncludeUpperLetters: true,
						IncludeDigits:       true,
						IncludeSymbols:      true,
					},
				},
			},
		},
		{
			name:    "prepareSecretGeneratorsQuery sql err",
			prepare: prepareSecretGeneratorsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareSecretGeneratorsStmt),
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
			name:    "prepareSecretGeneratorQuery no result",
			prepare: prepareSecretGeneratorQuery,
			want: want{
				sqlExpectations: mockQueries(
					prepareSecretGeneratorStmt,
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
			object: (*SecretGenerator)(nil),
		},
		{
			name:    "prepareSecretGeneratorQuery found",
			prepare: prepareSecretGeneratorQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareSecretGeneratorStmt),
					prepareSecretGeneratorCols,
					[]driver.Value{
						"agg-id",
						domain.SecretGeneratorTypeInitCode,
						testNow,
						testNow,
						"ro",
						uint64(20211108),
						4,
						time.Minute * 1,
						true,
						true,
						true,
						true,
					},
				),
			},
			object: &SecretGenerator{
				AggregateID:         "agg-id",
				GeneratorType:       domain.SecretGeneratorTypeInitCode,
				CreationDate:        testNow,
				ChangeDate:          testNow,
				ResourceOwner:       "ro",
				Sequence:            20211108,
				Length:              4,
				Expiry:              time.Minute * 1,
				IncludeLowerLetters: true,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      true,
			},
		},
		{
			name:    "prepareSecretGeneratorQuery sql err",
			prepare: prepareSecretGeneratorQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareSecretGeneratorStmt),
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
