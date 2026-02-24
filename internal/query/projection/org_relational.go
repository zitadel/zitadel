package projection

import (
	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceOrgRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq5R", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OrgColumnID, e.Aggregate().ID),
			handler.NewCol(OrgColumnName, e.Name),
			handler.NewCol(OrgColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(State, repoDomain.OrgStateActive),
			handler.NewCol(CreatedAt, e.CreationDate()),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceOrgRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *relationalTablesProjection) reduceOrgRelationalDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BApK5", "reduce.wrong.event.type %s", org.OrgDeactivatedEventType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(State, repoDomain.OrgStateInactive),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *relationalTablesProjection) reduceOrgRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-o38DE", "reduce.wrong.event.type %s", org.OrgReactivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(State, repoDomain.OrgStateActive),
			handler.NewCol(UpdatedAt, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *relationalTablesProjection) reduceOrgRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DGm9g", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(OrgColumnID, e.Aggregate().ID),
			handler.NewCond(OrgColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
