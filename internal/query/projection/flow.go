package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	FlowTriggerTable             = "projections.flow_triggers2"
	FlowTypeCol                  = "flow_type"
	FlowChangeDateCol            = "change_date"
	FlowSequenceCol              = "sequence"
	FlowTriggerTypeCol           = "trigger_type"
	FlowResourceOwnerCol         = "resource_owner"
	FlowInstanceIDCol            = "instance_id"
	FlowActionTriggerSequenceCol = "trigger_sequence"
	FlowActionIDCol              = "action_id"
	FlowOwnerRemovedCol          = "owner_removed"
)

type flowProjection struct {
	crdb.StatementHandler
}

func newFlowProjection(ctx context.Context, config crdb.StatementHandlerConfig) *flowProjection {
	p := new(flowProjection)
	config.ProjectionName = FlowTriggerTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(FlowTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(FlowChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(FlowSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(FlowTriggerTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(FlowResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(FlowInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(FlowActionTriggerSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(FlowActionIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(FlowOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(FlowInstanceIDCol, FlowTypeCol, FlowTriggerTypeCol, FlowResourceOwnerCol, FlowActionIDCol),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{FlowOwnerRemovedCol})),
		),
	)
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
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.TriggerActionsSetEventType)
	}
	stmts := make([]func(reader eventstore.Event) crdb.Exec, len(e.ActionIDs)+1)
	stmts[0] = crdb.AddDeleteStatement(
		[]handler.Condition{
			handler.NewCond(FlowTypeCol, e.FlowType),
			handler.NewCond(FlowTriggerTypeCol, e.TriggerType),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(FlowInstanceIDCol, e.Aggregate().InstanceID),
		},
	)
	for i, id := range e.ActionIDs {
		stmts[i+1] = crdb.AddCreateStatement(
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
	return crdb.NewMultiStatement(e, stmts...), nil
}

func (p *flowProjection) reduceFlowClearedEventType(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.FlowClearedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.FlowClearedEventType)
	}
	return crdb.NewDeleteStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Yd7WC", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FlowInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(FlowResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
