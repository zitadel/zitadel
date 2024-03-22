package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserGrantProjectionTable = "projections.user_grants5"

	UserGrantID                   = "id"
	UserGrantCreationDate         = "creation_date"
	UserGrantChangeDate           = "change_date"
	UserGrantSequence             = "sequence"
	UserGrantState                = "state"
	UserGrantResourceOwner        = "resource_owner"
	UserGrantInstanceID           = "instance_id"
	UserGrantUserID               = "user_id"
	UserGrantResourceOwnerUser    = "resource_owner_user"
	UserGrantProjectID            = "project_id"
	UserGrantResourceOwnerProject = "resource_owner_project"
	UserGrantGrantID              = "grant_id"
	UserGrantGrantedOrg           = "granted_org"
	UserGrantRoles                = "roles"
)

type userGrantProjection struct {
	es handler.EventStore
}

func newUserGrantProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &userGrantProjection{es: config.Eventstore})
}

func (*userGrantProjection) Name() string {
	return UserGrantProjectionTable
}

func (*userGrantProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(UserGrantID, handler.ColumnTypeText),
			handler.NewColumn(UserGrantCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserGrantChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserGrantSequence, handler.ColumnTypeInt64),
			handler.NewColumn(UserGrantState, handler.ColumnTypeEnum),
			handler.NewColumn(UserGrantResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(UserGrantInstanceID, handler.ColumnTypeText),
			handler.NewColumn(UserGrantUserID, handler.ColumnTypeText),
			handler.NewColumn(UserGrantResourceOwnerUser, handler.ColumnTypeText),
			handler.NewColumn(UserGrantProjectID, handler.ColumnTypeText),
			handler.NewColumn(UserGrantResourceOwnerProject, handler.ColumnTypeText),
			handler.NewColumn(UserGrantGrantID, handler.ColumnTypeText),
			handler.NewColumn(UserGrantGrantedOrg, handler.ColumnTypeText),
			handler.NewColumn(UserGrantRoles, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(UserGrantInstanceID, UserGrantID),
			handler.WithIndex(handler.NewIndex("user_id", []string{UserGrantUserID})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{UserGrantResourceOwner})),
		),
	)
}

func (p *userGrantProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: usergrant.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-MQHVB", "reduce.wrong.event.type %s", usergrant.UserGrantAddedType)
	}

	ctx := setUserGrantContext(e.Aggregate())
	userOwner, projectOwner, grantOwner, err := getUserGrantResourceOwners(ctx, p.es, e.Aggregate().InstanceID, e.UserID, e.ProjectID, e.ProjectGrantID)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserGrantID, e.Aggregate().ID),
			handler.NewCol(UserGrantResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserGrantInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(UserGrantCreationDate, e.CreatedAt()),
			handler.NewCol(UserGrantChangeDate, e.CreatedAt()),
			handler.NewCol(UserGrantSequence, e.Sequence()),
			handler.NewCol(UserGrantUserID, e.UserID),
			handler.NewCol(UserGrantResourceOwnerUser, userOwner),
			handler.NewCol(UserGrantProjectID, e.ProjectID),
			handler.NewCol(UserGrantResourceOwnerProject, projectOwner),
			handler.NewCol(UserGrantGrantID, e.ProjectGrantID),
			handler.NewCol(UserGrantGrantedOrg, grantOwner),
			handler.NewCol(UserGrantRoles, database.TextArray[string](e.RoleKeys)),
			handler.NewCol(UserGrantState, domain.UserGrantStateActive),
		},
	), nil
}

