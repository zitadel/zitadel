package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	instanceTable = table{
		name: projection.InstanceProjectionTable,
	}
	InstanceColumnID = Column{
		name:  projection.InstanceColumnID,
		table: instanceTable,
	}
	InstanceColumnChangeDate = Column{
		name:  projection.InstanceColumnChangeDate,
		table: instanceTable,
	}
	InstanceColumnSequence = Column{
		name:  projection.InstanceColumnSequence,
		table: instanceTable,
	}
	InstanceColumnGlobalOrgID = Column{
		name:  projection.InstanceColumnGlobalOrgID,
		table: instanceTable,
	}
	InstanceColumnProjectID = Column{
		name:  projection.InstanceColumnProjectID,
		table: instanceTable,
	}
	InstanceColumnSetupStarted = Column{
		name:  projection.InstanceColumnSetUpStarted,
		table: instanceTable,
	}
	InstanceColumnSetupDone = Column{
		name:  projection.InstanceColumnSetUpDone,
		table: instanceTable,
	}
	InstanceColumnDefaultLanguage = Column{
		name:  projection.InstanceColumnDefaultLanguage,
		table: instanceTable,
	}
)

type Instance struct {
	ID         string
	ChangeDate time.Time
	Sequence   uint64

	GlobalOrgID     string
	IAMProjectID    string
	DefaultLanguage language.Tag
	SetupStarted    domain.Step
	SetupDone       domain.Step
}

type InstanceSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *InstanceSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) Instance(ctx context.Context) (*Instance, error) {
	stmt, scan := prepareIAMQuery()
	query, args, err := stmt.Where(sq.Eq{
		InstanceColumnID.identifier(): authz.GetInstance(ctx).ID,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-d9ngs", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) GetDefaultLanguage(ctx context.Context) language.Tag {
	iam, err := q.Instance(ctx)
	if err != nil {
		return language.Und
	}
	return iam.DefaultLanguage
}

func prepareIAMQuery() (sq.SelectBuilder, func(*sql.Row) (*Instance, error)) {
	return sq.Select(
			InstanceColumnID.identifier(),
			InstanceColumnChangeDate.identifier(),
			InstanceColumnSequence.identifier(),
			InstanceColumnGlobalOrgID.identifier(),
			InstanceColumnProjectID.identifier(),
			InstanceColumnSetupStarted.identifier(),
			InstanceColumnSetupDone.identifier(),
			InstanceColumnDefaultLanguage.identifier(),
		).
			From(instanceTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Instance, error) {
			iam := new(Instance)
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
