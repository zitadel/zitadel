package projection

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	v3 "github.com/zitadel/zitadel/internal/eventstore/handler/v3"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

const (
	UserGrantProjectionTable = "projections.user_grants2"

	UserGrantID            = "id"
	UserGrantCreationDate  = "creation_date"
	UserGrantChangeDate    = "change_date"
	UserGrantState         = "state"
	UserGrantResourceOwner = "resource_owner"
	UserGrantInstanceID    = "instance_id"
	UserGrantUserID        = "user_id"
	UserGrantProjectID     = "project_id"
	UserGrantGrantID       = "grant_id"
	UserGrantRoles         = "roles"
)

type userGrantProjection struct{}

func newUserGrantProjection(ctx context.Context, config v3.Config) *v3.IDProjection {
	p := new(userGrantProjection)

	config.Check = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserGrantID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserGrantChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserGrantState, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserGrantResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantUserID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantProjectID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantGrantID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantRoles, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(UserGrantInstanceID, UserGrantID),
			crdb.WithIndex(crdb.NewIndex("user_grant_user_idx", []string{UserGrantUserID})),
			crdb.WithIndex(crdb.NewIndex("user_grant_ro_idx", []string{UserGrantResourceOwner})),
		),
	)
	config.Reduces = map[eventstore.AggregateType][]v3.Reducer{
		usergrant.AggregateType: {
			{
				Event:  usergrant.UserGrantAddedType,
				Reduce: p.reduceAdded,
			},
			{
				Event:          usergrant.UserGrantChangedType,
				Reduce:         p.reduceChanged,
				PreviousEvents: p.previousEventsChanged,
			},
			{
				Event:          usergrant.UserGrantCascadeChangedType,
				Reduce:         p.reduceChanged,
				PreviousEvents: p.previousEventsChanged,
			},
			{
				Event:          usergrant.UserGrantRemovedType,
				Reduce:         p.reduceRemoved,
				PreviousEvents: p.previousEventsRemoved,
			},
			{
				Event:          usergrant.UserGrantCascadeRemovedType,
				Reduce:         p.reduceRemoved,
				PreviousEvents: p.previousEventsRemoved,
			},
			{
				Event:          usergrant.UserGrantDeactivatedType,
				Reduce:         p.reduceDeactivated,
				PreviousEvents: p.previousEventsDeactivated,
			},
			{
				Event:          usergrant.UserGrantReactivatedType,
				Reduce:         p.reduceReactivated,
				PreviousEvents: p.previousEventsReactivated,
			},
		},
		user.AggregateType: {
			{
				Event:          user.UserRemovedType,
				Reduce:         p.reduceUserRemoved,
				PreviousEvents: p.previousEventsUser,
			},
		},
		project.AggregateType: {
			{
				Event:          project.ProjectRemovedType,
				Reduce:         p.reduceProjectRemoved,
				PreviousEvents: p.previousEventsProject,
			},
			{
				Event:          project.GrantRemovedType,
				Reduce:         p.reduceProjectGrantRemoved,
				PreviousEvents: p.previousEventsProject,
			},
			{
				Event:          project.RoleRemovedType,
				Reduce:         p.reduceRoleRemoved,
				PreviousEvents: p.previousEventsProject,
			},
			{
				Event:          project.GrantChangedType,
				Reduce:         p.reduceProjectGrantChanged,
				PreviousEvents: p.previousEventsProject,
			},
			{
				Event:          project.GrantCascadeChangedType,
				Reduce:         p.reduceProjectGrantChanged,
				PreviousEvents: p.previousEventsProject,
			},
		},
	}

	return v3.StartSubscriptionIDProjection(ctx, UserGrantProjectionTable, config)
}

func (p *userGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-MQHVB", "reduce.wrong.event.type %s", usergrant.UserGrantAddedType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserGrantID, e.Aggregate().ID),
			handler.NewCol(UserGrantResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserGrantInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(UserGrantCreationDate, e.CreationDate()),
			handler.NewCol(UserGrantChangeDate, e.CreationDate()),
			handler.NewCol(UserGrantUserID, e.UserID),
			handler.NewCol(UserGrantProjectID, e.ProjectID),
			handler.NewCol(UserGrantGrantID, e.ProjectGrantID),
			handler.NewCol(UserGrantRoles, database.StringArray(e.RoleKeys)),
			handler.NewCol(UserGrantState, domain.UserGrantStateActive),
		},
	), nil
}

