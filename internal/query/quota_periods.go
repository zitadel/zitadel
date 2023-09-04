package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	zitadel_errors "github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/api/call"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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

func (q *Queries) GetQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit, periodStart time.Time) (usage uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, scan := prepareQuotaUsageQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.Eq{
			QuotaPeriodColumnInstanceID.identifier(): instanceID,
			QuotaPeriodColumnUnit.identifier():       unit,
			QuotaPeriodColumnStart.identifier():      periodStart,
		}).
		ToSql()
	if err != nil {
		return 0, zitadel_errors.ThrowInternal(err, "QUERY-mEZX2", "Errors.Query.SQLStatement")
	}
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		usage, err = scan(row)
		return err
	}, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return usage, err
}

func (q *Queries) GetRemainingQuotaUsage(ctx context.Context, instanceID string, unit quota.Unit) (remaining *uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, scan := prepareRemainingQuotaUsageQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				QuotaPeriodColumnInstanceID.identifier(): instanceID,
				QuotaPeriodColumnUnit.identifier():       unit,
				QuotaColumnLimit.identifier():            true,
			},
			//sq.Lt{
			//	//"age(" + QuotaPeriodColumnStart.identifier() + ")": QuotaColumnInterval.identifier,
			//	//:                "now()",
			//},
			sq.Expr("age(" + QuotaPeriodColumnStart.identifier() + ") < " + QuotaColumnInterval.identifier()),
			sq.Expr(QuotaPeriodColumnStart.identifier() + " < now()"),
		}).
		ToSql()
	if err != nil {
		return nil, zitadel_errors.ThrowInternal(err, "QUERY-FSA3g", "Errors.Query.SQLStatement")
	}
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		remaining, err = scan(row)
		return err
	}, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return remaining, err
}

func prepareQuotaUsageQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (uint64, error)) {
	return sq.
			Select(QuotaPeriodColumnUsage.identifier()).
			From(quotaPeriodsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (usage uint64, err error) {
			return usage, row.Scan(&usage)
		}
}

func prepareRemainingQuotaUsageQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*uint64, error)) {
	return sq.
			Select(
				"greatest(0, (" + QuotaPeriodColumnUsage.identifier() + "-" + QuotaColumnAmount.identifier() + "))",
			).
			From(quotaPeriodsTable.identifier()).
			Join(join(QuotaColumnUnit, QuotaPeriodColumnUnit) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (usage *uint64, err error) {
			return usage, row.Scan(&usage)
		}
}
