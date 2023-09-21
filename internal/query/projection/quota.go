package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	zitadel_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
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
	crdb.StatementHandler
	client *database.DB
}

func newQuotaProjection(ctx context.Context, config crdb.StatementHandlerConfig) *quotaProjection {
	p := new(quotaProjection)
	config.ProjectionName = QuotasProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaColumnID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaColumnInstanceID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaColumnUnit, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaColumnAmount, crdb.ColumnTypeInt64, crdb.Nullable()),
				crdb.NewColumn(QuotaColumnFrom, crdb.ColumnTypeTimestamp, crdb.Nullable()),
				crdb.NewColumn(QuotaColumnInterval, crdb.ColumnTypeInterval, crdb.Nullable()),
				crdb.NewColumn(QuotaColumnLimit, crdb.ColumnTypeBool, crdb.Nullable()),
			},
			crdb.NewPrimaryKey(QuotaColumnInstanceID, QuotaColumnUnit),
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaPeriodColumnInstanceID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaPeriodColumnUnit, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaPeriodColumnStart, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaPeriodColumnUsage, crdb.ColumnTypeInt64),
			},
			crdb.NewPrimaryKey(QuotaPeriodColumnInstanceID, QuotaPeriodColumnUnit, QuotaPeriodColumnStart),
			quotaPeriodsTableSuffix,
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaNotificationColumnInstanceID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationColumnUnit, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaNotificationColumnID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationColumnCallURL, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationColumnPercent, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaNotificationColumnRepeat, crdb.ColumnTypeBool),
				crdb.NewColumn(QuotaNotificationColumnLatestDuePeriodStart, crdb.ColumnTypeTimestamp, crdb.Nullable()),
				crdb.NewColumn(QuotaNotificationColumnNextDueThreshold, crdb.ColumnTypeInt64, crdb.Nullable()),
			},
			crdb.NewPrimaryKey(QuotaNotificationColumnInstanceID, QuotaNotificationColumnUnit, QuotaNotificationColumnID),
			quotaNotificationsTableSuffix,
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.client = config.Client
	return p
}

func (q *quotaProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: q.reduceInstanceRemoved,
				},
			},
		},
		{
			Aggregate: quota.AggregateType,
			EventRedusers: []handler.EventReducer{
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
	return crdb.NewNoOpStatement(event), nil
}

func (q *quotaProjection) reduceQuotaSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.SetEvent](event)
	if err != nil {
		return nil, err
	}
	var statements []func(e eventstore.Event) crdb.Exec

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
		statements = append(statements, crdb.AddUpsertStatement(quotaConflictColumns, quotaUpdateCols))
	}

	// 2. Delete existing notifications
	if e.Notifications == nil {
		return crdb.NewMultiStatement(e, statements...), nil
	}
	statements = append(statements, crdb.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
		},
		crdb.WithTableSuffix(quotaNotificationsTableSuffix),
	))
	notifications := *e.Notifications
	for i := range notifications {
		notification := notifications[i]
		statements = append(statements, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(QuotaNotificationColumnUnit, e.Unit),
				handler.NewCol(QuotaNotificationColumnID, notification.ID),
				handler.NewCol(QuotaNotificationColumnCallURL, notification.CallURL),
				handler.NewCol(QuotaNotificationColumnPercent, notification.Percent),
				handler.NewCol(QuotaNotificationColumnRepeat, notification.Repeat),
			},
			crdb.WithTableSuffix(quotaNotificationsTableSuffix),
		))
	}
	return crdb.NewMultiStatement(e, statements...), nil
}

func (q *quotaProjection) reduceQuotaNotificationDue(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.NotificationDueEvent](event)
	if err != nil {
		return nil, err
	}
	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(QuotaNotificationColumnLatestDuePeriodStart, e.PeriodStart),
			handler.NewCol(QuotaNotificationColumnNextDueThreshold, e.Threshold+100), // next due_threshold is always the reached + 100 => percent (e.g. 90) in the next bucket (e.g. 190)
		},
		[]handler.Condition{
			handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
			handler.NewCond(QuotaNotificationColumnID, e.ID),
		},
		crdb.WithTableSuffix(quotaNotificationsTableSuffix),
		// The notification could have been removed in the meantime
		crdb.WithIgnoreNotFound(),
	), nil
}

func (q *quotaProjection) reduceQuotaRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.RemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaPeriodColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaPeriodColumnUnit, e.Unit),
			},
			crdb.WithTableSuffix(quotaPeriodsTableSuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaNotificationColumnUnit, e.Unit),
			},
			crdb.WithTableSuffix(quotaNotificationsTableSuffix),
		),
		crdb.AddDeleteStatement(
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
	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaPeriodColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(quotaPeriodsTableSuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotaNotificationColumnInstanceID, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(quotaNotificationsTableSuffix),
		),
		crdb.AddDeleteStatement(
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
		return 0, zitadel_errors.ThrowInternalf(err, "PROJ-SJL3h", "incrementing usage for unit %d failed for at least one quota period", unit)
	}
	return sum, err
}
