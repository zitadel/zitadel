package query

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/zitadel/zitadel/internal/repository/instance"

	"github.com/zitadel/zitadel/internal/query/projection"

	"github.com/Masterminds/squirrel"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

func GetInstanceQuota(ctx context.Context, client *sql.DB, instanceID string, unit quota.Unit) (*Quota, error) {

	stmt, args, err := squirrel.Select(projection.QuotaAmountCol, projection.QuotaLimitCol, projection.QuotaFromCol, projection.QuotaIntervalCol).
		From(projection.QuotaTable /* + " AS OF SYSTEM TIME '-20s'"*/). // TODO: Incomment
		Where(squirrel.Eq{
			projection.QuotaInstanceIDCol: instanceID,
			projection.QuotaUnitCol:       unit,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-V9Sde", "Errors.Internal")
	}

	quota := Quota{
		InstanceId: instanceID,
		Unit:       unit,
	}
	if err = client.
		QueryRowContext(ctx, stmt, args...).
		Scan(&quota.Amount, &quota.Limit, &quota.from, &quota.Interval); err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Quota.ScanFailed")
	}
	quota.refreshPeriod()
	return &quota, nil
}

type Quota struct {
	Amount      int64
	Limit       bool
	InstanceId  string
	Unit        quota.Unit
	from        time.Time
	Interval    time.Duration
	PeriodStart time.Time
	PeriodEnd   time.Time
}

func (q *Quota) refreshPeriod() {
	q.PeriodStart = pushFrom(q.from, q.Interval, time.Now())
	q.PeriodEnd = q.PeriodStart.Add(q.Interval)
}

func pushFrom(from time.Time, interval time.Duration, now time.Time) time.Time {
	next := from.Add(interval)
	if next.After(now) {
		return from
	}
	return pushFrom(next, interval, now)
}

func GetDueInstanceQuotaNotifications(ctx context.Context, dbClient *sql.DB, quota *Quota, usedAbs uint64) ([]*instance.QuotaNotifiedEvent, error) {

	usedRel := int64(math.Floor(float64(usedAbs*100) / float64(quota.Amount)))

	thresholdExpr := fmt.Sprintf("%d - %d %% %s", usedRel, usedRel, projection.QuotaNotificationPercentCol)
	// TODO: Is it possible to reuse the scalar expression in the where clause somehow?
	stmt, args, err := squirrel.Select(projection.QuotaNotificationIdCol, projection.QuotaNotificationCallURLCol, fmt.Sprintf("%s as threshold", thresholdExpr)).
		From(fmt.Sprintf("%s_%s  AS OF SYSTEM TIME '-10s'", projection.QuotaTable, projection.QuotaNotificationsTableSuffix)).
		Where(squirrel.And{
			squirrel.Eq{
				projection.QuotaNotificationInstanceIDCol: quota.InstanceId,
				projection.QuotaNotificationUnitCol:       quota.Unit,
			},
			squirrel.Lt{
				projection.QuotaNotificationPercentCol: usedRel,
			},
			squirrel.Or{
				squirrel.Eq{
					projection.QuotaNotificationLastCallDateCol: nil,
				},
				squirrel.Lt{
					projection.QuotaNotificationLastCallDateCol: quota.PeriodStart,
				},
				squirrel.Or{
					squirrel.And{
						squirrel.Eq{
							projection.QuotaNotificationRepeatCol: true,
						},
						squirrel.Or{
							squirrel.Eq{
								projection.QuotaNotificationLastCallThresholdCol: nil,
							},
							squirrel.Expr(fmt.Sprintf("%s > %s", thresholdExpr, projection.QuotaNotificationLastCallThresholdCol)),
						},
					},
				},
			},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-V9Sde", "Errors.Internal")
	}

	var notifications []*instance.QuotaNotifiedEvent
	rows, err := dbClient.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-SV9LW", "Errors.Quota.QueryFailed")
	}
	for rows.Next() {
		row := struct {
			id        string
			callUrl   string
			threshold uint64
		}{}
		if rows.Scan(&row.id, &row.callUrl, &row.threshold); err != nil {
			return nil, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Quota.ScanFailed")
		}
		notifications = append(notifications, instance.NewQuotaNotifiedEvent(
			ctx,
			&instance.NewAggregate(quota.InstanceId).Aggregate,
			quota.Unit,
			row.id,
			row.callUrl,
			quota.PeriodStart,
			row.threshold,
			usedAbs,
		))
	}
	return notifications, nil
}
