package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	QuotaTable                    = "projections.instance_quotas"
	QuotaNotificationsTableSuffix = "notifications"

	QuotaIDCol            = "id"
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

	QuotaNotificationIDCol                = "id"
	QuotaNotificationQuotaIDCol           = "quota_id"
	QuotaNotificationInstanceIDCol        = "instance_id"
	QuotaNotificationUnitCol              = "unit"
	QuotaNotificationCallURLCol           = "call_url"
	QuotaNotificationPercentCol           = "percent"
	QuotaNotificationRepeatCol            = "repeat"
	QuotaNotificationLastCallDateCol      = "last_call_date"
	QuotaNotificationLastCallThresholdCol = "last_call_threshold"
)

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
				crdb.NewColumn(QuotaIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaSequenceCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaUnitCol, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaFromCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaIntervalCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaAmountCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaLimitCol, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotaIDCol),
			crdb.WithIndex(crdb.NewIndex("quotas_ro_idx", []string{QuotaResourceOwnerCol})),
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaNotificationIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationQuotaIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationInstanceIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationUnitCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationCallURLCol, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaNotificationPercentCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotaNotificationRepeatCol, crdb.ColumnTypeBool),
				crdb.NewColumn(QuotaNotificationLastCallDateCol, crdb.ColumnTypeTimestamp, crdb.Nullable()),
				crdb.NewColumn(QuotaNotificationLastCallThresholdCol, crdb.ColumnTypeInt64, crdb.Nullable()),
			},
			crdb.NewPrimaryKey(QuotaNotificationIDCol),
			QuotaNotificationsTableSuffix,
			crdb.WithForeignKey(
				crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
			),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, esHandlerConfig)
	return p
}

func esReducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: quota.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.AddedEventType,
					Reduce: reduceQuotaAdded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.NotifiedEventType,
					Reduce: reduceQuotaNotified,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.RemovedEventType,
					Reduce: reduceQuotaRemoved,
				},
			},
		},
	}
}

func reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*quota.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff21", "reduce.wrong.event.type% s", quota.AddedEventType)
	}

	execFuncs := []func(eventstore.Event) crdb.Exec{
		crdb.AddCreateStatement([]handler.Column{
			handler.NewCol(QuotaIDCol, e.Aggregate().ID),
			handler.NewCol(QuotaCreationDateCol, e.CreationDate()),
			handler.NewCol(QuotaInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(QuotaResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(QuotaChangeDateCol, e.CreationDate()),
			handler.NewCol(QuotaSequenceCol, e.Sequence()),
			handler.NewCol(QuotaUnitCol, e.Unit),
			handler.NewCol(QuotaFromCol, e.From),
			handler.NewCol(QuotaIntervalCol, e.ResetInterval),
			handler.NewCol(QuotaAmountCol, e.Amount),
			handler.NewCol(QuotaLimitCol, e.Limit),
		}),
	}

	for _, notification := range e.Notifications {

		execFuncs = append(execFuncs, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(QuotaNotificationIDCol, notification.ID),
				handler.NewCol(QuotaNotificationQuotaIDCol, e.Aggregate().ID),
				handler.NewCol(QuotaNotificationInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(QuotaNotificationUnitCol, e.Unit),
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
	e, ok := event.(*quota.NotifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-wmQPo", "reduce.wrong.event.type% s", quota.NotifiedEventType)
	}

	return crdb.NewUpdateStatement(event, []handler.Column{
		handler.NewCol(QuotaNotificationLastCallDateCol, e.CreationDate()),
		handler.NewCol(QuotaNotificationLastCallThresholdCol, e.Threshold),
	}, []handler.Condition{
		handler.NewCond(QuotaNotificationIDCol, e.ID),
	}, crdb.WithTableSuffix(QuotaNotificationsTableSuffix)), nil
}

func reduceQuotaRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*quota.RemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DlFsg", "reduce.wrong.event.type% s", quota.NotifiedEventType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaNotificationInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationQuotaIDCol, e.Aggregate().ID),
		}, crdb.WithTableSuffix(QuotaNotificationsTableSuffix)),
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCond(QuotaIDCol, e.Aggregate().ID),
		}),
	), nil
}

func reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	return crdb.NewMultiStatement(
		event,
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaNotificationInstanceIDCol, event.Aggregate().InstanceID),
		}, crdb.WithTableSuffix(QuotaNotificationsTableSuffix)),
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaInstanceIDCol, event.Aggregate().InstanceID),
		}),
	), nil
}
