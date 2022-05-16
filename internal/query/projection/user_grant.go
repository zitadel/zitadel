package projection

import (
	"context"

	"github.com/lib/pq"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type userGrantProjection struct {
	crdb.StatementHandler
}

const (
	UserGrantProjectionTable = "zitadel.projections.user_grants"
)

func newUserGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *userGrantProjection {
	p := &userGrantProjection{}
	config.ProjectionName = UserGrantProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *userGrantProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: usergrant.AggregateType,
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
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
				},
				{
					Event:  project.RoleRemovedType,
					Reduce: p.reduceRoleRemoved,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: p.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: p.reduceProjectGrantChanged,
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

func (p *userGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userGrantProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userGrantProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userGrantProjection) reduceReactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userGrantProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*user.UserRemovedEvent); !ok {
		logging.LogWithFields("PROJE-Vfeg3", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Bner2a", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantUserID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*project.ProjectRemovedEvent); !ok {
		logging.LogWithFields("PROJE-Vfeg3", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Bne2a", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantProjectID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-DGfe2", "seq", event.Sequence(), "expectedType", project.GrantRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-dGr2a", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantGrantID, e.GrantID),
		},
	), nil
}

func (p *userGrantProjection) reduceRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-Edg22", "seq", event.Sequence(), "expectedType", project.RoleRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-dswg2", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			crdb.NewArrayRemoveCol(UserGrantRoles, e.Key),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	var grantID string
	var keys []string
	switch e := event.(type) {
	case *project.GrantChangedEvent:
		grantID = e.GrantID
		keys = e.RoleKeys
	case *project.GrantCascadeChangedEvent:
		grantID = e.GrantID
		keys = e.RoleKeys
	default:
		logging.LogWithFields("PROJE-FGgw2", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.GrantChangedType, project.GrantCascadeChangedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Fh3gw", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			crdb.NewArrayIntersectCol(UserGrantRoles, pq.StringArray(keys)),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantGrantID, grantID),
		},
	), nil
}
