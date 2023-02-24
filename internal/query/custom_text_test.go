package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"golang.org/x/text/language"

	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareCustomTextsStmt = `SELECT projections.custom_texts2.aggregate_id,` +
		` projections.custom_texts2.sequence,` +
		` projections.custom_texts2.creation_date,` +
		` projections.custom_texts2.change_date,` +
		` projections.custom_texts2.language,` +
		` projections.custom_texts2.template,` +
		` projections.custom_texts2.key,` +
		` projections.custom_texts2.text,` +
		` COUNT(*) OVER ()` +
		` FROM projections.custom_texts2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareCustomTextsCols = []string{
		"aggregate_id",
		"sequence",
		"creation_date",
		"change_date",
		"language",
		"template",
		"key",
		"text",
		"count",
	}
)

func Test_CustomTextPrepares(t *testing.T) {
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
			name:    "prepareCustomTextQuery no result",
			prepare: prepareCustomTextsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareCustomTextsStmt),
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
			object: &CustomTexts{CustomTexts: []*CustomText{}},
		},
		{
			name:    "prepareCustomTextQuery one result",
			prepare: prepareCustomTextsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareCustomTextsStmt),
					prepareCustomTextsCols,
					[][]driver.Value{
						{
							"agg-id",
							uint64(20211109),
							testNow,
							testNow,
							"en",
							"template",
							"key",
							"text",
						},
					},
				),
			},
			object: &CustomTexts{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				CustomTexts: []*CustomText{
					{
						AggregateID:  "agg-id",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Language:     language.English,
						Template:     "template",
						Key:          "key",
						Text:         "text",
					},
				},
			},
		},
		{
			name:    "prepareCustomTextQuery multiple result",
			prepare: prepareCustomTextsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareCustomTextsStmt),
					prepareCustomTextsCols,
					[][]driver.Value{
						{
							"agg-id",
							uint64(20211109),
							testNow,
							testNow,
							"en",
							"template",
							"key",
							"text",
						},
						{
							"agg-id",
							uint64(20211109),
							testNow,
							testNow,
							"en",
							"template",
							"key2",
							"text",
						},
					},
				),
			},
			object: &CustomTexts{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				CustomTexts: []*CustomText{
					{
						AggregateID:  "agg-id",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Language:     language.English,
						Template:     "template",
						Key:          "key",
						Text:         "text",
					},
					{
						AggregateID:  "agg-id",
						CreationDate: testNow,
						ChangeDate:   testNow,
						Sequence:     20211109,
						Language:     language.English,
						Template:     "template",
						Key:          "key2",
						Text:         "text",
					},
				},
			},
		},
		{
			name:    "prepareCustomTextQuery sql err",
			prepare: prepareCustomTextsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareCustomTextsStmt),
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
