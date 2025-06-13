package projection

import (
	"context"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	OrgRelationProjectionTable = "zitadel.organizations"
)

type orgRelationalProjection struct{}

func (*orgRelationalProjection) Name() string {
	return OrgRelationProjectionTable
}

func newOrgRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgRelationalProjection))
}

// Init implements [handler.initializer]
func (p *orgRelationalProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OrgColumnID, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnName, handler.ColumnTypeText),
			handler.NewColumn(OrgColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(State, handler.ColumnTypeEnum),
			handler.NewColumn(CreatedAt, handler.ColumnTypeTimestamp),
			handler.NewColumn(UpdatedAt, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeletedAt, handler.ColumnTypeTimestamp),
		},
			handler.NewPrimaryKey(OrgColumnInstanceID, OrgColumnID),
			// handler.WithIndex(handler.NewIndex("domain", []string{OrgColumnDomain})),
			handler.WithIndex(handler.NewIndex("name", []string{OrgColumnName})),
		),
	)
}

func (p *orgRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgRelationalAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgRelationalChanged,
				},
				{
					Event:  org.OrgDeactivatedEventType,
					Reduce: p.reduceOrgRelationalDeactivated,
				},
				{
					Event:  org.OrgReactivatedEventType,
					Reduce: p.reduceOrgRelationalReactivated,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRelationalRemoved,
				},
				// TODO
				// {
				// 	Event:  org.OrgDomainPrimarySetEventType,
				// 	Reduce: p.reducePrimaryDomainSetRelational,
				// },
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

func (p *orgRelationalProjection) reduceOrgRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq5R", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}

	// need to convert old state (int) to new state (enum)
	state := repoDomain.State[domain.OrgStateActive-1]

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnID, e.Aggregate().ID),
			handler.NewCol(OrgColumnName, e.Name),
			handler.NewCol(OrgColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(State, state),
			handler.NewCol(CreatedAt, e.CreationDate()),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
	), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bg9om", "reduce.wrong.event.type %s", org.OrgChangedEventType)
	}
	if e.Name == "" {
		return handler.NewNoOpStatement(e), nil
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnName, e.Name),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BApK5", "reduce.wrong.event.type %s", org.OrgDeactivatedEventType)
	}

	// need to convert old state (int) to new state (enum)
	state := repoDomain.State[domain.OrgStateInactive-1]
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(State, state),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-o38DE", "reduce.wrong.event.type %s", org.OrgReactivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(State, domain.OrgStateActive),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

// TODO
// func (p *orgRelationalProjection) reducePrimaryDomainSetRelational(event eventstore.Event) (*handler.Statement, error) {
// 	e, ok := event.(*org.DomainPrimarySetEvent)
// 	if !ok {
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-3Tbkt", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
// 	}
// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(OrgColumnChangeDate, e.CreationDate()),
// 			handler.NewCol(OrgColumnSequence, e.Sequence()),
// 			handler.NewCol(OrgColumnDomain, e.Domain),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(OrgColumnID, e.Aggregate().ID),
// 			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
// 		},
// 	), nil
// }

func (p *orgRelationalProjection) reduceOrgRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DGm9g", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UpdatedAt, e.CreationDate()),
			handler.NewCol(DeletedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
