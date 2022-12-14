package projection

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/zitadel/logging"

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

	QuotaNotificationIdCol                = "id"
	QuotaNotificationInstanceIDCol        = "instance_id"
	QuotaNotificationUnitCol              = "quota_unit"
	QuotaNotificationCallURLCol           = "call_url"
	QuotaNotificationPercentCol           = "percent"
	QuotaNotificationRepeatCol            = "repeat"
	QuotaNotificationLastCallDateCol      = "last_call_date"
	QuotaNotificationLastCallThresholdCol = "last_call_threshold"
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

type notificationFailed struct {
	err          error
	notification *quotaNotification
}

// Report calls notification hooks if necessary and returns if usage should be limited
func (q *Quota) Report(ctx context.Context, es *eventstore.Eventstore, used uint64) bool {

	doLimit := q.Limit && int64(used) > q.Amount

	dueNotifications, getNotificationsErr := getDueInstanceQuotaNotifications(ctx, q, used)
	if getNotificationsErr != nil {
		// TODO: log warning
		return doLimit
	}

	var errs []*notificationFailed
	for _, notification := range dueNotifications {

		alreadyNotified, alreadyNotifiedErr := isAlreadNotified(ctx, es, notification, q.PeriodStart)
		if alreadyNotifiedErr != nil {
			errs = append(errs, &notificationFailed{err: alreadyNotifiedErr, notification: notification})
			continue
		}

		if alreadyNotified {
			logging.Debugf(
				"quota notification with ID %s and threshold %d was already notified in this period",
				notification.notifiedEvent.ID,
				notification.notifiedEvent.Threshold,
			)
			continue
		}

		if notifyErr := notify(ctx, notification); notifyErr != nil {
			errs = append(errs, &notificationFailed{err: notifyErr, notification: notification})
			continue
		}

		if _, pushErr := es.Push(ctx, notification.notifiedEvent); pushErr != nil {
			errs = append(errs, &notificationFailed{err: pushErr, notification: notification})
			continue
		}
	}

	if severeErr := emitFailedEvents(ctx, q.client, errs); severeErr != nil {
		panic(severeErr)
	}

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
	notifiedEvent *instance.QuotaNotifiedEvent
	callUrl       string
}

func getDueInstanceQuotaNotifications(ctx context.Context, quota *Quota, usedAbs uint64) ([]*quotaNotification, error) {

	usedRel := int64(math.Floor(float64(usedAbs*100) / float64(quota.Amount)))

	thresholdExpr := fmt.Sprintf("%d - %d %% %s", usedRel, usedRel, QuotaNotificationPercentCol)
	// TODO: Is it possible to reuse the scalar expression in the where clause somehow?
	stmt, args, err := squirrel.Select(QuotaNotificationIdCol, QuotaNotificationCallURLCol, fmt.Sprintf("%s as threshold", thresholdExpr)).
		From(fmt.Sprintf("%s_%s  AS OF SYSTEM TIME '-10s'", QuotaTable, QuotaNotificationsTableSuffix)).
		Where(squirrel.And{
			squirrel.Eq{
				QuotaNotificationInstanceIDCol: quota.instanceId,
				QuotaNotificationUnitCol:       quota.unit,
			},
			squirrel.Lt{
				QuotaNotificationPercentCol: usedRel,
			},
			squirrel.Or{
				squirrel.Eq{
					QuotaNotificationLastCallDateCol: nil,
				},
				squirrel.Lt{
					QuotaNotificationLastCallDateCol: quota.PeriodStart,
				},
				squirrel.Or{
					squirrel.And{
						squirrel.Eq{
							QuotaNotificationRepeatCol: true,
						},
						squirrel.Or{
							squirrel.Eq{
								QuotaNotificationLastCallThresholdCol: nil,
							},
							squirrel.Expr(fmt.Sprintf("%s > %s", thresholdExpr, QuotaNotificationLastCallThresholdCol)),
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

	var notifications []*quotaNotification
	rows, err := quota.client.QueryContext(ctx, stmt, args...)
	for rows.Next() {
		row := struct {
			id        string
			callUrl   string
			threshold uint64
		}{}
		if rows.Scan(&row.id, &row.callUrl, &row.threshold); err != nil {
			return nil, caos_errors.ThrowInternal(err, "QUOTA-pBPrM", "Errors.Quota.ScanFailed")
		}
		notifications = append(notifications, &quotaNotification{
			notifiedEvent: instance.NewQuotaNotifiedEvent(
				ctx,
				&instance.NewAggregate(quota.instanceId).Aggregate,
				quota.unit,
				row.id,
				row.threshold,
				usedAbs,
			),
			callUrl: row.callUrl,
		})
	}
	return notifications, nil
}

func isAlreadNotified(ctx context.Context, es *eventstore.Eventstore, notification *quotaNotification, periodStart time.Time) (bool, error) {

	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(notification.notifiedEvent.Aggregate().InstanceID).
			AddQuery().
			AggregateTypes(instance.AggregateType).
			AggregateIDs(notification.notifiedEvent.Aggregate().ID).
			SequenceGreater(notification.notifiedEvent.Sequence()).
			EventTypes(quota.NotifiedEventType).
			CreationDateAfter(periodStart).
			EventData(map[string]interface{}{
				"id":        notification.notifiedEvent.ID,
				"threshold": notification.notifiedEvent.Threshold,
			}).
			Builder(),
	)
	return len(events) > 0, err
}

func notify(ctx context.Context, notification *quotaNotification) error {

	payload, err := json.Marshal(notification.notifiedEvent)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notification.callUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("calling url %s returned %s", notification.callUrl, resp.Status)
	}

	return nil
}

func emitFailedEvents(ctx context.Context, client *sql.DB, errs []*notificationFailed) error {
	for _, err := range errs {
		// TODO: Implement
		logging.WithError(err.err).Warn("creating failed event for notification failures not implemented, yet")
	}
	return nil
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
				crdb.NewColumn(QuotaNotificationLastCallDateCol, crdb.ColumnTypeTimestamp, crdb.Nullable()),
				crdb.NewColumn(QuotaNotificationLastCallThresholdCol, crdb.ColumnTypeInt64, crdb.Nullable()),
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
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.QuotaNotifiedEventType,
					Reduce: reduceQuotaNotified,
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
				handler.NewCol(QuotaNotificationLastCallDateCol, nil),
				handler.NewCol(QuotaNotificationLastCallThresholdCol, nil),
			},
			crdb.WithTableSuffix(QuotaNotificationsTableSuffix),
		))
	}

	return crdb.NewMultiStatement(e, execFuncs...), nil
}

func reduceQuotaNotified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.QuotaNotifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-wmQPo", "reduce.wrong.event.type% s", quota.NotifiedEventType)
	}

	return crdb.NewUpdateStatement(event, []handler.Column{
		handler.NewCol(QuotaNotificationLastCallDateCol, e.CreationDate()),
		handler.NewCol(QuotaNotificationLastCallThresholdCol, e.Threshold),
	}, []handler.Condition{
		handler.NewCond(QuotaNotificationInstanceIDCol, e.Aggregate().InstanceID),
		handler.NewCond(QuotaNotificationUnitCol, e.Unit),
		handler.NewCond(QuotaNotificationIdCol, e.ID),
	}, crdb.WithTableSuffix(QuotaNotificationsTableSuffix)), nil
}
