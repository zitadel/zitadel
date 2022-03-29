package projection

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/action"
)

const (
	ActionTable            = "projections.actions"
	ActionIDCol            = "id"
	ActionCreationDateCol  = "creation_date"
	ActionChangeDateCol    = "change_date"
	ActionResourceOwnerCol = "resource_owner"
	ActionInstanceIDCol    = "instance_id"
	ActionStateCol         = "action_state"
	ActionSequenceCol      = "sequence"
	ActionNameCol          = "name"
	ActionScriptCol        = "script"
	ActionTimeoutCol       = "timeout"
	ActionAllowedToFailCol = "allowed_to_fail"
)

type ActionProjection struct {
	crdb.StatementHandler
}

func NewActionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ActionProjection {
	p := new(ActionProjection)
	config.ProjectionName = ActionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(ActionIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(ActionCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ActionChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ActionResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(ActionInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(ActionStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(ActionSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(ActionNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(ActionScriptCol, crdb.ColumnTypeText, crdb.Default("")),
			crdb.NewColumn(ActionTimeoutCol, crdb.ColumnTypeInt64, crdb.Default(0)),
			crdb.NewColumn(ActionAllowedToFailCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(ActionInstanceIDCol, ActionIDCol),
			crdb.WithIndex(crdb.NewIndex("ro_idx", []string{ActionResourceOwnerCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ActionProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: action.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  action.AddedEventType,
					Reduce: p.reduceActionAdded,
				},
				{
					Event:  action.ChangedEventType,
					Reduce: p.reduceActionChanged,
				},
				{
					Event:  action.DeactivatedEventType,
					Reduce: p.reduceActionDeactivated,
				},
				{
					Event:  action.ReactivatedEventType,
					Reduce: p.reduceActionReactivated,
				},
				{
					Event:  action.RemovedEventType,
					Reduce: p.reduceActionRemoved,
				},
			},
		},
	}
}

func (p *ActionProjection) reduceActionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff21", "reduce.wrong.event.type% s", action.AddedEventType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ActionIDCol, e.Aggregate().ID),
			handler.NewCol(ActionCreationDateCol, e.CreationDate()),
			handler.NewCol(ActionChangeDateCol, e.CreationDate()),
			handler.NewCol(ActionResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(ActionInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCol(ActionSequenceCol, e.Sequence()),
			handler.NewCol(ActionNameCol, e.Name),
			handler.NewCol(ActionScriptCol, e.Script),
			handler.NewCol(ActionTimeoutCol, e.Timeout),
			handler.NewCol(ActionAllowedToFailCol, e.AllowedToFail),
			handler.NewCol(ActionStateCol, domain.ActionStateActive),
		},
	), nil
}

func (p *ActionProjection) reduceActionChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.ChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Gg43d", "reduce.wrong.event.type %s", action.ChangedEventType)
	}
	values := []handler.Column{
		handler.NewCol(ActionChangeDateCol, e.CreationDate()),
		handler.NewCol(ActionSequenceCol, e.Sequence()),
	}
	if e.Name != nil {
		values = append(values, handler.NewCol(ActionNameCol, *e.Name))
	}
	if e.Script != nil {
		values = append(values, handler.NewCol(ActionScriptCol, *e.Script))
	}
	if e.Timeout != nil {
		values = append(values, handler.NewCol(ActionTimeoutCol, *e.Timeout))
	}
	if e.AllowedToFail != nil {
		values = append(values, handler.NewCol(ActionAllowedToFailCol, *e.AllowedToFail))
	}
	return crdb.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.DeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Fgh32", "reduce.wrong.event.type %s", action.DeactivatedEventType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ActionChangeDateCol, e.CreationDate()),
			handler.NewCol(ActionSequenceCol, e.Sequence()),
			handler.NewCol(ActionStateCol, domain.ActionStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.ReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-hwdqa", "reduce.wrong.event.type% s", action.ReactivatedEventType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ActionChangeDateCol, e.CreationDate()),
			handler.NewCol(ActionSequenceCol, e.Sequence()),
			handler.NewCol(ActionStateCol, domain.ActionStateActive),
		},
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.RemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dgh2d", "reduce.wrong.event.type% s", action.RemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
		},
	), nil
}
