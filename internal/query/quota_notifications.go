package query

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	quotaNotificationsTable = table{
		name:          projection.QuotaNotificationsTable,
		instanceIDCol: projection.QuotaNotificationColumnInstanceID,
	}
	QuotaNotificationColumnInstanceID = Column{
		name:  projection.QuotaNotificationColumnInstanceID,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnUnit = Column{
		name:  projection.QuotaNotificationColumnUnit,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnID = Column{
		name:  projection.QuotaNotificationColumnID,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnCallURL = Column{
		name:  projection.QuotaNotificationColumnCallURL,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnPercent = Column{
		name:  projection.QuotaNotificationColumnPercent,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnRepeat = Column{
		name:  projection.QuotaNotificationColumnRepeat,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnLatestDuePeriodStart = Column{
		name:  projection.QuotaNotificationColumnLatestDuePeriodStart,
		table: quotaNotificationsTable,
	}
	QuotaNotificationColumnNextDueThreshold = Column{
		name:  projection.QuotaNotificationColumnNextDueThreshold,
		table: quotaNotificationsTable,
	}
)

func (q *Queries) GetDueQuotaNotifications(ctx context.Context, instanceID string, unit quota.Unit, qu *Quota, periodStart time.Time, usedAbs uint64) (dueNotifications []*quota.NotificationDueEvent, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	usedRel := uint16(math.Floor(float64(usedAbs*100) / float64(qu.Amount)))
	query, scan := prepareQuotaNotificationsQuery()
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				QuotaNotificationColumnInstanceID.identifier(): instanceID,
				QuotaNotificationColumnUnit.identifier():       unit,
			},
			sq.Or{
				// If the relative usage is greater than the next due threshold in the current period, it's clear we can notify
				sq.And{
					sq.Eq{QuotaNotificationColumnLatestDuePeriodStart.identifier(): periodStart},
					sq.LtOrEq{QuotaNotificationColumnNextDueThreshold.identifier(): usedRel},
				},
				// In case we haven't seen a due notification for this quota period, we compare against the configured percent
				sq.And{
					sq.Or{
						sq.Expr(QuotaNotificationColumnLatestDuePeriodStart.identifier() + " IS NULL"),
						sq.NotEq{QuotaNotificationColumnLatestDuePeriodStart.identifier(): periodStart},
					},
					sq.LtOrEq{QuotaNotificationColumnPercent.identifier(): usedRel},
				},
			},
		},
	).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-XmYn9", "Errors.Query.SQLStatement")
	}
	var notifications *QuotaNotifications
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		notifications, err = scan(rows)
		return err
	}, stmt, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications.Configs {
		reachedThreshold := calculateThreshold(usedRel, notification.Percent)
		if !notification.Repeat && notification.Percent < reachedThreshold {
			continue
		}
		dueNotifications = append(
			dueNotifications,
			quota.NewNotificationDueEvent(
				ctx,
				&quota.NewAggregate(qu.ID, instanceID).Aggregate,
				unit,
				notification.ID,
				notification.CallURL,
				periodStart,
				reachedThreshold,
				usedAbs,
			),
		)
	}
	return dueNotifications, nil
}

type QuotaNotification struct {
	ID               string
	CallURL          string
	Percent          uint16
	Repeat           bool
	NextDueThreshold uint16
}

type QuotaNotifications struct {
	SearchResponse
	Configs []*QuotaNotification
}

// calculateThreshold calculates the nearest reached threshold.
// It makes sure that the percent configured on the notification is calculated within the "current" 100%,
// e.g. when configuring 80%, the thresholds are 80, 180, 280, ...
// so 170% use is always 70% of the current bucket, with the above config, the reached threshold would be 80.
func calculateThreshold(usedRel, notificationPercent uint16) uint16 {
	// check how many times we reached 100%
	times := math.Floor(float64(usedRel) / 100)
	// check how many times we reached the percent configured with the "current" 100%
	percent := math.Floor(float64(usedRel%100) / float64(notificationPercent))
	// If neither is reached, directly return 0.
	// This way we don't end up in some wrong uint16 range in the calculation below.
	if times == 0 && percent == 0 {
		return 0
	}
	return uint16(times+percent-1)*100 + notificationPercent
}

func prepareQuotaNotificationsQuery() (sq.SelectBuilder, func(*sql.Rows) (*QuotaNotifications, error)) {
	return sq.Select(
			QuotaNotificationColumnID.identifier(),
			QuotaNotificationColumnCallURL.identifier(),
			QuotaNotificationColumnPercent.identifier(),
			QuotaNotificationColumnRepeat.identifier(),
			QuotaNotificationColumnNextDueThreshold.identifier(),
		).
			From(quotaNotificationsTable.identifier()).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) (*QuotaNotifications, error) {
			cfgs := &QuotaNotifications{Configs: []*QuotaNotification{}}
			for rows.Next() {
				cfg := new(QuotaNotification)
				var nextDueThreshold sql.NullInt16
				err := rows.Scan(&cfg.ID, &cfg.CallURL, &cfg.Percent, &cfg.Repeat, &nextDueThreshold)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return nil, zerrors.ThrowNotFound(err, "QUERY-bbqWb", "Errors.QuotaNotification.NotExisting")
					}
					return nil, zerrors.ThrowInternal(err, "QUERY-8copS", "Errors.Internal")
				}
				if nextDueThreshold.Valid {
					cfg.NextDueThreshold = uint16(nextDueThreshold.Int16)
				}
				cfgs.Configs = append(cfgs.Configs, cfg)
			}
			return cfgs, nil
		}
}
