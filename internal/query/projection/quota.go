package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

const (
	QuotasProjectionTable = "projections.quotas"

	QuotasColumnInstanceID = "instance_id"
	QuotasColumnUnit       = "unit"
	QuotasColumnAmount     = "amount"
	QuotasColumnFrom       = "from_anchor"
	QuotasColumnInterval   = "interval"
	QuotasColumnLimit      = "limit_usage"

	periodsTableSuffix           = "periods"
	QuotaPeriodsColumnInstanceID = "instance_id"
	QuotaPeriodsColumnUnit       = "unit"
	QuotaPeriodsColumnStart      = "start"
	QuotaPeriodsColumnUsage      = "usage"
)

type quotaProjection struct {
	crdb.StatementHandler
}

func newQuotaProjection(ctx context.Context, config crdb.StatementHandlerConfig) *quotaProjection {
	p := new(quotaProjection)
	config.ProjectionName = QuotasProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotasColumnInstanceID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotasColumnUnit, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotasColumnAmount, crdb.ColumnTypeInt64),
				crdb.NewColumn(QuotasColumnFrom, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotasColumnInterval, crdb.ColumnTypeInterval),
				crdb.NewColumn(QuotasColumnLimit, crdb.ColumnTypeBool),
			},
			crdb.NewPrimaryKey(QuotasColumnInstanceID, QuotasColumnUnit),
		),
		crdb.NewSuffixedTable(
			[]*crdb.Column{
				crdb.NewColumn(QuotaPeriodsColumnInstanceID, crdb.ColumnTypeText),
				crdb.NewColumn(QuotaPeriodsColumnUnit, crdb.ColumnTypeEnum),
				crdb.NewColumn(QuotaPeriodsColumnStart, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(QuotaPeriodsColumnUsage, crdb.ColumnTypeInt64),
			},
			crdb.NewPrimaryKey(QuotaPeriodsColumnInstanceID, QuotaPeriodsColumnUnit, QuotaPeriodsColumnStart),
			periodsTableSuffix,
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (q *quotaProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(QuotasColumnInstanceID),
				},
			},
		},
		{
			Aggregate: quota.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.AddedEventType,
					Reduce: q.reduceQuotaAdded,
				},
			},
		},
		{
			Aggregate: quota.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  quota.RemovedEventType,
					Reduce: q.reduceQuotaRemoved,
				},
			},
		},
	}
}

func (q *quotaProjection) reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*quota.AddedEvent](event)
	if err != nil {
		return nil, err
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(QuotasColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(QuotasColumnUnit, e.Unit),
			handler.NewCol(QuotasColumnAmount, e.Amount),
			handler.NewCol(QuotasColumnFrom, e.From),
			handler.NewCol(QuotasColumnInterval, e.ResetInterval),
			handler.NewCol(QuotasColumnLimit, e.Limit),
		},
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
				handler.NewCond(QuotaPeriodsColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotaPeriodsColumnUnit, e.Unit),
			},
			crdb.WithTableSuffix(periodsTableSuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(QuotasColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(QuotasColumnUnit, e.Unit),
			},
		),
	), nil
}

func (q *quotaProjection) IncrementAccessLogs(records []*record.AccessLog) error {
	byInstance := make(map[string]uint64)
	for _, r := range records {
		if r.IsAuthenticated() {
			byInstance[r.InstanceID]++
		}
	}
	for instanceID, count := range byInstance {
		if err := incrementUsage(quota.RequestsAllAuthenticated, instanceID, count); err != nil {
			return err
		}
	}
	return nil
}

func incrementUsage(unit quota.Unit, instanceID string, count uint64) error {
	/*	crdb.NewUpsertStatement(
		pseudo.ScheduledEvent{},
		[]handler.Column{},
	)*/
	return nil
}
