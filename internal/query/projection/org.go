package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgProjection struct {
	crdb.StatementHandler
}

const (
	OrgProjectionTable = "zitadel.projections.orgs"
)

func NewOrgProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgProjection {
	p := &OrgProjection{}
	config.ProjectionName = OrgProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgProjection) reducers() []handler.AggregateReducer {
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

type OrgColumn int32

const (
	OrgColumnCreationDate OrgColumn = iota + 1
	OrgColumnChangeDate
	OrgColumnResourceOwner
	OrgColumnState
	OrgColumnSequence
	OrgColumnName
	OrgColumnDomain
	OrgColumnID
)

func (c OrgColumn) ColumnName() string {
	switch c {
	case OrgColumnID:
		return "id"
	case OrgColumnCreationDate:
		return "creation_date"
	case OrgColumnChangeDate:
		return "change_date"
	case OrgColumnResourceOwner:
		return "resource_owner"
	case OrgColumnState:
		return "org_state"
	case OrgColumnSequence:
		return "sequence"
	case OrgColumnName:
		return "name"
	case OrgColumnDomain:
		return "primary_domain"
	default:
		return ""
	}
}

func (c OrgColumn) FullColumnName() string {
	return OrgProjectionTable + "." + c.ColumnName()
}

func (p *OrgProjection) reduceOrgAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgAddedEventType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-uYq4r", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnID.ColumnName(), e.Aggregate().ID),
			handler.NewCol(OrgColumnCreationDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnChangeDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnResourceOwner.ColumnName(), e.Aggregate().ResourceOwner),
			handler.NewCol(OrgColumnSequence.ColumnName(), e.Sequence()),
			handler.NewCol(OrgColumnName.ColumnName(), e.Name),
			handler.NewCol(OrgColumnState.ColumnName(), domain.OrgStateActive),
		},
	), nil
}

func (p *OrgProjection) reduceOrgChanged(event eventstore.EventReader) (*handler.Statement, error) {
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
			handler.NewCol(OrgColumnChangeDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnSequence.ColumnName(), e.Sequence()),
			handler.NewCol(OrgColumnName.ColumnName(), e.Name),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID.ColumnName(), e.Aggregate().ID),
		},
	), nil
}

func (p *OrgProjection) reduceOrgDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1gwdc", "seq", event.Sequence(), "expectedType", org.OrgDeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BApK4", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnSequence.ColumnName(), e.Sequence()),
			handler.NewCol(OrgColumnState.ColumnName(), domain.OrgStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID.ColumnName(), e.Aggregate().ID),
		},
	), nil
}

func (p *OrgProjection) reduceOrgReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Vjwiy", "seq", event.Sequence(), "expectedType", org.OrgReactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o37De", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnSequence.ColumnName(), e.Sequence()),
			handler.NewCol(OrgColumnState.ColumnName(), domain.OrgStateActive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID.ColumnName(), e.Aggregate().ID),
		},
	), nil
}

func (p *OrgProjection) reducePrimaryDomainSet(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-79OhB", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4TbKT", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate.ColumnName(), e.CreationDate()),
			handler.NewCol(OrgColumnSequence.ColumnName(), e.Sequence()),
			handler.NewCol(OrgColumnDomain.ColumnName(), e.Domain),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID.ColumnName(), e.Aggregate().ID),
		},
	), nil
}
