package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	QuotassProjectionTable = "projections.quotas"

	QuotasColumnInstanceID = "instance_id"
	QuotasColumnUnit       = "unit"
	QuotasColumnAmount     = "amount"
	QuotasColumnFrom       = "from_anchor"
	QuotasColumnInterval   = "interval"
	QuotasColumnLimit      = "limit_usage"
)

type quotaProjection struct {
	crdb.StatementHandler
}

func newQuotaProjection(ctx context.Context, config crdb.StatementHandlerConfig) *quotaProjection {
	p := new(quotaProjection)
	config.ProjectionName = QuotassProjectionTable
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
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *quotaProjection) reducers() []handler.AggregateReducer {
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
					Reduce: p.reduceQuotaAdded,
				},
			},
		},
	}
}

func (p *quotaProjection) reduceQuotaAdded(event eventstore.Event) (*handler.Statement, error) {
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
