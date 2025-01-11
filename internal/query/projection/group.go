package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupProjectionTable = "projections.group"

	GroupColumnID            = "id"
	GroupColumnCreationDate  = "creation_date"
	GroupColumnChangeDate    = "change_date"
	GroupColumnSequence      = "sequence"
	GroupColumnState         = "state"
	GroupColumnResourceOwner = "resource_owner"
	GroupColumnInstanceID    = "instance_id"
	GroupColumnName          = "name"
	GroupColumnDescription   = "description"
)

type groupProjection struct{}

func newGroupProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupProjection))
}

func (*groupProjection) Name() string {
	return GroupProjectionTable
}

func (*groupProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupColumnID, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(GroupColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(GroupColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnName, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnDescription, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(GroupColumnInstanceID, GroupColumnID),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{GroupColumnResourceOwner})),
		),
	)
}

func (g *groupProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupAddedType,
					Reduce: g.reduceGroupAdded,
				},
				{
					Event:  group.GroupChangedType,
					Reduce: g.reduceGroupChanged,
				},
				{
					Event:  group.GroupDeactivatedType,
					Reduce: g.reduceGroupDeactivated,
				},
				{
					Event:  group.GroupReactivatedType,
					Reduce: g.reduceGroupReactivated,
				},
				{
					Event:  group.GroupRemovedType,
					Reduce: g.reduceGroupRemoved,
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
					Reduce: reduceInstanceRemovedHelper(GroupColumnInstanceID),
				},
			},
		},
	}
}

func (g *groupProjection) reduceGroupAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-l000G", "reduce.wrong.event.type %s", group.GroupAddedType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupColumnID, e.Aggregate().ID),
			handler.NewCol(GroupColumnCreationDate, e.CreationDate()),
			handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
			handler.NewCol(GroupColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupColumnSequence, e.Sequence()),
			handler.NewCol(GroupColumnName, e.Name),
			handler.NewCol(GroupColumnDescription, e.Description),
			handler.NewCol(GroupColumnState, domain.GroupStateActive),
		},
	), nil
}

func (g *groupProjection) reduceGroupChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupChangeEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-u00Eu", "reduce.wrong.event.type %s", group.GroupChangedType)
	}
	if e.Name == nil {
		return handler.NewNoOpStatement(e), nil
	}

	columns := make([]handler.Column, 0, 7)
	columns = append(columns, handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
		handler.NewCol(GroupColumnSequence, e.Sequence()))
	if e.Name != nil {
		columns = append(columns, handler.NewCol(GroupColumnName, *e.Name))
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(GroupColumnDescription, *e.Description))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceGroupDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KKp0f", "reduce.wrong.event.type %s", group.GroupDeactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
			handler.NewCol(GroupColumnSequence, e.Sequence()),
			handler.NewCol(GroupColumnState, domain.GroupStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceGroupReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-8I78f", "reduce.wrong.event.type %s", group.GroupReactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
			handler.NewCol(GroupColumnSequence, e.Sequence()),
			handler.NewCol(GroupColumnState, domain.GroupStateActive),
		},
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-6M5fs", "reduce.wrong.event.type %s", group.GroupRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "GROUP-tchsv", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
