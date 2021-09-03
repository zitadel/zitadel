package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/action"
)

type ActionProjection struct {
	crdb.StatementHandler
}

func NewActionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ActionProjection {
	p := &ActionProjection{}
	config.ProjectionName = "projections.actions"
	config.Reducers = p.reducers()
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

const (
	actionIDCol            = "id"
	actionCreationDateCol  = "creation_date"
	actionChangeDateCol    = "change_date"
	actionResourceOwnerCol = "resource_owner"
	actionStateCol         = "action_state"
	actionSequenceCol      = "sequence"
	actionNameCol          = "name"
)

func (p *ActionProjection) reduceActionAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.AddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence, "expectedType", action.AddedEventType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(actionIDCol, e.Aggregate().ID),
			handler.NewCol(actionCreationDateCol, e.CreationDate()),
			handler.NewCol(actionChangeDateCol, e.CreationDate()),
			handler.NewCol(actionResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(actionSequenceCol, e.Sequence()),
			handler.NewCol(actionNameCol, e.Name),
			handler.NewCol(actionStateCol, domain.ActionStateActive),
		},
	), nil
}

func (p *ActionProjection) reduceActionChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.ChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q4oq8", "seq", event.Sequence, "expected", action.ChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	values := []handler.Column{
		handler.NewCol(actionChangeDateCol, e.CreationDate()),
		handler.NewCol(actionSequenceCol, e.Sequence()),
	}
	if e.Name != nil {
		values = append(values, handler.NewCol(actionNameCol, e.Name))
	}
	return crdb.NewUpdateStatement(
		e,
		values,
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.DeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1gwdc", "seq", event.Sequence, "expectedType", action.DeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BApK4", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(actionChangeDateCol, e.CreationDate()),
			handler.NewCol(actionSequenceCol, e.Sequence()),
			handler.NewCol(actionStateCol, domain.ActionStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.ReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Vjwiy", "seq", event.Sequence, "expectedType", action.ReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o37De", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(actionChangeDateCol, e.CreationDate()),
			handler.NewCol(actionSequenceCol, e.Sequence()),
			handler.NewCol(actionStateCol, domain.ActionStateActive),
		},
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ActionProjection) reduceActionRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*action.RemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence, "expectedType", action.RemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(actionIDCol, e.Aggregate().ID),
		},
	), nil
}
