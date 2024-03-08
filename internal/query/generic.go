package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SearchResponses interface {
	*Targets | *Executions
	Stateful
}

func genericSearch[R SearchResponses](
	q *Queries,
	ctx context.Context,
	projection table,
	prepareQuery func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (R, error)),
	toQuery func(query sq.SelectBuilder) sq.SelectBuilder,
) (resp R, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareQuery(ctx, q.client)
	stmt, args, err := toQuery(query).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SDgwg", "Errors.Query.InvalidRequest")
	}
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		resp, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SDfr52", "Errors.Internal")
	}
	state, err := q.latestState(ctx, projection)
	if err != nil {
		return nil, err
	}
	resp.SetState(state)
	return resp, err
}

type GetByIDResponse interface {
	*Target | *Execution
}

func genericGetByID[R GetByIDResponse](
	q *Queries,
	ctx context.Context,
	prepareQuery func(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (R, error)),
	toQuery func(query sq.SelectBuilder) sq.SelectBuilder,
) (resp R, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareQuery(ctx, q.client)
	stmt, args, err := toQuery(query).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgff3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		resp, err = scan(row)
		return err
	}, stmt, args...)
	return resp, err
}

func where(eq interface{}) func(query sq.SelectBuilder) sq.SelectBuilder {
	return func(query sq.SelectBuilder) sq.SelectBuilder {
		return query.Where(eq)
	}
}

func whereWrapper(toQuery func(query sq.SelectBuilder) sq.SelectBuilder, eq interface{}) func(query sq.SelectBuilder) sq.SelectBuilder {
	return func(query sq.SelectBuilder) sq.SelectBuilder {
		return toQuery(query).Where(eq)
	}
}
