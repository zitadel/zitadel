package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupManagerRolesProjectionTable = "projections.group_manager_roles1"

	GroupManagerRolesGroupID       = "group_id"
	GroupManagerRolesResourceOwner = "resource_owner"
	GroupManagerRolesInstanceID    = "instance_id"
	GroupManagerRolesRoles         = "roles"
	GroupManagerRolesCreationDate  = "creation_date"
	GroupManagerRolesChangeDate    = "change_date"
	GroupManagerRolesSequence      = "sequence"
)

type groupManagerRolesProjection struct{}

func newGroupManagerRolesProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupManagerRolesProjection))
}

func (*groupManagerRolesProjection) Name() string {
	return GroupManagerRolesProjectionTable
}

func (*groupManagerRolesProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupManagerRolesGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupManagerRolesResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupManagerRolesInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupManagerRolesRoles, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(GroupManagerRolesCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupManagerRolesChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupManagerRolesSequence, handler.ColumnTypeInt64),
		},
			handler.NewPrimaryKey(GroupManagerRolesInstanceID, GroupManagerRolesGroupID),
		),
	)
}

func (p *groupManagerRolesProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupManagerRolesSetEventType,
					Reduce: p.reduceSet,
				},
				{
					Event:  group.GroupRemovedEventType,
					Reduce: p.reduceGroupRemoved,
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
					Reduce: reduceInstanceRemovedHelper(GroupManagerRolesInstanceID),
				},
			},
		},
	}
}

func (p *groupManagerRolesProjection) reduceSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupManagerRolesSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-wL3kVd", "reduce.wrong.event.type %s", group.GroupManagerRolesSetEventType)
	}

	if len(e.Roles) == 0 {
		return handler.NewDeleteStatement(
			e,
			[]handler.Condition{
				handler.NewCond(GroupManagerRolesGroupID, e.Aggregate().ID),
				handler.NewCond(GroupManagerRolesInstanceID, e.Aggregate().InstanceID),
			},
		), nil
	}

	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupManagerRolesInstanceID, nil),
			handler.NewCol(GroupManagerRolesGroupID, nil),
		},
		[]handler.Column{
			handler.NewCol(GroupManagerRolesGroupID, e.Aggregate().ID),
			handler.NewCol(GroupManagerRolesResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupManagerRolesInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupManagerRolesRoles, database.TextArray[string](e.Roles)),
			handler.NewCol(GroupManagerRolesCreationDate, handler.OnlySetValueOnInsert(GroupManagerRolesProjectionTable, e.CreatedAt())),
			handler.NewCol(GroupManagerRolesChangeDate, e.CreatedAt()),
			handler.NewCol(GroupManagerRolesSequence, e.Sequence()),
		},
	), nil
}

func (p *groupManagerRolesProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, ok := event.(*group.GroupRemovedEvent); !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-nB6wQs", "reduce.wrong.event.type %s", group.GroupRemovedEventType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupManagerRolesGroupID, event.Aggregate().ID),
			handler.NewCond(GroupManagerRolesInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupManagerRolesProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-c8XwLd", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupManagerRolesInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupManagerRolesResourceOwner, e.Aggregate().ID),
		},
	), nil
}
