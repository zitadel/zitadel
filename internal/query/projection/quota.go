package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	QuotasProjectionTable       = "projections.quotas"
	QuotaPeriodsProjectionTable = QuotasProjectionTable + "_" + quotaPeriodsTableSuffix
	QuotaNotificationsTable     = QuotasProjectionTable + "_" + quotaNotificationsTableSuffix

	QuotaColumnID         = "id"
	QuotaColumnInstanceID = "instance_id"
	QuotaColumnUnit       = "unit"
	QuotaColumnAmount     = "amount"
	QuotaColumnFrom       = "from_anchor"
	QuotaColumnInterval   = "interval"
	QuotaColumnLimit      = "limit_usage"

	quotaPeriodsTableSuffix     = "periods"
	QuotaPeriodColumnInstanceID = "instance_id"
	QuotaPeriodColumnUnit       = "unit"
	QuotaPeriodColumnStart      = "start"
	QuotaPeriodColumnUsage      = "usage"

	quotaNotificationsTableSuffix               = "notifications"
	QuotaNotificationColumnInstanceID           = "instance_id"
	QuotaNotificationColumnUnit                 = "unit"
	QuotaNotificationColumnID                   = "id"
	QuotaNotificationColumnCallURL              = "call_url"
	QuotaNotificationColumnPercent              = "percent"
	QuotaNotificationColumnRepeat               = "repeat"
	QuotaNotificationColumnLatestDuePeriodStart = "latest_due_period_start"
	QuotaNotificationColumnNextDueThreshold     = "next_due_threshold"
)

const (
	incrementQuotaStatement = `INSERT INTO projections.quotas_periods` +
		` (instance_id, unit, start, usage)` +
		` VALUES ($1, $2, $3, $4) ON CONFLICT (instance_id, unit, start)` +
		` DO UPDATE SET usage = projections.quotas_periods.usage + excluded.usage RETURNING usage`
)

type quotaProjection struct {
	handler *handler.Handler
	client  *database.DB
}

func newQuotaProjection(ctx context.Context, config handler.Config) *quotaProjection {
	p := &quotaProjection{
		client: config.Client,
	}
	p.handler = handler.NewHandler(ctx, &config, p)
	return p
}

func (*quotaProjection) Name() string {
	return QuotasProjectionTable
}

func (*quotaProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable(
			[]*handler.InitColumn{
				handler.NewColumn(QuotaColumnID, handler.ColumnTypeText),
				handler.NewColumn(QuotaColumnInstanceID, handler.ColumnTypeText),
				handler.NewColumn(QuotaColumnUnit, handler.ColumnTypeEnum),
				handler.NewColumn(QuotaColumnAmount, handler.ColumnTypeInt64, handler.Nullable()),
				handler.NewColumn(QuotaColumnFrom, handler.ColumnTypeTimestamp, handler.Nullable()),
				handler.NewColumn(QuotaColumnInterval, handler.ColumnTypeInterval, handler.Nullable()),
				handler.NewColumn(QuotaColumnLimit, handler.ColumnTypeBool, handler.Nullable()),
			},
			handler.NewPrimaryKey(QuotaColumnInstanceID, QuotaColumnUnit),
		),
		handler.NewSuffixedTable(
			[]*handler.InitColumn{
				handler.NewColumn(QuotaPeriodColumnInstanceID, handler.ColumnTypeText),
				handler.NewColumn(QuotaPeriodColumnUnit, handler.ColumnTypeEnum),
				handler.NewColumn(QuotaPeriodColumnStart, handler.ColumnTypeTimestamp),
				handler.NewColumn(QuotaPeriodColumnUsage, handler.ColumnTypeInt64),
			},
			handler.NewPrimaryKey(QuotaPeriodColumnInstanceID, QuotaPeriodColumnUnit, QuotaPeriodColumnStart),
			quotaPeriodsTableSuffix,
		),
		handler.NewSuffixedTable(
			[]*handler.InitColumn{
				handler.NewColumn(QuotaNotificationColumnInstanceID, handler.ColumnTypeText),
				handler.NewColumn(QuotaNotificationColumnUnit, handler.ColumnTypeEnum),
				handler.NewColumn(QuotaNotificationColumnID, handler.ColumnTypeText),
				handler.NewColumn(QuotaNotificationColumnCallURL, handler.ColumnTypeText),
				handler.NewColumn(QuotaNotificationColumnPercent, handler.ColumnTypeInt64),
				handler.NewColumn(QuotaNotificationColumnRepeat, handler.ColumnTypeBool),
				handler.NewColumn(QuotaNotificationColumnLatestDuePeriodStart, handler.ColumnTypeTimestamp, handler.Nullable()),
				handler.NewColumn(QuotaNotificationColumnNextDueThreshold, handler.ColumnTypeInt64, handler.Nullable()),
			},
			handler.NewPrimaryKey(QuotaNotificationColumnInstanceID, QuotaNotificationColumnUnit, QuotaNotificationColumnID),
			quotaNotificationsTableSuffix,
		),
	)
}

