package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

type actionProjection struct{}

func newActionProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(actionProjection))
}

func (*actionProjection) Name() string {
	return ActionTable
}

func (*actionProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ActionIDCol, handler.ColumnTypeText),
			handler.NewColumn(ActionCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(ActionChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(ActionResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(ActionInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(ActionStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(ActionSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(ActionNameCol, handler.ColumnTypeText),
			handler.NewColumn(ActionScriptCol, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(ActionTimeoutCol, handler.ColumnTypeInt64, handler.Default(0)),
			handler.NewColumn(ActionAllowedToFailCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(ActionOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(ActionInstanceIDCol, ActionIDCol),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{ActionResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{ActionOwnerRemovedCol})),
		),
	)
}

func (p *actionProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: action.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
	return handler.NewCreateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewDeleteStatement(
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
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ActionInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(ActionResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
