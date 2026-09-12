package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupGrantProjectionTable = "projections.group_grants1"

	GroupGrantID            = "id"
	GroupGrantCreationDate  = "creation_date"
	GroupGrantChangeDate    = "change_date"
	GroupGrantSequence      = "sequence"
	GroupGrantResourceOwner = "resource_owner"
	GroupGrantInstanceID    = "instance_id"
	GroupGrantGroupID       = "group_id"
	GroupGrantProjectID     = "project_id"
	GroupGrantGrantID       = "grant_id"
	GroupGrantRoles         = "roles"
)

type groupGrantProjection struct{}

func newGroupGrantProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupGrantProjection))
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
			handler.NewColumn(GroupGrantResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantProjectID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantGrantID, handler.ColumnTypeText),
			handler.NewColumn(GroupGrantRoles, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GroupGrantInstanceID, GroupGrantID),
			handler.WithIndex(handler.NewIndex("group_id", []string{GroupGrantGroupID})),
			handler.WithIndex(handler.NewIndex("project_id", []string{GroupGrantProjectID})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{GroupGrantResourceOwner})),
		),
	)
}

func (p *groupGrantProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: groupgrant.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  groupgrant.GroupGrantAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  groupgrant.GroupGrantChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  groupgrant.GroupGrantRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  groupgrant.GroupGrantCascadeRemovedType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupRemovedEventType,
					Reduce: p.reduceGroupRemoved,
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
					Reduce: reduceInstanceRemovedHelper(GroupGrantInstanceID),
				},
			},
		},
	}
}

func (p *groupGrantProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupgrant.GroupGrantAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-h0J9Wq", "reduce.wrong.event.type %s", groupgrant.GroupGrantAddedType)
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
			handler.NewCol(GroupGrantGrantID, e.ProjectGrantID),
			handler.NewCol(GroupGrantRoles, database.TextArray[string](e.RoleKeys)),
		},
	), nil
}

func (p *groupGrantProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupgrant.GroupGrantChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pLx2Vn", "reduce.wrong.event.type %s", groupgrant.GroupGrantChangedType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupGrantChangeDate, e.CreatedAt()),
			handler.NewCol(GroupGrantRoles, database.TextArray[string](e.RoleKeys)),
			handler.NewCol(GroupGrantSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(GroupGrantID, e.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupGrantProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *groupgrant.GroupGrantRemovedEvent, *groupgrant.GroupGrantCascadeRemovedEvent:
		// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-V3xJqf", "reduce.wrong.event.type %v", []eventstore.EventType{groupgrant.GroupGrantRemovedType, groupgrant.GroupGrantCascadeRemovedType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupGrantProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*group.GroupRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-q8WtVl", "reduce.wrong.event.type %s", group.GroupRemovedEventType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantGroupID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupGrantProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*project.ProjectRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-jYp3Wd", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantProjectID, event.Aggregate().ID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupGrantProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-mN4Xqk", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupGrantGrantID, e.GrantID),
			handler.NewCond(GroupGrantInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupGrantProjection) reduceRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-tB7Hkz", "reduce.wrong.event.type %s", project.RoleRemovedType)
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

func (p *groupGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-x9HwLs", "reduce.wrong.event.type %v", []eventstore.EventType{project.GrantChangedType, project.GrantCascadeChangedType})
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

func (p *groupGrantProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-w2LqRv", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupGrantInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupGrantResourceOwner, e.Aggregate().ID),
		},
	), nil
}