func (q *quotaProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: q.reduceInstanceRemoved,
				},
			},
		},
		{
			Aggregate: quota.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  quota.AddedEventType,
					Reduce: q.reduceQuotaSet,
				},
				{
					Event:  quota.SetEventType,
					Reduce: q.reduceQuotaSet,
				},
				{
					Event:  quota.RemovedEventType,
					Reduce: q.reduceQuotaRemoved,
				},
				{
					Event:  quota.NotificationDueEventType,
					Reduce: q.reduceQuotaNotificationDue,
				},
				{
					Event:  quota.NotifiedEventType,
					Reduce: q.reduceQuotaNotified,
				},
			},
		},
	}
}

func (q *quotaProjection) reduceQuotaNotified(event eventstore.Event) (*handler.Statement, error) {
	return handler.NewNoOpStatement(event), nil
}

func (q *quotaProjection) reduceQuotaSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.SetEvent](event)
	if err != nil {
		return nil, err
	}
	var statements []func(e eventstore.Event) handler.Exec

	// 1. Insert or update quota if the event has not only notification changes
	quotaConflictColumns := []handler.Column{
		handler.NewCol(QuotaColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(QuotaColumnUnit, e.Unit),
	}
	quotaUpdateCols := make([]handler.Column, 0, 4+1+len(quotaConflictColumns))
	if e.Limit != nil {
		quotaUpdateCols = append(quotaUpdateCols, handler.NewCol(QuotaColumnLimit, *e.Limit))
	}
	if e.Amount != nil {
		quotaUpdateCols = append(quotaUpdateCols, handler.NewCol(QuotaColumnAmount, *e.Amount))
	}
	if e.From != nil {
		quotaUpdateCols = append(quotaUpdateCols, handler.NewCol(QuotaColumnFrom, *e.From))
	}
	if e.ResetInterval != nil {
		quotaUpdateCols = append(quotaUpdateCols, handler.NewCol(QuotaColumnInterval, *e.ResetInterval))
	}
	if len(quotaUpdateCols) > 0 {
		// TODO: Add the quota ID to the primary key in a migration?
		quotaUpdateCols = append(quotaUpdateCols, handler.NewCol(QuotaColumnID, e.Aggregate().ID))
		quotaUpdateCols = append(quotaUpdateCols, quotaConflictColumns...)
		statements = append(statements, handler.AddUpsertStatement(quotaConflictColumns, quotaUpdateCols))
	}

	// 2. Delete existing notifications
	if e.Notifications == nil {
		return handler.NewMultiStatement(e, statements...), nil
	}
	statements = append(statements, handler.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
		},
		handler.WithTableSuffix(quotaNotificationsTableSuffix),
	))
	notifications := *e.Notifications
	for i := range notifications {
		notification := notifications[i]
		statements = append(statements, handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(QuotaNotificationColumnUnit, e.Unit),
				handler.NewCol(QuotaNotificationColumnID, notification.ID),
				handler.NewCol(QuotaNotificationColumnCallURL, notification.CallURL),
				handler.NewCol(QuotaNotificationColumnPercent, notification.Percent),
				handler.NewCol(QuotaNotificationColumnRepeat, notification.Repeat),
			},
			handler.WithTableSuffix(quotaNotificationsTableSuffix),
		))
	}
	return handler.NewMultiStatement(e, statements...), nil
}

func (q *quotaProjection) reduceQuotaNotificationDue(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.NotificationDueEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(QuotaNotificationColumnLatestDuePeriodStart, e.PeriodStart),
			handler.NewCol(QuotaNotificationColumnNextDueThreshold, e.Threshold+100), // next due_threshold is always the reached + 100 => percent (e.g. 90) in the next bucket (e.g. 190)
		},
		[]handler.Condition{
			handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
			handler.NewCond(QuotaNotificationColumnID, e.ID),
		},
		handler.WithTableSuffix(quotaNotificationsTableSuffix),
	), nil
}

func (q *quotaProjection) reduceQuotaRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.RemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewMultiStatement(
		e,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaPeriodColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaPeriodColumnUnit, e.Unit),
			},
			handler.WithTableSuffix(quotaPeriodsTableSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
			},
			handler.WithTableSuffix(quotaNotificationsTableSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaColumnUnit, e.Unit),
			},
		),
	), nil
}

func (q *quotaProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	// we only assert the event to make sure it is the correct type
	e, err := assertEvent[*instance.InstanceRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewMultiStatement(
		e,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaPeriodColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(quotaPeriodsTableSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(quotaNotificationsTableSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (q *quotaProjection) IncrementUsage(ctx context.Context, unit quota.Unit, instanceID string, periodStart time.Time, count uint64) (sum uint64, err error) {
	if count == 0 {
		return 0, nil
	}

	err = q.client.DB.QueryRowContext(
		ctx,
		incrementQuotaStatement,
		instanceID, unit, periodStart, count,
	).Scan(&sum)
	if err != nil {
		return 0, zerrors.ThrowInternalf(err, "PROJ-SJL3h", "incrementing usage for unit %d failed for at least one quota period", unit)
	}
	return sum, err
}
