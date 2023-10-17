package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/limits"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	LimitsProjectionTable = "projections.limits"

	LimitsColumnAggregateID   = "aggregate_id"
	LimitsColumnCreationDate  = "creation_date"
	LimitsColumnChangeDate    = "change_date"
	LimitsColumnResourceOwner = "resource_owner"
	LimitsColumnInstanceID    = "instance_id"
	LimitsColumnSequence      = "sequence"

	LimitsColumnAuditLogRetention = "audit_log_retention"
)

type limitsProjection struct {
	crdb.StatementHandler
}

func newLimitsProjection(ctx context.Context, config crdb.StatementHandlerConfig) *limitsProjection {
	p := new(limitsProjection)
	config.ProjectionName = LimitsProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(LimitsColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(LimitsColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LimitsColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LimitsColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(LimitsColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(LimitsColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(LimitsColumnAuditLogRetention, crdb.ColumnTypeInterval, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(LimitsColumnInstanceID, LimitsColumnResourceOwner),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *limitsProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: limits.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  limits.SetEventType,
					Reduce: p.reduceLimitsSet,
				},
				{
					Event:  limits.ResetEventType,
					Reduce: p.reduceLimitsReset,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(LimitsColumnInstanceID),
				},
			},
		},
	}
}

func (p *limitsProjection) reduceLimitsSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*limits.SetEvent](event)
	if err != nil {
		return nil, err
	}
	conflictCols := []handler.Column{
		handler.NewCol(LimitsColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(LimitsColumnResourceOwner, e.Aggregate().ResourceOwner),
	}
	updateCols := append(conflictCols,
		handler.NewCol(LimitsColumnCreationDate, e.CreationDate()),
		handler.NewCol(LimitsColumnChangeDate, e.CreationDate()),
		handler.NewCol(LimitsColumnSequence, e.Sequence()),
		handler.NewCol(LimitsColumnAggregateID, e.Aggregate().ID),
	)
	if e.AuditLogRetention != nil {
		updateCols = append(updateCols, handler.NewCol(LimitsColumnAuditLogRetention, *e.AuditLogRetention))
	}
	return crdb.NewUpsertStatement(e, conflictCols, updateCols), nil
}

func (p *limitsProjection) reduceLimitsReset(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*limits.ResetEvent](event)
	if err != nil {
		return nil, err
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LimitsColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(LimitsColumnResourceOwner, e.Aggregate().ResourceOwner),
		},
	), nil
}
