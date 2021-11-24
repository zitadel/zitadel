package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"golang.org/x/text/language"
)

type CustomTexts struct {
	SearchResponse
	CustomTexts []*CustomText
}

type CustomText struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time

	Template string
	Language language.Tag
	Key      string
	Text     string
}

var (
	customTextTable = table{
		name: projection.CustomTextTable,
	}
	CustomTextColAggregateID = Column{
		name: projection.CustomTextAggregateIDCol,
	}
	CustomTextColSequence = Column{
		name: projection.CustomTextSequenceCol,
	}
	CustomTextColCreationDate = Column{
		name: projection.CustomTextCreationDateCol,
	}
	CustomTextColChangeDate = Column{
		name: projection.CustomTextChangeDateCol,
	}
	CustomTextColTemplate = Column{
		name: projection.CustomTextTemplateCol,
	}
	CustomTextColLanguage = Column{
		name: projection.CustomTextLanguageCol,
	}
	CustomTextColKey = Column{
		name: projection.CustomTextKeyCol,
	}
	CustomTextColText = Column{
		name: projection.CustomTextTextCol,
	}
)

func (q *Queries) CustomTextList(ctx context.Context, aggregateID, template, language string) (texts *CustomTexts, err error) {
	stmt, scan := prepareCustomTextsQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			CustomTextColAggregateID.identifier(): aggregateID,
			CustomTextColTemplate.identifier():    template,
			CustomTextColLanguage.identifier():    language,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M9gse", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2j00f", "Errors.Internal")
	}
	texts, err = scan(rows)
	if err != nil {
		return nil, err
	}
	texts.LatestSequence, err = q.latestSequence(ctx, projectsTable)
	return texts, err
}

func prepareCustomTextQuery() (sq.SelectBuilder, func(*sql.Row) (*CustomText, error)) {
	return sq.Select(
			CustomTextColAggregateID.identifier(),
			CustomTextColSequence.identifier(),
			CustomTextColCreationDate.identifier(),
			CustomTextColChangeDate.identifier(),
			CustomTextColLanguage.identifier(),
			CustomTextColTemplate.identifier(),
			CustomTextColKey.identifier(),
			CustomTextColText.identifier(),
		).
			From(customTextTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*CustomText, error) {
			msg := new(CustomText)
			lang := ""
			err := row.Scan(
				&msg.AggregateID,
				&msg.Sequence,
				&msg.CreationDate,
				&msg.ChangeDate,
				&lang,
				&msg.Template,
				&msg.Key,
				&msg.Text,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-3K0ge", "Errors.CustomText.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-2m9gR", "Errors.Internal")
			}
			msg.Language = language.Make(lang)
			return msg, nil
		}
}

func prepareCustomTextsQuery() (sq.SelectBuilder, func(*sql.Rows) (*CustomTexts, error)) {
	return sq.Select(
			CustomTextColAggregateID.identifier(),
			CustomTextColSequence.identifier(),
			CustomTextColCreationDate.identifier(),
			CustomTextColChangeDate.identifier(),
			CustomTextColLanguage.identifier(),
			CustomTextColTemplate.identifier(),
			CustomTextColKey.identifier(),
			CustomTextColText.identifier(),
			countColumn.identifier()).
			From(customTextTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*CustomTexts, error) {
			customTexts := make([]*CustomText, 0)
			var count uint64
			for rows.Next() {
				customText := new(CustomText)
				lang := ""
				err := rows.Scan(
					&customText.AggregateID,
					&customText.Sequence,
					&customText.CreationDate,
					&customText.ChangeDate,
					&lang,
					&customText.Template,
					&customText.Key,
					&customText.Text,
					&count,
				)
				if err != nil {
					return nil, err
				}
				customText.Language = language.Make(lang)
				customTexts = append(customTexts, customText)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-3n9fs", "Errors.Query.CloseRows")
			}

			return &CustomTexts{
				CustomTexts: customTexts,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
