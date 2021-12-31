package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
	"golang.org/x/text/language"
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
					regexp.QuoteMeta(`SELECT zitadel.projections.message_texts.aggregate_id,`+
						` zitadel.projections.message_texts.sequence,`+
						` zitadel.projections.message_texts.creation_date,`+
						` zitadel.projections.message_texts.change_date,`+
						` zitadel.projections.message_texts.state,`+
						` zitadel.projections.message_texts.type,`+
						` zitadel.projections.message_texts.language,`+
						` zitadel.projections.message_texts.title,`+
						` zitadel.projections.message_texts.pre_header,`+
						` zitadel.projections.message_texts.subject,`+
						` zitadel.projections.message_texts.greeting,`+
						` zitadel.projections.message_texts.text,`+
						` zitadel.projections.message_texts.button_text,`+
						` zitadel.projections.message_texts.footer_text`+
						` FROM zitadel.projections.message_texts`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.message_texts.aggregate_id,`+
						` zitadel.projections.message_texts.sequence,`+
						` zitadel.projections.message_texts.creation_date,`+
						` zitadel.projections.message_texts.change_date,`+
						` zitadel.projections.message_texts.state,`+
						` zitadel.projections.message_texts.type,`+
						` zitadel.projections.message_texts.language,`+
						` zitadel.projections.message_texts.title,`+
						` zitadel.projections.message_texts.pre_header,`+
						` zitadel.projections.message_texts.subject,`+
						` zitadel.projections.message_texts.greeting,`+
						` zitadel.projections.message_texts.text,`+
						` zitadel.projections.message_texts.button_text,`+
						` zitadel.projections.message_texts.footer_text`+
						` FROM zitadel.projections.message_texts`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.message_texts.aggregate_id,`+
						` zitadel.projections.message_texts.sequence,`+
						` zitadel.projections.message_texts.creation_date,`+
						` zitadel.projections.message_texts.change_date,`+
						` zitadel.projections.message_texts.state,`+
						` zitadel.projections.message_texts.type,`+
						` zitadel.projections.message_texts.language,`+
						` zitadel.projections.message_texts.title,`+
						` zitadel.projections.message_texts.pre_header,`+
						` zitadel.projections.message_texts.subject,`+
						` zitadel.projections.message_texts.greeting,`+
						` zitadel.projections.message_texts.text,`+
						` zitadel.projections.message_texts.button_text,`+
						` zitadel.projections.message_texts.footer_text`+
						` FROM zitadel.projections.message_texts`),
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
