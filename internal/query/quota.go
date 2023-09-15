package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	quotasTable = table{
		name:          projection.QuotasProjectionTable,
		instanceIDCol: projection.QuotaColumnInstanceID,
	}
	QuotaColumnID = Column{
		name:  projection.QuotaColumnID,
		table: quotasTable,
	}
	QuotaColumnInstanceID = Column{
		name:  projection.QuotaColumnInstanceID,
		table: quotasTable,
	}
	QuotaColumnUnit = Column{
		name:  projection.QuotaColumnUnit,
		table: quotasTable,
	}
	QuotaColumnAmount = Column{
		name:  projection.QuotaColumnAmount,
		table: quotasTable,
	}
	QuotaColumnLimit = Column{
		name:  projection.QuotaColumnLimit,
		table: quotasTable,
	}
	QuotaColumnInterval = Column{
		name:  projection.QuotaColumnInterval,
		table: quotasTable,
	}
	QuotaColumnFrom = Column{
		name:  projection.QuotaColumnFrom,
		table: quotasTable,
	}
)

type Quota struct {
	ID                 string
	From               time.Time
	ResetInterval      time.Duration
	Amount             uint64
	Limit              bool
	CurrentPeriodStart time.Time
}

func (q *Queries) GetQuota(ctx context.Context, instanceID string, unit quota.Unit) (qu *Quota, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, scan := prepareQuotaQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.Eq{
			QuotaColumnInstanceID.identifier(): instanceID,
			QuotaColumnUnit.identifier():       unit,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-XmYn9", "Errors.Query.SQLStatement")
	}
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		qu, err = scan(row)
		return err
	}, stmt, args...)
	return qu, err
}

func prepareQuotaQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*Quota, error)) {
	return sq.
			Select(
				QuotaColumnID.identifier(),
				QuotaColumnFrom.identifier(),
				QuotaColumnInterval.identifier(),
				QuotaColumnAmount.identifier(),
				QuotaColumnLimit.identifier(),
				"now()",
			).
			From(quotasTable.identifier()).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*Quota, error) {
			q := new(Quota)
			var interval database.Duration
			var now time.Time
			err := row.Scan(&q.ID, &q.From, &interval, &q.Amount, &q.Limit, &now)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-rDTM6", "Errors.Quota.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-LqySK", "Errors.Internal")
			}
			q.ResetInterval = time.Duration(interval)
			q.CurrentPeriodStart = pushPeriodStart(q.From, q.ResetInterval, now)
			return q, nil
		}
}

func pushPeriodStart(from time.Time, interval time.Duration, now time.Time) time.Time {
	if now.IsZero() {
		now = time.Now()
	}
	for {
		next := from.Add(interval)
		if next.After(now) {
			return from
		}
		from = next
	}
}
