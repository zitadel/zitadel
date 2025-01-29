package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupGrantProjectionTable = "projections.group_grants16"

	GroupGrantCreationDate  = "creation_date"
	GroupGrantChangeDate    = "change_date"
	GroupGrantSequence      = "sequence"
	GroupGrantResourceOwner = "resource_owner"
	GroupGrantInstanceID    = "instance_id"
	// GroupGrantResourceOwnerUser    = "resource_owner_user" // Why -- ? Not needed as per current understanding.

	GroupGrantID                   = "id"
	GroupGrantState                = "state"
	GroupGrantGroupID              = "group_id"
	GroupGrantProjectID            = "project_id"
	GroupGrantResourceOwnerProject = "resource_owner_project"
	GroupGrantGrantID              = "grant_id"
	GroupGrantGrantedOrg           = "granted_org"
	GroupGrantRoles                = "roles"
)

type groupGrantProjection struct {
	es handler.EventStore
}

func newGroupGrantProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &groupGrantProjection{es: config.Eventstore})
}

func (*groupGrantProjection) Name() string {
	return GroupGrantProjectionTable
}

func (*groupGrantProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupGrantID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupGrantChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupGrantSequence, handler.ColumnTypeInt64),
			handler.NewColumn(GroupGrantState, handler.ColumnTypeEnum),
			handler.NewColumn(GroupGrantResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantProjectID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantResourceOwnerProject, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantGrantID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantGrantedOrg, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantRoles, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GroupGrantInstanceID, GroupGrantID),
			handler.WithIndex(handler.NewIndex("group_id", []string{GroupGrantGroupID})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{UserGrantResourceOwner})),
		),
	)
}

func (g *groupGrantProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: groupgrant.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  groupgrant.GroupGrantAddedType,
					Reduce: g.reduceAdded,
				},
				{
					Event:  groupgrant.GroupGrantChangedType,
					Reduce: g.reduceChanged,
				},
				{
					Event:  groupgrant.GroupGrantCascadeChangedType,
					Reduce: g.reduceChanged,
				},
				{
					Event:  groupgrant.GroupGrantRemovedType,
					Reduce: g.reduceRemoved,
				},
				{
					Event:  groupgrant.GroupGrantCascadeRemovedType,
					Reduce: g.reduceRemoved,
				},
				{
					Event:  groupgrant.GroupGrantDeactivatedType,
					Reduce: g.reduceDeactivated,
				},
				{
					Event:  groupgrant.GroupGrantReactivatedType,
					Reduce: g.reduceReactivated,
				},
			},
		},
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupRemovedType,
					Reduce: g.reduceGroupRemoved,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectRemovedType,
					Reduce: g.reduceProjectRemoved,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: g.reduceProjectGrantRemoved,
				},
				{
					Event:  project.RoleRemovedType,
					Reduce: g.reduceRoleRemoved,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: g.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: g.reduceProjectGrantChanged,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: g.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(GroupGrantInstanceID),
				},
			},
		},
	}
}

func (g *groupGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupgrant.GroupGrantAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "GGRANT-NPGVB", "reduce.wrong.event.type %s", groupgrant.GroupGrantAddedType)
	}

	ctx := setGroupGrantContext(e.Aggregate())
	_, projectOwner, grantOwner, err := getGroupGrantResourceOwners(ctx, g.es, e.Aggregate().InstanceID, e.GroupID, e.ProjectID, e.ProjectGrantID)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupGrantID, e.Aggregate().ID),
			handler.NewCol(GroupGrantResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupGrantInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupGrantCreationDate, e.CreatedAt()),
			handler.NewCol(GroupGrantChangeDate, e.CreatedAt()),
			handler.NewCol(GroupGrantSequence, e.Sequence()),
			handler.NewCol(GroupGrantGroupID, e.GroupID),
			handler.NewCol(GroupGrantProjectID, e.ProjectID),
			handler.NewCol(GroupGrantResourceOwnerProject, projectOwner),
			handler.NewCol(GroupGrantGrantID, e.ProjectGrantID),
			handler.NewCol(GroupGrantGrantedOrg, grantOwner),
			handler.NewCol(GroupGrantRoles, database.TextArray[string](e.RoleKeys)),
			handler.NewCol(GroupGrantState, domain.GroupGrantStateActive),
		},
	), nil
}

