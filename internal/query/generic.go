package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func genericSearch[R Stateful](
	q *Queries,
	ctx context.Context,
	projection table,
	prepareQuery func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (R, error)),
	toQuery func(query sq.SelectBuilder) sq.SelectBuilder,
) (resp R, err error) {
	var rnil R
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareQuery(ctx, q.client)
	stmt, args, err := toQuery(query).ToSql()
	if err != nil {
		return rnil, zerrors.ThrowInvalidArgument(err, "QUERY-SDgwg", "Errors.Query.InvalidRequest")
	}
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		resp, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return rnil, zerrors.ThrowInternal(err, "QUERY-SDfr52", "Errors.Internal")
	}
	state, err := q.latestState(ctx, projection)
	if err != nil {
		return rnil, err
	}
	resp.SetState(state)
	return resp, err
}

func genericGetByID[R any](
	q *Queries,
	ctx context.Context,
	prepareQuery func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (R, error)),
	toQuery func(query sq.SelectBuilder) sq.SelectBuilder,
) (resp R, err error) {
	var rnil R
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareQuery(ctx, q.client)
	stmt, args, err := toQuery(query).ToSql()
	if err != nil {
		return rnil, zerrors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		resp, err = scan(row)
		return err
	}, stmt, args...)
	return resp, err
}

func queryToWhereStmt(eq interface{}) func(query sq.SelectBuilder) sq.SelectBuilder {
	return func(query sq.SelectBuilder) sq.SelectBuilder {
		return query.Where(eq)
	}
}

func combineToWhereStmt(toQuery func(query sq.SelectBuilder) sq.SelectBuilder, eq interface{}) func(query sq.SelectBuilder) sq.SelectBuilder {
	return func(query sq.SelectBuilder) sq.SelectBuilder {
		return toQuery(query).Where(eq)
	}
}
