package projection

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/action"
)

const (
	ActionTable            = "zitadel.projections.actions"
	ActionIDCol            = "id"
	ActionCreationDateCol  = "creation_date"
	ActionChangeDateCol    = "change_date"
	ActionResourceOwnerCol = "resource_owner"
	ActionStateCol         = "action_state"
	ActionSequenceCol      = "sequence"
	ActionNameCol          = "name"
	ActionScriptCol        = "script"
	ActionTimeoutCol       = "timeout"
	ActionAllowedToFailCol = "allowed_to_fail"
)

type actionProjection struct {
	crdb.StatementHandler
}

func newActionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *actionProjection {
	p := &actionProjection{}
	config.ProjectionName = ActionTable
	config.Reducers = p.reducers()
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
	}
}

func (p *actionProjection) reduceActionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.AddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Sgg31", "seq", event.Sequence, "expectedType", action.AddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dff21", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ActionIDCol, e.Aggregate().ID),
			handler.NewCol(ActionCreationDateCol, e.CreationDate()),
			handler.NewCol(ActionChangeDateCol, e.CreationDate()),
			handler.NewCol(ActionResourceOwnerCol, e.Aggregate().ResourceOwner),
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
		logging.LogWithFields("HANDL-Dg2th", "seq", event.Sequence, "expected", action.ChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Gg43d", "reduce.wrong.event.type")
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

func (p *actionProjection) reduceActionDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.DeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Fhhjd", "seq", event.Sequence, "expectedType", action.DeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Fgh32", "reduce.wrong.event.type")
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

func (p *actionProjection) reduceActionReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.ReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Fg4r3", "seq", event.Sequence, "expectedType", action.ReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-hwdqa", "reduce.wrong.event.type")
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

func (p *actionProjection) reduceActionRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*action.RemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Dgwh2", "seq", event.Sequence, "expectedType", action.RemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dgh2d", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ActionIDCol, e.Aggregate().ID),
		},
	), nil
}
