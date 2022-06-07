package projection

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type orgProjection struct {
	crdb.StatementHandler
}

const (
	OrgProjectionTable = "zitadel.projections.orgs"
)

func newOrgProjection(ctx context.Context, config crdb.StatementHandlerConfig) *orgProjection {
	p := &orgProjection{}
	config.ProjectionName = OrgProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *orgProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  org.OrgDeactivatedEventType,
					Reduce: p.reduceOrgDeactivated,
				},
				{
					Event:  org.OrgReactivatedEventType,
					Reduce: p.reduceOrgReactivated,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimaryDomainSet,
				},
			},
		},
	}
}

type OrgColumn string

const (
	OrgColumnID            = "id"
	OrgColumnCreationDate  = "creation_date"
	OrgColumnChangeDate    = "change_date"
	OrgColumnResourceOwner = "resource_owner"
	OrgColumnState         = "org_state"
	OrgColumnSequence      = "sequence"
	OrgColumnName          = "name"
	OrgColumnDomain        = "primary_domain"
)

func (p *orgProjection) reduceOrgAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnID, e.Aggregate().ID),
			handler.NewCol(OrgColumnCreationDate, e.CreationDate()),
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnName, e.Name),
			handler.NewCol(OrgColumnState, domain.OrgStateActive),
		},
	), nil
}

func (p *orgProjection) reduceOrgChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q4oq8", "seq", event.Sequence(), "expected", org.OrgChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bg8oM", "reduce.wrong.event.type")
	}
	if e.Name == "" {
		return crdb.NewNoOpStatement(e), nil
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnName, e.Name),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *orgProjection) reduceOrgDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1gwdc", "seq", event.Sequence(), "expectedType", org.OrgDeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BApK4", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnState, domain.OrgStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *orgProjection) reduceOrgReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Vjwiy", "seq", event.Sequence(), "expectedType", org.OrgReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o37De", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnState, domain.OrgStateActive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *orgProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnDomain, e.Domain),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
		},
	), nil
}
