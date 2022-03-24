package projection

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/org"
)

const (
	OrgProjectionTable = "projections.orgs"

	OrgColumnID            = "id"
	OrgColumnCreationDate  = "creation_date"
	OrgColumnChangeDate    = "change_date"
	OrgColumnResourceOwner = "resource_owner"
	OrgColumnInstanceID    = "instance_id"
	OrgColumnState         = "org_state"
	OrgColumnSequence      = "sequence"
	OrgColumnName          = "name"
	OrgColumnDomain        = "primary_domain"
)

type OrgProjection struct {
	crdb.StatementHandler
}

func NewOrgProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgProjection {
	p := new(OrgProjection)
	config.ProjectionName = OrgProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(OrgColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(OrgColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OrgColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OrgColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(OrgColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(OrgColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(OrgColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(OrgColumnName, crdb.ColumnTypeText),
			crdb.NewColumn(OrgColumnDomain, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(OrgColumnInstanceID, OrgColumnID),
			crdb.WithIndex(crdb.NewIndex("domain_idx", []string{OrgColumnDomain})),
			crdb.WithIndex(crdb.NewIndex("name_idx", []string{OrgColumnName})),
		),
	)
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

func (p *OrgProjection) reduceOrgAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnID, e.Aggregate().ID),
			handler.NewCol(OrgColumnCreationDate, e.CreationDate()),
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(OrgColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnName, e.Name),
			handler.NewCol(OrgColumnState, domain.OrgStateActive),
		},
	), nil
}

func (p *OrgProjection) reduceOrgChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bg8oM", "reduce.wrong.event.type %s", org.OrgChangedEventType)
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

func (p *OrgProjection) reduceOrgDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-BApK4", "reduce.wrong.event.type %s", org.OrgDeactivatedEventType)
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

func (p *OrgProjection) reduceOrgReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-o37De", "reduce.wrong.event.type %s", org.OrgReactivatedEventType)
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

func (p *OrgProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4TbKT", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
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
