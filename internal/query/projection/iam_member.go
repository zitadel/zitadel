package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMMemberProjection struct {
	crdb.StatementHandler
}

const (
	IAMMemberProjectionTable = "zitadel.projections.iam_members"
)

func NewIAMMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *IAMMemberProjection {
	p := &IAMMemberProjection{}
	config.ProjectionName = IAMMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *IAMMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.MemberAddedEventType,
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  iam.MemberChangedEventType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  iam.MemberCascadeRemovedEventType,
					Reduce: p.reduceIAMMemberCascadeRemoved,
				},
				{
					Event:  iam.MemberRemovedEventType,
					Reduce: p.reduceIAMMemberRemoved,
				},
			},
		},
	}
}

type IAMMemberColumn string

const (
	IAMMemberOrgIDCol  = "iam_id"
	IAMMemberUserIDCol = "user_id"
	IAMMemberRolesCol  = "roles"
)

func (p *IAMMemberProjection) reduceOrgAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GG4Tq", "seq", event.Sequence(), "expectedType", iam.MemberAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-8cz08", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IAMMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCol(IAMMemberUserIDCol, e.UserID),
			handler.NewCol(IAMMemberRolesCol, e.Roles),
		},
	), nil
}

func (p *IAMMemberProjection) reduceOrgChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-cssba", "seq", event.Sequence(), "expected", iam.MemberChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-oDv1T", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IAMMemberRolesCol, e.Roles),
		},
		[]handler.Condition{
			handler.NewCond(IAMMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(IAMMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *IAMMemberProjection) reduceIAMMemberCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QEvzT", "seq", event.Sequence(), "expected", iam.MemberCascadeRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-sWq0V", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(IAMMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(IAMMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *IAMMemberProjection) reduceIAMMemberRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*iam.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-QPcBD", "seq", event.Sequence(), "expected", iam.MemberRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-rIE9Z", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(IAMMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(IAMMemberUserIDCol, e.UserID),
		},
	), nil
}
