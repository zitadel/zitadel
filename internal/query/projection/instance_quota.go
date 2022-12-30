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
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.QuotaRemovedEventType,
					Reduce: reduceQuotaRemoved,
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

func reduceQuotaRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.QuotaRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DlFsg", "reduce.wrong.event.type% s", quota.NotifiedEventType)
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaNotificationInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(QuotaNotificationUnitCol, e.Unit),
		}, crdb.WithTableSuffix(QuotaNotificationsTableSuffix)),
		crdb.AddDeleteStatement([]handler.Condition{
			handler.NewCond(QuotaInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(QuotaUnitCol, e.Unit),
		}),
	), nil
}
