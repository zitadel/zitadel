package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func genericRowsQuery[R any](
	ctx context.Context,
	client *database.DB,
	query sq.SelectBuilder,
	scan func(rows *sql.Rows) (R, error),
) (resp R, err error) {
	var rnil R
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, args, err := query.ToSql()
	if err != nil {
		return rnil, zerrors.ThrowInvalidArgument(err, "QUERY-05wf2q36ji", "Errors.Query.InvalidRequest")
	}
	err = client.QueryContext(ctx, func(rows *sql.Rows) error {
		resp, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return rnil, zerrors.ThrowInternal(err, "QUERY-y2u7vctrha", "Errors.Internal")
	}
	return resp, err
}

func genericRowsQueryWithState[R Stateful](
	ctx context.Context,
	client *database.DB,
	projection table,
	query sq.SelectBuilder,
	scan func(rows *sql.Rows) (R, error),
) (resp R, err error) {
	var rnil R
	resp, err = genericRowsQuery(ctx, client, query, scan)
	if err != nil {
		return rnil, err
	}
	state, err := latestState(ctx, client, projection)
	if err != nil {
		return rnil, err
	}
	resp.SetState(state)
	return resp, err
}

func latestState(ctx context.Context, client *database.DB, projections ...table) (state *State, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareLatestState()
	or := make(sq.Or, len(projections))
	for i, projection := range projections {
		or[i] = sq.Eq{CurrentStateColProjectionName.identifier(): projection.name}
	}
	stmt, args, err := query.
		Where(or).
		Where(sq.Eq{CurrentStateColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}).
		OrderBy(CurrentStateColEventDate.identifier() + " DESC").
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-5CfX9", "Errors.Query.SQLStatement")
	}

	err = client.QueryRowContext(ctx, func(row *sql.Row) error {
		state, err = scan(row)
		return err
	}, stmt, args...)

	return state, err
}

func genericRowQuery[R any](
	ctx context.Context,
	client *database.DB,
	query sq.SelectBuilder,
	scan func(row *sql.Row) (R, error),
) (resp R, err error) {
	var rnil R
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, args, err := query.ToSql()
	if err != nil {
		return rnil, zerrors.ThrowInternal(err, "QUERY-s969t763z4", "Errors.Query.SQLStatement")
	}

	err = client.QueryRowContext(ctx, func(row *sql.Row) error {
		resp, err = scan(row)
		return err
	}, stmt, args...)
	return resp, err
}

func combineToWhereStmt(query sq.SelectBuilder, toQuery func(query sq.SelectBuilder) sq.SelectBuilder, eq interface{}) sq.SelectBuilder {
	return toQuery(query).Where(eq)
}
