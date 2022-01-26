package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"golang.org/x/text/language"
)

var (
	iamTable = table{
		name: projection.IAMProjectionTable,
	}
	IAMColumnID = Column{
		name:  projection.IAMColumnID,
		table: iamTable,
	}
	IAMColumnChangeDate = Column{
		name:  projection.IAMColumnChangeDate,
		table: iamTable,
	}
	IAMColumnSequence = Column{
		name:  projection.IAMColumnSequence,
		table: iamTable,
	}
	IAMColumnGlobalOrgID = Column{
		name:  projection.IAMColumnGlobalOrgID,
		table: iamTable,
	}
	IAMColumnProjectID = Column{
		name:  projection.IAMColumnProjectID,
		table: iamTable,
	}
	IAMColumnSetupStarted = Column{
		name:  projection.IAMColumnSetUpStarted,
		table: iamTable,
	}
	IAMColumnSetupDone = Column{
		name:  projection.IAMColumnSetUpDone,
		table: iamTable,
	}
	IAMColumnDefaultLanguage = Column{
		name:  projection.IAMColumnDefaultLanguage,
		table: iamTable,
	}
)

type IAM struct {
	ID         string
	ChangeDate time.Time
	Sequence   uint64

	GlobalOrgID     string
	IAMProjectID    string
	DefaultLanguage language.Tag
	SetupStarted    domain.Step
	SetupDone       domain.Step
}

type IAMSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *IAMSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) IAMByID(ctx context.Context, id string) (*IAM, error) {
	stmt, scan := prepareIAMQuery()
	query, args, err := stmt.Where(sq.Eq{
		IAMColumnID.identifier(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-d9ngs", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) GetDefaultLanguage(ctx context.Context) language.Tag {
	iam, err := q.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return language.Und
	}
	return iam.DefaultLanguage
}

func prepareIAMQuery() (sq.SelectBuilder, func(*sql.Row) (*IAM, error)) {
	return sq.Select(
			IAMColumnID.identifier(),
			IAMColumnChangeDate.identifier(),
			IAMColumnSequence.identifier(),
			IAMColumnGlobalOrgID.identifier(),
			IAMColumnProjectID.identifier(),
			IAMColumnSetupStarted.identifier(),
			IAMColumnSetupDone.identifier(),
			IAMColumnDefaultLanguage.identifier(),
		).
			From(iamTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IAM, error) {
			iam := new(IAM)
			lang := ""
			err := row.Scan(
				&iam.ID,
				&iam.ChangeDate,
				&iam.Sequence,
				&iam.GlobalOrgID,
				&iam.IAMProjectID,
				&iam.SetupStarted,
				&iam.SetupDone,
				&lang,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-n0wng", "Errors.IAM.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-d9nw", "Errors.Internal")
			}
			iam.DefaultLanguage = language.Make(lang)
			return iam, nil
		}
}