func (g *groupGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles database.TextArray[string]

	switch e := event.(type) {
	case *groupgrant.GroupGrantChangedEvent:
		roles = e.RoleKeys
	case *groupgrant.GroupGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-gpR1E", "reduce.wrong.event.type %v", []eventstore.EventType{groupgrant.GroupGrantChangedType, groupgrant.GroupGrantCascadeChangedType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(GroupGrantChangeDate, event.CreatedAt()),
			handler.NewCol(GroupGrantRoles, roles),
			handler.NewCol(GroupGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *groupgrant.GroupGrantRemovedEvent, *groupgrant.GroupGrantCascadeRemovedEvent:
		// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-8PCEC", "reduce.wrong.event.type %v", []eventstore.EventType{groupgrant.GroupGrantRemovedType, groupgrant.GroupGrantCascadeRemovedType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceDeactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*groupgrant.GroupGrantDeactivatedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pQ8Gm", "reduce.wrong.event.type %s", groupgrant.GroupGrantDeactivatedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(GroupGrantChangeDate, event.CreatedAt()),
			handler.NewCol(GroupGrantState, domain.GroupGrantStateInactive),
			handler.NewCol(GroupGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceReactivated(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*groupgrant.GroupGrantReactivatedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EFtLh", "reduce.wrong.event.type %s", groupgrant.GroupGrantReactivatedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(GroupGrantChangeDate, event.CreatedAt()),
			handler.NewCol(GroupGrantState, domain.GroupGrantStateActive),
			handler.NewCol(GroupGrantSequence, event.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*group.GroupRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CNfr3a", "reduce.wrong.event.type %s", group.GroupRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantGroupID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*project.ProjectRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Cne2a", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantProjectID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-cGr2a", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantGrantID, e.GrantID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-dswg3", "reduce.wrong.event.type %s", project.RoleRemovedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewArrayRemoveCol(GroupGrantRoles, e.Key),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantProjectID, e.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Fh3gx", "reduce.wrong.event.type %v", []eventstore.EventType{project.GrantChangedType, project.GrantCascadeChangedType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewArrayIntersectCol(GroupGrantRoles, keys),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantGrantID, grantID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupGrantProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-jpIvq", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(GroupGrantResourceOwner, e.Aggregate().ID),
			},
		),
		// handler.AddDeleteStatement(
		// 	[]handler.Condition{
		// 		handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
		// 		handler.NewCond(GroupGrantResourceOwnerUser, e.Aggregate().ID),
		// 	},
		// ),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(GroupGrantResourceOwnerProject, e.Aggregate().ID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(GroupGrantGrantedOrg, e.Aggregate().ID),
			},
		),
	), nil
}

func getGroupResourceOwner(ctx context.Context, es handler.EventStore, instanceID, groupID string) (string, error) {
	groupRO, _, _, err := getGroupResourceOwners(ctx, es, instanceID, groupID, "", "")
	if groupRO == "" {
		return "", zerrors.ThrowNotFound(nil, "PROJ-uahkkord22", "Errors.NotFound")
	}
	return groupRO, err
}

func getGroupGrantResourceOwners(ctx context.Context, es handler.EventStore, instanceID, groupID, projectID, grantID string) (string, string, string, error) {
	groupRO, projectRO, grantedOrg, err := getGroupResourceOwners(ctx, es, instanceID, groupID, projectID, grantID)
	if err != nil {
		return "", "", "", err
	}
	// group grant always has a user defined
	if groupRO == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-9y6behx5ky", "Errors.NotFound")
	}
	// either a projectID
	if projectID != "" && projectRO == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-2pep36o3by", "Errors.NotFound")
	}
	// or a grantID
	if grantID != "" && grantedOrg == "" {
		return "", "", "", zerrors.ThrowNotFound(nil, "PROJ-1mGq5dcn86", "Errors.NotFound")
	}
	return groupRO, projectRO, grantedOrg, nil
}

func getGroupResourceOwners(ctx context.Context, es handler.EventStore, instanceID, groupID, projectID, grantID string) (groupRO string, projectRO string, grantedOrg string, err error) {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		InstanceID(instanceID).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(groupID).
		EventTypes(group.GroupAddedType)

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
		case *group.GroupAddedEvent:
			groupRO = e.Aggregate().ResourceOwner
		}
	}
	return groupRO, projectRO, grantedOrg, nil
}

func setGroupGrantContext(aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), aggregate.InstanceID)
}
