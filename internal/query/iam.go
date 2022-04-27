package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
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
)

type IAM struct {
	ID         string
	ChangeDate time.Time
	Sequence   uint64

	GlobalOrgID  string
	IAMProjectID string
	SetupStarted domain.Step
	SetupDone    domain.Step
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

func prepareIAMQuery() (sq.SelectBuilder, func(*sql.Row) (*IAM, error)) {
	return sq.Select(
			IAMColumnID.identifier(),
			IAMColumnChangeDate.identifier(),
			IAMColumnSequence.identifier(),
			IAMColumnGlobalOrgID.identifier(),
			IAMColumnProjectID.identifier(),
			IAMColumnSetupStarted.identifier(),
			IAMColumnSetupDone.identifier(),
		).
			From(iamTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IAM, error) {
			o := new(IAM)
			err := row.Scan(
				&o.ID,
				&o.ChangeDate,
				&o.Sequence,
				&o.GlobalOrgID,
				&o.IAMProjectID,
				&o.SetupStarted,
				&o.SetupDone,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-n0wng", "Errors.IAM.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-d9nw", "Errors.Internal")
			}
			return o, nil
		}
}
