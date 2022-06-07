package projection

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	FlowTriggerTable             = "zitadel.projections.flows_triggers"
	FlowTypeCol                  = "flow_type"
	FlowTriggerTypeCol           = "trigger_type"
	FlowResourceOwnerCol         = "resource_owner"
	FlowActionTriggerSequenceCol = "trigger_sequence"
	FlowActionIDCol              = "action_id"
)

type flowProjection struct {
	crdb.StatementHandler
}

func newFlowProjection(ctx context.Context, config crdb.StatementHandlerConfig) *flowProjection {
	p := &flowProjection{}
	config.ProjectionName = FlowTriggerTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *flowProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.TriggerActionsSetEventType,
					Reduce: p.reduceTriggerActionsSetEventType,
				},
				{
					Event:  org.FlowClearedEventType,
					Reduce: p.reduceFlowClearedEventType,
				},
			},
		},
	}
}

func (p *flowProjection) reduceTriggerActionsSetEventType(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.TriggerActionsSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", org.TriggerActionsSetEventType).Error("was not an trigger actions set event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	stmts := make([]func(reader eventstore.Event) crdb.Exec, len(e.ActionIDs)+1)
	stmts[0] = crdb.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(FlowTypeCol, e.FlowType),
			handler.NewCond(FlowTriggerTypeCol, e.TriggerType),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
	)
	for i, id := range e.ActionIDs {
		stmts[i+1] = crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(FlowTypeCol, e.FlowType),
				handler.NewCol(FlowTriggerTypeCol, e.TriggerType),
				handler.NewCol(FlowActionIDCol, id),
				handler.NewCol(FlowActionTriggerSequenceCol, i),
			},
		)
	}
	return crdb.NewMultiStatement(e, stmts...), nil
}

func (p *flowProjection) reduceFlowClearedEventType(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.FlowClearedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", org.FlowClearedEventType).Error("was not a flow cleared event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FlowTypeCol, e.FlowType),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
	), nil
}
