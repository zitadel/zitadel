package projection

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/errors"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	QuotaTable                    = "projections.instance_quotas"
	QuotaNotificationsTableSuffix = "notifications"

	QuotaCreationDateCol  = "creation_date"
	QuotaChangeDateCol    = "change_date"
	QuotaResourceOwnerCol = "resource_owner"
	QuotaInstanceIDCol    = "instance_id"
	QuotaSequenceCol      = "sequence"
	QuotaUnitCol          = "unit"
	QuotaFromCol          = "active_from"
	QuotaIntervalCol      = "interval"
	QuotaAmountCol        = "amount"
	QuotaLimitCol         = "do_limit"

	QuotaNotificationIdCol              = "id"
	QuotaNotificationInstanceIDCol      = "instance_id"
	QuotaNotificationUnitCol            = "quota_unit"
	QuotaNotificationCallURLCol         = "call_url"
	QuotaNotificationPercentCol         = "percent"
	QuotaNotificationRepeatCol          = "repeat"
	QuotaNotificationLastCallDateCol    = "last_call_date"
	QuotaNotificationLastCallPercentCol = "last_call_percent"
)

type Quota struct {
	Amount      int64
	Limit       bool
	client      *sql.DB
	instanceId  string
	unit        quota.Unit
	from        time.Time
	interval    time.Duration
	PeriodStart time.Time
	PeriodEnd   time.Time
}

// Report calls notification hooks if necessary and returns if usage should be limited
func (q *Quota) Report(ctx context.Context, used uint64) bool {

	var errs []error
	dueNotifications, getNotificationsErr := getDueInstanceQuotaNotifications(ctx, client, q, int64(used))
	if getNotificationsErr != nil {
		errs = append(errs, getNotificationsErr)
	}
	for _, notification := range dueNotifications {
		alreadyNotified, alreadyNotifiedErr := isNotified(ctx, client, q, notification)
		if alreadyNotifiedErr != nil {
			errs = append(errs, alreadyNotifiedErr)
			continue
		}
		if alreadyNotified {
			continue
		}

		if notifyErr := notify(ctx, q, notification); notifyErr != nil {
			errs = append(errs, notifyErr)
			continue
		}

		if emitMotifiedErr := emitNotifiedEvent(ctx, client, q, notification); emitMotifiedErr != nil {
			errs = append(errs, emitMotifiedErr)
			continue
		}
	}

	if severeErr := emitFailedEvents(ctx, client, errs); severeErr != nil {
		panic(severeErr)
	}

	doLimit := q.Limit && int64(used) > q.Amount
	return doLimit
}

func (q *Quota) refreshPeriod() {
	q.PeriodStart = pushFrom(q.from, q.interval, time.Now())
	q.PeriodEnd = q.PeriodStart.Add(q.interval)
}

func pushFrom(from time.Time, interval time.Duration, now time.Time) time.Time {
	next := from.Add(interval)
	if next.After(now) {
		return from
	}
	return pushFrom(next, interval, now)
}

type quotaNotification struct {
}

func getDueInstanceQuotaNotifications(ctx context.Context, client *sql.DB, q *Quota, used int64) ([]*quotaNotification, error) {

	usagePercent := int64(math.Floor(float64(used*100) / float64(q.Amount)))

	stmt, args, err := squirrel.Select(QuotaAmountCol, QuotaLimitCol, QuotaFromCol, QuotaIntervalCol).
		From(fmt.Sprintf("%s_%s  AS OF SYSTEM TIME '-10s'", QuotaTable, QuotaNotificationsTableSuffix)).
		Where(squirrel.And{
			squirrel.Eq{
				QuotaNotificationInstanceIDCol: q.instanceId,
				QuotaNotificationUnitCol:       q.unit,
			},
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{
						QuotaNotificationRepeatCol: false,
					},
					squirrel.Lt{
						QuotaNotificationPercentCol:      usagePercent,
						QuotaNotificationLastCallDateCol: q.PeriodStart,
					},
				},
				squirrel.And{
					squirrel.Eq{
						QuotaNotificationRepeatCol: true,
					},
					squirrel.Expr(),
					squirrel.Lt{
						QuotaNotificationPercentCol:       usagePercent,
						QuotaNotificationLastCallUsageCol: q.PeriodStart,
					},
				},
			},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-V9Sde", "Errors.Internal")
	}

	notififcation := quotaNotification{}
	rows, err := client.QueryContext(ctx, stmt, args...)
	for rows.Next() {
		if rows.Scan(&notififcation.Amount, &notififcation.Limit, &notififcation.from, &notififcation.interval); err != nil {

		}
	}
	return nil, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Quota.ScanFailed")

	return quota.refreshPeriod(), nil
}

func isNotified(ctx context.Context, client *sql.DB, quota *Quota, notification *quotaNotification) (bool, error) {
	return false, errors.ThrowError(nil, "", "not implemented")
}

func notify(ctx context.Context, quota *Quota, notification *quotaNotification) error {
	return errors.ThrowError(nil, "", "not implemented")
}

