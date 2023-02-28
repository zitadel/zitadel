package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareMessageTextStmt = `SELECT projections.message_texts2.aggregate_id,` +
		` projections.message_texts2.sequence,` +
		` projections.message_texts2.creation_date,` +
		` projections.message_texts2.change_date,` +
		` projections.message_texts2.state,` +
		` projections.message_texts2.type,` +
		` projections.message_texts2.language,` +
		` projections.message_texts2.title,` +
		` projections.message_texts2.pre_header,` +
		` projections.message_texts2.subject,` +
		` projections.message_texts2.greeting,` +
		` projections.message_texts2.text,` +
		` projections.message_texts2.button_text,` +
		` projections.message_texts2.footer_text` +
		` FROM projections.message_texts2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareMessgeTextCols = []string{
		"aggregate_id",
		"sequence",
		"creation_date",
		"change_date",
		"state",
		"type",
		"language",
		"title",
		"pre_header",
		"subject",
		"greeting",
		"text",
		"button_text",
		"footer_text",
	}
)

func Test_MessageTextPrepares(t *testing.T) {
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
			name:    "prepareMessageTextQuery no result",
			prepare: prepareMessageTextQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareMessageTextStmt),
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
			object: (*MessageText)(nil),
		},
		{
			name:    "prepareMesssageTextQuery found",
			prepare: prepareMessageTextQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareMessageTextStmt),
					prepareMessgeTextCols,
					[]driver.Value{
						"agg-id",
						uint64(20211109),
						testNow,
						testNow,
						domain.PolicyStateActive,
						"type",
						"en",
						"title",
						"pre_header",
						"subject",
						"greeting",
						"text",
						"button_text",
						"footer_text",
					},
				),
			},
			object: &MessageText{
				AggregateID:  "agg-id",
				CreationDate: testNow,
				ChangeDate:   testNow,
				Sequence:     20211109,
				State:        domain.PolicyStateActive,
				Type:         "type",
				Language:     language.English,
				Title:        "title",
				PreHeader:    "pre_header",
				Subject:      "subject",
				Greeting:     "greeting",
				Text:         "text",
				ButtonText:   "button_text",
				Footer:       "footer_text",
			},
		},
		{
			name:    "prepareMessageTextQuery sql err",
			prepare: prepareMessageTextQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareMessageTextStmt),
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