func (p *userGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles database.TextArray[string]

	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		roles = e.RoleKeys
	case *usergrant.UserGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-hOr1E", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreatedAt()),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-7OBEC", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantRemovedType, usergrant.UserGrantCascadeRemovedType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userGrantProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*usergrant.UserGrantDeactivatedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-oP7Gm", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreatedAt()),
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
	if _, ok := event.(*usergrant.UserGrantReactivatedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DGsKh", "reduce.wrong.event.type %s", usergrant.UserGrantReactivatedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserGrantChangeDate, event.CreatedAt()),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bner2a", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserGrantUserID, event.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userGrantProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*project.ProjectRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bne2a", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}

	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dGr2a", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}

	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dswg2", "reduce.wrong.event.type %s", project.RoleRemovedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewArrayRemoveCol(UserGrantRoles, e.Key),
		},
		[]handler.Condition{
			handler.NewCond(UserGrantProjectID, e.Aggregate().ID),
			handler.NewCond(UserGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	var grantID string
	var keys database.TextArray[string]
	switch e := event.(type) {
	case *project.GrantChangedEvent:
		grantID = e.GrantID
		keys = e.RoleKeys
	case *project.GrantCascadeChangedEvent:
		grantID = e.GrantID
		keys = e.RoleKeys
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Fh3gw", "reduce.wrong.event.type %v", []eventstore.EventType{project.GrantChangedType, project.GrantCascadeChangedType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewArrayIntersectCol(UserGrantRoles, keys),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-jpIvp", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwner, e.Aggregate().ID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwnerUser, e.Aggregate().ID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantResourceOwnerProject, e.Aggregate().ID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(UserGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(UserGrantGrantedOrg, e.Aggregate().ID),
			},
		),
	), nil
}

func getUserResourceOwner(ctx context.Context, es handler.EventStore, instanceID, userID string) (string, error) {
	userRO, _, _, err := getResourceOwners(ctx, es, instanceID, userID, "", "")
	if userRO == "" {
		return "", zerrors.ThrowNotFound(nil, "PROJ-uahkkord22", "Errors.NotFound")
	}
	return userRO, err
}

func getUserGrantResourceOwners(ctx context.Context, es handler.EventStore, instanceID, userID, projectID, grantID string) (string, string, string, error) {
	userRO, projectRO, grantedOrg, err := getResourceOwners(ctx, es, instanceID, userID, projectID, grantID)
	if err != nil {
		return "", "", "", err
	}
	// user grant always has a user defined
	if userRO == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-8x5behx5jy", "Errors.NotFound")
	}
	// either a projectID
	if projectID != "" && projectRO == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-1ldp25o3bx", "Errors.NotFound")
	}
	// or a grantID
	if grantID != "" && grantedOrg == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-9ngp5dcn76", "Errors.NotFound")
	}
	return userRO, projectRO, grantedOrg, nil
}

func getResourceOwners(ctx context.Context, es handler.EventStore, instanceID, userID, projectID, grantID string) (userRO string, projectRO string, grantedOrg string, err error) {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		InstanceID(instanceID).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		EventTypes(user.HumanRegisteredType, user.HumanAddedType, user.MachineAddedEventType)

	// if it's a project grant then we only need the resourceowner for the projectgrant, else the project
	if grantID != "" {
		builder = builder.Or().
			AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventTypes(project.GrantAddedType).
			EventData(map[string]interface{}{
				"grantId": grantID,
			})
	}
	if projectID != "" {
		builder = builder.Or().
			AggregateTypes(project.AggregateType).
			AggregateIDs(projectID).
			EventTypes(project.ProjectAddedType)
	}

	events, err := es.Filter(
		ctx,
		builder.Builder(),
	)
	if err != nil {
		return "", "", "", err
	}

	// sorted ascending
	for _, event := range events {
		switch e := event.(type) {
		case *project.GrantAddedEvent:
			grantedOrg = e.GrantedOrgID
		case *project.ProjectAddedEvent:
			projectRO = e.Aggregate().ResourceOwner
		case *user.HumanRegisteredEvent, *user.HumanAddedEvent, *user.MachineAddedEvent:
			userRO = e.Aggregate().ResourceOwner
		}
	}
	return userRO, projectRO, grantedOrg, nil
}

func setUserGrantContext(aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), aggregate.InstanceID)
}
