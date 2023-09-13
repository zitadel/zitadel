package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	OrgProjectionTable = "projections.orgs1"

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

type orgProjection struct{}

func (*orgProjection) Name() string {
	return OrgProjectionTable
}

func newOrgProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgProjection))
}

// Init implements [handler.initializer]
func (p *orgProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OrgColumnID, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OrgColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(OrgColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(OrgColumnName, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnDomain, handler.ColumnTypeText, handler.Default("")),
		},
			handler.NewPrimaryKey(OrgColumnInstanceID, OrgColumnID),
			handler.WithIndex(handler.NewIndex("domain", []string{OrgColumnDomain})),
			handler.WithIndex(handler.NewIndex("name", []string{OrgColumnName})),
		),
	)
}

func (p *orgProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimaryDomainSet,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(OrgColumnInstanceID),
				},
			},
		},
	}
}

func (p *orgProjection) reduceOrgAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-uYq4r", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}
	return handler.NewCreateStatement(
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

func (p *orgProjection) reduceOrgChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bg8oM", "reduce.wrong.event.type %s", org.OrgChangedEventType)
	}
	if e.Name == "" {
		return handler.NewNoOpStatement(e), nil
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnName, e.Name),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgProjection) reduceOrgDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-BApK4", "reduce.wrong.event.type %s", org.OrgDeactivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnState, domain.OrgStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgProjection) reduceOrgReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-o37De", "reduce.wrong.event.type %s", org.OrgReactivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnState, domain.OrgStateActive),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4TbKT", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
			handler.NewCol(OrgColumnSequence, e.Sequence()),
			handler.NewCol(OrgColumnDomain, e.Domain),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-DgMSg", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
