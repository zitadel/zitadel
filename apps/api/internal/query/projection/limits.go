package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/limits"
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
	LimitsColumnBlock             = "block"
)

type limitsProjection struct{}

func newLimitsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &limitsProjection{})
}

func (*limitsProjection) Name() string {
	return LimitsProjectionTable
}

func (*limitsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(LimitsColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(LimitsColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(LimitsColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(LimitsColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(LimitsColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(LimitsColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(LimitsColumnAuditLogRetention, handler.ColumnTypeInterval, handler.Nullable()),
			handler.NewColumn(LimitsColumnBlock, handler.ColumnTypeBool, handler.Nullable()),
		},
			handler.NewPrimaryKey(LimitsColumnInstanceID, LimitsColumnResourceOwner),
		),
	)
}

func (p *limitsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: limits.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
	updateCols := []handler.Column{
		handler.NewCol(LimitsColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(LimitsColumnResourceOwner, e.Aggregate().ResourceOwner),
		handler.NewCol(LimitsColumnCreationDate, handler.OnlySetValueOnInsert(LimitsProjectionTable, e.CreationDate())),
		handler.NewCol(LimitsColumnChangeDate, e.CreationDate()),
		handler.NewCol(LimitsColumnSequence, e.Sequence()),
		handler.NewCol(LimitsColumnAggregateID, e.Aggregate().ID),
	}
	if e.AuditLogRetention != nil {
		updateCols = append(updateCols, handler.NewCol(LimitsColumnAuditLogRetention, *e.AuditLogRetention))
	}
	if e.Block != nil {
		updateCols = append(updateCols, handler.NewCol(LimitsColumnBlock, *e.Block))
	}
	return handler.NewUpsertStatement(e, conflictCols, updateCols), nil
}

func (p *limitsProjection) reduceLimitsReset(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*limits.ResetEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LimitsColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(LimitsColumnResourceOwner, e.Aggregate().ResourceOwner),
		},
	), nil
}
