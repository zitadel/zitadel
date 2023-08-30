package query

import (
	"context"
	"database/sql"
	errs "errors"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/zitadel/internal/api/call"
	zitadel_errors "github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/query/projection"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
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
	query, scan := prepareQuotaNotificationsQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				QuotaNotificationColumnInstanceID.identifier(): instanceID,
				QuotaNotificationColumnUnit.identifier():       unit,
			},
			sq.Or{
				// If the relative usage is greater than the next due threshold, it's clear we can notify,
				// because min next threshold is equal to the percent field or null (LtOrEq to Null comparison is null -> false)
				sq.LtOrEq{QuotaNotificationColumnNextDueThreshold.identifier(): usedRel},
				sq.And{
					// In case we haven't seen a due notification for this quota period, we compare against the configured percent
					sq.Or{
						sq.Eq{QuotaNotificationColumnLatestDuePeriodStart.identifier(): nil},
						sq.Lt{QuotaNotificationColumnLatestDuePeriodStart.identifier(): periodStart},
					},
					sq.LtOrEq{QuotaNotificationColumnPercent.identifier(): usedRel},
				},
			},
		},
	).ToSql()
	if err != nil {
		return nil, zitadel_errors.ThrowInternal(err, "QUERY-XmYn9", "Errors.Query.SQLStatement")
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
		reachedThreshold := notification.NextDueThreshold
		if reachedThreshold == nil {
			reachedThreshold = &notification.Percent
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
				*reachedThreshold,
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
	NextDueThreshold *uint16
}

type QuotaNotifications struct {
	SearchResponse
	Configs []*QuotaNotification
}

func prepareQuotaNotificationsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*QuotaNotifications, error)) {
	return sq.Select(
			QuotaNotificationColumnID.identifier(),
			QuotaNotificationColumnCallURL.identifier(),
			QuotaNotificationColumnPercent.identifier(),
			QuotaNotificationColumnRepeat.identifier(),
			QuotaNotificationColumnNextDueThreshold.identifier(),
		).
			From(quotaNotificationsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(rows *sql.Rows) (*QuotaNotifications, error) {
			cfgs := &QuotaNotifications{Configs: []*QuotaNotification{}}
			for rows.Next() {
				cfg := new(QuotaNotification)
				var nextDueThreshold sql.NullInt16
				err := rows.Scan(&cfg.ID, &cfg.CallURL, &cfg.Percent, &cfg.Repeat, &nextDueThreshold)
				if err != nil {
					if errs.Is(err, sql.ErrNoRows) {
						return nil, zitadel_errors.ThrowNotFound(err, "QUERY-bbqWb", "Errors.QuotaNotification.NotExisting")
					}
					return nil, zitadel_errors.ThrowInternal(err, "QUERY-8copS", "Errors.Internal")
				}
				if nextDueThreshold.Valid {
					n := uint16(nextDueThreshold.Int16)
					cfg.NextDueThreshold = &n
				}
				cfgs.Configs = append(cfgs.Configs, cfg)
			}
			return cfgs, nil
		}
}

func (q *Queries) getQuotaNotificationsReadModel(ctx context.Context, aggregate eventstore.Aggregate, periodStart time.Time) (*quotaNotificationsReadModel, error) {
	wm := newQuotaNotificationsReadModel(aggregate.ID, aggregate.InstanceID, aggregate.ResourceOwner, periodStart)
	return wm, q.eventstore.FilterToQueryReducer(ctx, wm)
}