func emitNotifiedEvent(ctx context.Context, client *sql.DB, quota *Quota, notification *quotaNotification) error {
	return errors.ThrowError(nil, "", "not implemented")
}

func emitFailedEvents(ctx context.Context, client *sql.DB, errs []error) error {
	return errors.ThrowError(nil, "", "not implemented")
}

func GetInstanceQuota(ctx context.Context, client *sql.DB, instanceID string, unit quota.Unit) (*Quota, error) {

	stmt, args, err := squirrel.Select(QuotaAmountCol, QuotaLimitCol, QuotaFromCol, QuotaIntervalCol).
		From(QuotaTable + " AS OF SYSTEM TIME '-20s'").
		Where(squirrel.Eq{
			QuotaInstanceIDCol: instanceID,
			QuotaUnitCol:       unit,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-V9Sde", "Errors.Internal")
	}

	quota := Quota{
		client:     client,
		instanceId: instanceID,
		unit:       unit,
	}
	if err = client.
		QueryRowContext(ctx, stmt, args...).
		Scan(&quota.Amount, &quota.Limit, &quota.from, &quota.interval); err != nil {
		return nil, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Quota.ScanFailed")
	}
	quota.refreshPeriod()
	return &quota, nil
}

// TODO: Why not return *StatementHandler?
type quotaProjection struct {
	crdb.StatementHandler
}

func newQuotaProjection(ctx context.Context, esHandlerConfig crdb.StatementHandlerConfig) *quotaProjection {
	p := new(quotaProjection)
	esHandlerConfig.ProjectionName = QuotaTable
	esHandlerConfig.Reducers = esReducers()
	esHandlerConfig.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaCreationDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaChangeDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaResourceOwnerCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaSequenceCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaUnitCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaFromCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaIntervalCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaAmountCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitCol, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotaInstanceIDCol, QuotaUnitCol),
			crdb.WithIndex(crdb.NewIndex("quotas_ro_idx", []string{QuotaResourceOwnerCol})),
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaNotificationIdCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationUnitCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaNotificationCallURLCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationPercentCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaNotificationRepeatCol, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotaNotificationInstanceIDCol, QuotaNotificationUnitCol, QuotaNotificationIdCol),
			QuotaNotificationsTableSuffix,
			crdb.WithForeignKey(
				crdb.NewForeignKey(
					"fk_instance_quotas_notifications_ref_instance_quotas",
					[]string{QuotaNotificationInstanceIDCol, QuotaNotificationUnitCol},
					[]string{QuotaInstanceIDCol, QuotaUnitCol},
				),
			),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, esHandlerConfig)
	return p
}

func esReducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.QuotaAddedEventType,
					Reduce: reduceQuotaAdded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(QuotaInstanceIDCol),
				},
			},
		},
	}
}

func reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.QuotaAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff21", "reduce.wrong.event.type% s", quota.AddedEventType)
	}

	execFuncs := []func(eventstore.Event) crdb.Exec{
		crdb.AddCreateStatement([]handler.Column{
			handler.NewCol(QuotaCreationDateCol, e.CreationDate()),
			handler.NewCol(QuotaChangeDateCol, e.CreationDate()),
			handler.NewCol(QuotaResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(QuotaInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(QuotaSequenceCol, e.Sequence()),
			handler.NewCol(QuotaUnitCol, e.Unit),
			handler.NewCol(QuotaFromCol, e.From),
			handler.NewCol(QuotaIntervalCol, e.Interval),
			handler.NewCol(QuotaAmountCol, e.Amount),
			handler.NewCol(QuotaLimitCol, e.Limit),
		}),
	}

	for _, notification := range e.Notifications {

		execFuncs = append(execFuncs, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(QuotaNotificationInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(QuotaNotificationUnitCol, e.Unit),
				handler.NewCol(QuotaNotificationIdCol, notification.ID),
				handler.NewCol(QuotaNotificationPercentCol, notification.Percent),
				handler.NewCol(QuotaNotificationRepeatCol, notification.Repeat),
				handler.NewCol(QuotaNotificationCallURLCol, notification.CallURL),
			},
			crdb.WithTableSuffix(QuotaNotificationsTableSuffix),
		))
	}

	return crdb.NewMultiStatement(e, execFuncs...), nil
}

type nextNotification int64

const (
	unknown nextNotification = -1
	due     nextNotification = -2
	done    nextNotification = -3
)

// TODO: think
func nextNotificationCall(amount, percent, used, currentThreshold int64, repeat bool) int64 {

	if nextNotification(currentThreshold) == due || nextNotification(currentThreshold) == done {
		return currentThreshold
	}

	nextThreshold := int64(math.Floor(float64(amount/100))) * percent
	if used < nextThreshold && currentThreshold != -1 {
		return nextThreshold
	}

	if repeat {
		return nextNotificationCall(amount, percent*2, used, currentThreshold, repeat)
	}
	return -1
}
