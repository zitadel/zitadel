package query

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	quotaPeriodsTable = table{
		name:          projection.QuotaPeriodsProjectionTable,
		instanceIDCol: projection.QuotaColumnInstanceID,
	}
	QuotaPeriodColumnInstanceID = Column{
		name:  projection.QuotaPeriodColumnInstanceID,
		table: quotaPeriodsTable,
	}
	QuotaPeriodColumnUnit = Column{
		name:  projection.QuotaPeriodColumnUnit,
		table: quotaPeriodsTable,
	}
	QuotaPeriodColumnStart = Column{
		name:  projection.QuotaPeriodColumnStart,
		table: quotaPeriodsTable,
	}
	QuotaPeriodColumnUsage = Column{
		name:  projection.QuotaPeriodColumnUsage,
		table: quotaPeriodsTable,
	}
)

func (q *Queries) GetRemainingQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit) (remaining *uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	stmt, scan := prepareRemainingQuotaUsageQuery(ctx, q.client)
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				QuotaPeriodColumnInstanceID.identifier(): instanceID,
				QuotaPeriodColumnUnit.identifier():       unit,
				QuotaColumnLimit.identifier():            true,
			},
			sq.Expr("age(" + QuotaPeriodColumnStart.identifier() + ") < " + QuotaColumnInterval.identifier()),
			sq.Expr(QuotaPeriodColumnStart.identifier() + " <= now()"),
			sq.Expr(QuotaPeriodColumnStart.identifier() + " >= " + QuotaColumnFrom.identifier()),
		}).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-FSA3g", "Errors.Query.SQLStatement")
	}
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		remaining, err = scan(row)
		return err
	}, query, args...)
	if zerrors.IsNotFound(err) {
		return nil, nil
	}
	return remaining, err
}

func prepareRemainingQuotaUsageQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*uint64, error)) {
	return sq.
			Select(
				"greatest(0, " + QuotaColumnAmount.identifier() + "-" + QuotaPeriodColumnUsage.identifier() + ")",
			).
			From(quotaPeriodsTable.identifier()).
			Join(join(QuotaColumnUnit, QuotaPeriodColumnUnit) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*uint64, error) {
			remaining := new(uint64)
			err := row.Scan(remaining)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-quiowi2", "Errors.Internal")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-81j1jn2", "Errors.Internal")
			}
			return remaining, nil
		}
}