func (p *userGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles database.StringArray

	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		roles = e.RoleKeys
	case *usergrant.UserGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-hOr1E", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantRoles, roles),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) previousEventsChanged(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	var grantID string

	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		grantID = e.Aggregate().ID
	case *usergrant.UserGrantCascadeChangedEvent:
		grantID = e.Aggregate().ID
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-RePmo", "previous.events.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return p.previousEvents(grantID, event.Aggregate().InstanceID, tx)
}

func (p *userGrantProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *usergrant.UserGrantRemovedEvent, *usergrant.UserGrantCascadeRemovedEvent:
		// ok
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-7OBEC", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantRemovedType, usergrant.UserGrantCascadeRemovedType})
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) previousEventsRemoved(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	var grantID string

	switch e := event.(type) {
	case *usergrant.UserGrantRemovedEvent:
		grantID = e.Aggregate().ID
	case *usergrant.UserGrantCascadeRemovedEvent:
		grantID = e.Aggregate().ID
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-RePmo", "previous.events.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return p.previousEvents(grantID, event.Aggregate().InstanceID, tx)
}

func (p *userGrantProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*usergrant.UserGrantDeactivatedEvent); !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-oP7Gm", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantState, domain.UserGrantStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) previousEventsDeactivated(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*usergrant.UserGrantDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-oP7Gm", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}

	return p.previousEvents(e.Aggregate().ID, event.Aggregate().InstanceID, tx)
}

func (p *userGrantProjection) reduceReactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*usergrant.UserGrantDeactivatedEvent); !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-DGsKh", "reduce.wrong.event.type %s", usergrant.UserGrantReactivatedType)
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreationDate()),
			handler.NewCol(UserGrantState, domain.UserGrantStateActive),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
		},
	), nil
}

func (p *userGrantProjection) previousEventsReactivated(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*usergrant.UserGrantReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-WEpat", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}

	return p.previousEvents(e.Aggregate().ID, event.Aggregate().InstanceID, tx)
}

func (p *userGrantProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*user.UserRemovedEvent); !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Bner2a", "reduce.wrong.event.type %s", user.UserRemovedType)
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Bne2a", "reduce.wrong.event.type %s", project.ProjectRemovedType)
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-dGr2a", "reduce.wrong.event.type %s", project.GrantRemovedType)
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-dswg2", "reduce.wrong.event.type %s", project.RoleRemovedType)
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Fh3gw", "reduce.wrong.event.type %v", []eventstore.EventType{project.GrantChangedType, project.GrantCascadeChangedType})
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			crdb.NewArrayIntersectCol(UserGrantRoles, database.StringArray(keys)),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantGrantID, grantID),
		},
	), nil
}

// previousEventsUserRemoved locks the rows of the user
func (p *userGrantProjection) previousEventsUser(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	_, err := tx.Exec("SELECT FOR UPDATE 1 FROM "+UserGrantProjectionTable+" WHERE "+UserGrantUserID+" = $1 AND "+UserGrantInstanceID+" = $2", event.Aggregate().ID, event.Aggregate().InstanceID)
	return nil, err
}

// previousEventsUserRemoved locks the rows of the project
func (p *userGrantProjection) previousEventsProject(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	_, err := tx.Exec("SELECT FOR UPDATE 1 FROM "+UserGrantProjectionTable+" WHERE "+UserGrantProjectID+" = $1 AND "+UserGrantInstanceID+" = $2", event.Aggregate().ID, event.Aggregate().InstanceID)
	return nil, err
}

// previousEventsUserRemoved locks the rows of the grant
func (p *userGrantProjection) previousEventsProjectGrant(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-jTfdc", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}

	_, err := tx.Exec("SELECT FOR UPDATE 1 FROM "+UserGrantProjectionTable+" WHERE "+UserGrantGrantID+" = $1 AND "+UserGrantInstanceID+" = $2", e.GrantID, e.Aggregate().InstanceID)
	return nil, err
}

func (p *userGrantProjection) previousEvents(grantID, instanceID string, tx *sql.Tx) (*eventstore.SearchQueryBuilder, error) {
	row := tx.QueryRow("SELECT FOR UPDATE "+UserGrantChangeDate+" FROM "+UserGrantProjectionTable+" WHERE "+UserGrantID+" = $1 AND "+UserGrantInstanceID+" = $2", grantID, instanceID)

	var changeDate time.Time

	if err := row.Scan(changeDate); err != nil {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(instanceID).
		SystemTime(changeDate).
		AddQuery().
		AggregateTypes(usergrant.AggregateType).
		AggregateIDs(grantID).
		CreationDateAfter(changeDate).
		Builder(), nil
}
