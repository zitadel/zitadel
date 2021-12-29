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
	"github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/lib/pq"
)

type UserGrantProjection struct {
	crdb.StatementHandler
}

const (
	UserGrantProjectionTable = "zitadel.projections.user_grants"
)

func NewUserGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *UserGrantProjection {
	p := &UserGrantProjection{}
	config.ProjectionName = UserGrantProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *UserGrantProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  usergrant.UserGrantAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  usergrant.UserGrantChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  usergrant.UserGrantCascadeChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  usergrant.UserGrantRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  usergrant.UserGrantCascadeRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  usergrant.UserGrantDeactivatedType,
					Reduce: p.reduceDeactivated,
				},
				{
					Event:  usergrant.UserGrantReactivatedType,
					Reduce: p.reduceReactivated,
				},
			},
		},
	}
}

type UserGrantColumn string

const (
	UserGrantID            = "id"
	UserGrantResourceOwner = "resource_owner"
	UserGrantCreationDate  = "creation_date"
	UserGrantChangeDate    = "change_date"
	UserGrantSequence      = "sequence"
	UserGrantUserID        = "user_id"
	UserGrantProjectID     = "project_id"
	UserGrantGrantID       = "grant_id"
	UserGrantRoles         = "roles"
	UserGrantState         = "state"
)

func (p *UserGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantAddedEvent)
	if !ok {
		logging.LogWithFields("PROJE-WYOHD", "seq", event.Sequence(), "expectedType", usergrant.UserGrantAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-MQHVB", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserGrantID, e.Aggregate().ID),
			handler.NewCol(UserGrantResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserGrantCreationDate, e.CreationDate()),
			handler.NewCol(UserGrantChangeDate, e.CreationDate()),
			handler.NewCol(UserGrantSequence, e.Sequence()),
			handler.NewCol(UserGrantUserID, e.UserID),
			handler.NewCol(UserGrantProjectID, e.ProjectID),
			handler.NewCol(UserGrantGrantID, e.ProjectGrantID),
			handler.NewCol(UserGrantRoles, pq.StringArray(e.RoleKeys)),
			handler.NewCol(UserGrantState, domain.UserGrantStateActive),
		},
	), nil
}

func (p *UserGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles pq.StringArray

	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		roles = e.RoleKeys
	case *usergrant.UserGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		logging.LogWithFields("PROJE-dIflx", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-hOr1E", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantRoles, roles),
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *UserGrantProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *usergrant.UserGrantRemovedEvent, *usergrant.UserGrantCascadeRemovedEvent:
		// ok
	default:
		logging.LogWithFields("PROJE-Nw0cR", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{usergrant.UserGrantRemovedType, usergrant.UserGrantCascadeRemovedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-7OBEC", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *UserGrantProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*usergrant.UserGrantDeactivatedEvent); !ok {
		logging.LogWithFields("PROJE-V3txf", "seq", event.Sequence(), "expectedType", usergrant.UserGrantDeactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-oP7Gm", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantState, domain.UserGrantStateInactive),
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *UserGrantProjection) reduceReactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*usergrant.UserGrantDeactivatedEvent); !ok {
		logging.LogWithFields("PROJE-ly6oe", "seq", event.Sequence(), "expectedType", usergrant.UserGrantReactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-DGsKh", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantState, domain.UserGrantStateActive),
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}
