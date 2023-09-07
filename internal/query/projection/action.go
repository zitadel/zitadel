package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	ActionTable            = "projections.actions3"
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
	ActionOwnerRemovedCol  = "owner_removed"
)

type actionProjection struct {
	crdb.StatementHandler
}

func newActionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *actionProjection {
	p := new(actionProjection)
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
			crdb.NewColumn(ActionOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(ActionInstanceIDCol, ActionIDCol),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{ActionResourceOwnerCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{ActionOwnerRemovedCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *actionProjection) reducers() []handler.AggregateReducer {
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
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
					Reduce: reduceInstanceRemovedHelper(ActionInstanceIDCol),
				},
			},
		},
	}
}

func (p *actionProjection) reduceActionAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *actionProjection) reduceActionChanged(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCond(ActionInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *actionProjection) reduceActionDeactivated(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCond(ActionInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *actionProjection) reduceActionReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.ReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-hwdqa", "reduce.wrong.event.type %s", action.ReactivatedEventType)
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
			handler.NewCond(ActionInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *actionProjection) reduceActionRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.RemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dgh2d", "reduce.wrong.event.type %s", action.RemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
			handler.NewCond(ActionInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *actionProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-mSmWM", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ActionInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(ActionResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
