package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

const (
	UserGrantProjectionTable = "projections.user_grants3"

	UserGrantID                   = "id"
	UserGrantCreationDate         = "creation_date"
	UserGrantChangeDate           = "change_date"
	UserGrantSequence             = "sequence"
	UserGrantState                = "state"
	UserGrantResourceOwner        = "resource_owner"
	UserGrantInstanceID           = "instance_id"
	UserGrantUserID               = "user_id"
	UserGrantResourceOwnerUser    = "resource_owner_user"
	UserGrantUserOwnerRemoved     = "user_owner_removed"
	UserGrantProjectID            = "project_id"
	UserGrantResourceOwnerProject = "resource_owner_project"
	UserGrantProjectOwnerRemoved  = "project_owner_removed"
	UserGrantGrantID              = "grant_id"
	UserGrantGrantedOrg           = "granted_org"
	UserGrantGrantedOrgRemoved    = "granted_org_removed"
	UserGrantRoles                = "roles"
	UserGrantOwnerRemoved         = "owner_removed"
)

type userGrantProjection struct {
	crdb.StatementHandler
}

func newUserGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *userGrantProjection {
	p := new(userGrantProjection)
	config.ProjectionName = UserGrantProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserGrantID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserGrantChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserGrantSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(UserGrantState, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserGrantResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantUserID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantResourceOwnerUser, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantUserOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(UserGrantProjectID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantResourceOwnerProject, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantProjectOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(UserGrantGrantID, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantGrantedOrg, crdb.ColumnTypeText),
			crdb.NewColumn(UserGrantGrantedOrgRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(UserGrantRoles, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(UserGrantOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(UserGrantInstanceID, UserGrantID),
			crdb.WithIndex(crdb.NewIndex("user_id", []string{UserGrantUserID})),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{UserGrantResourceOwner})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{UserGrantOwnerRemoved})),
			crdb.WithIndex(crdb.NewIndex("user_owner_removed", []string{UserGrantUserOwnerRemoved})),
			crdb.WithIndex(crdb.NewIndex("project_owner_removed", []string{UserGrantProjectOwnerRemoved})),
			crdb.WithIndex(crdb.NewIndex("granted_org_removed", []string{UserGrantGrantedOrgRemoved})),
		),
	)

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
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(UserGrantInstanceID),
				},
			},
		},
	}
}

func (p *userGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-MQHVB", "reduce.wrong.event.type %s", usergrant.UserGrantAddedType)
	}

	ctx := setUserGrantContext(e.Aggregate())
	userOwner, err := getResourceOwnerOfUser(ctx, p.Eventstore, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}

	projectOwner := ""
	grantOwner := ""
	if e.ProjectGrantID != "" {
		grantOwner, err = getGrantedOrgOfGrantedProject(ctx, p.Eventstore, e.Aggregate().InstanceID, e.ProjectID, e.ProjectGrantID)
		if err != nil {
			return nil, err
		}
	} else {
		projectOwner, err = getResourceOwnerOfProject(ctx, p.Eventstore, e.Aggregate().InstanceID, e.ProjectID)
		if err != nil {
			return nil, err
		}
	}

	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserGrantID, e.Aggregate().ID),
			handler.NewCol(UserGrantResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserGrantInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(UserGrantCreationDate, e.CreationDate()),
			handler.NewCol(UserGrantChangeDate, e.CreationDate()),
			handler.NewCol(UserGrantSequence, e.Sequence()),
			handler.NewCol(UserGrantUserID, e.UserID),
			handler.NewCol(UserGrantResourceOwnerUser, userOwner),
			handler.NewCol(UserGrantProjectID, e.ProjectID),
			handler.NewCol(UserGrantResourceOwnerProject, projectOwner),
			handler.NewCol(UserGrantGrantID, e.ProjectGrantID),
			handler.NewCol(UserGrantGrantedOrg, grantOwner),
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
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
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
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
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
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
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
			handler.NewCol(UserGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userGrantProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*user.UserRemovedEvent); !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Bner2a", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantUserID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
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
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
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
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
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
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
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
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userGrantProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-jpIvp", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwner, e.Aggregate().ID),
			},
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwnerUser, e.Aggregate().ID),
			},
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwnerProject, e.Aggregate().ID),
			},
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantGrantedOrg, e.Aggregate().ID),
			},
		),
	), nil
}

func getResourceOwnerOfUser(ctx context.Context, es *eventstore.Eventstore, instanceID, aggID string) (string, error) {
	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(instanceID).
			AddQuery().
			AggregateTypes(user.AggregateType).
			AggregateIDs(aggID).
			EventTypes(user.HumanRegisteredType, user.HumanAddedType, user.MachineAddedEventType).
			Builder(),
	)
	if err != nil {
		return "", err
	}
	if len(events) != 1 {
		return "", errors.ThrowNotFound(nil, "PROJ-0I92sp", "Errors.User.NotFound")
	}
	return events[0].Aggregate().ResourceOwner, nil
}

func getResourceOwnerOfProject(ctx context.Context, es *eventstore.Eventstore, instanceID, aggID string) (string, error) {
	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(instanceID).
			AddQuery().
			AggregateTypes(project.AggregateType).
			AggregateIDs(aggID).
			EventTypes(project.ProjectAddedType).
			Builder(),
	)
	if err != nil {
		return "", err
	}
	if len(events) != 1 {
		return "", errors.ThrowNotFound(nil, "PROJ-0I91sp", "Errors.Project.NotFound")
	}
	return events[0].Aggregate().ResourceOwner, nil
}

func getGrantedOrgOfGrantedProject(ctx context.Context, es *eventstore.Eventstore, instanceID, projectID, grantID string) (string, error) {
	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(instanceID).
			AddQuery().
			AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventTypes(project.GrantAddedType).
			EventData(map[string]interface{}{
				"grantId": grantID,
			}).
			Builder(),
	)
	if err != nil {
		return "", err
	}
	if len(events) != 1 {
		return "", errors.ThrowNotFound(nil, "PROJ-MoaSpw", "Errors.Grant.NotFound")
	}
	grantAddedEvent, ok := events[0].(*project.GrantAddedEvent)
	if !ok {
		return "", errors.ThrowNotFound(nil, "PROJ-P0s2o0", "Errors.Grant.NotFound")
	}
	return grantAddedEvent.GrantedOrgID, nil
}

func setUserGrantContext(event eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), event.InstanceID)
}
