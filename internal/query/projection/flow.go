package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	FlowTriggerTable             = "projections.flow_triggers3"
	FlowTypeCol                  = "flow_type"
	FlowChangeDateCol            = "change_date"
	FlowSequenceCol              = "sequence"
	FlowTriggerTypeCol           = "trigger_type"
	FlowResourceOwnerCol         = "resource_owner"
	FlowInstanceIDCol            = "instance_id"
	FlowActionTriggerSequenceCol = "trigger_sequence"
	FlowActionIDCol              = "action_id"
)

type flowProjection struct{}

func newFlowProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(flowProjection))
}

func (*flowProjection) Name() string {
	return FlowTriggerTable
}

func (*flowProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(FlowTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(FlowChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(FlowSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(FlowTriggerTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(FlowResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(FlowInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(FlowActionTriggerSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(FlowActionIDCol, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(FlowInstanceIDCol, FlowTypeCol, FlowTriggerTypeCol, FlowResourceOwnerCol, FlowActionIDCol),
		),
	)
}

func (p *flowProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.TriggerActionsSetEventType,
					Reduce: p.reduceTriggerActionsSetEventType,
				},
				{
					Event:  org.FlowClearedEventType,
					Reduce: p.reduceFlowClearedEventType,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(FlowInstanceIDCol),
				},
			},
		},
	}
}

func (p *flowProjection) reduceTriggerActionsSetEventType(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.TriggerActionsSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.TriggerActionsSetEventType)
	}
	stmts := make([]func(reader eventstore.Event) handler.Exec, len(e.ActionIDs)+1)
	stmts[0] = handler.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(FlowTypeCol, e.FlowType),
			handler.NewCond(FlowTriggerTypeCol, e.TriggerType),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(FlowInstanceIDCol, e.Aggregate().InstanceID),
		},
	)
	for i, id := range e.ActionIDs {
		stmts[i+1] = handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(FlowInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(FlowTypeCol, e.FlowType),
				handler.NewCol(FlowChangeDateCol, e.CreationDate()),
				handler.NewCol(FlowSequenceCol, e.Sequence()),
				handler.NewCol(FlowTriggerTypeCol, e.TriggerType),
				handler.NewCol(FlowActionIDCol, id),
				handler.NewCol(FlowActionTriggerSequenceCol, i),
			},
		)
	}
	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *flowProjection) reduceFlowClearedEventType(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.FlowClearedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.FlowClearedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FlowTypeCol, e.FlowType),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(FlowInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *flowProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Yd7WC", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FlowInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
