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
)

const (
	GroupProjectionTable = "projections.groups1"

	GroupColumnID            = "id"
	GroupColumnName          = "name"
	GroupColumnResourceOwner = "resource_owner"
	GroupColumnInstanceID    = "instance_id"
	GroupColumnCreationDate  = "creation_date"
	GroupColumnChangeDate    = "change_date"
	GroupColumnSequence      = "sequence"
	GroupColumnState         = "state"
	GroupColumnDescription   = "description"
)

type groupProjection struct{}

func (g *groupProjection) Name() string {
	return GroupProjectionTable
}

func newGroupProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupProjection))
}

func (g *groupProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupColumnID, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnName, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnDescription, handler.ColumnTypeText),
			handler.NewColumn(GroupColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(GroupColumnState, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(GroupColumnInstanceID, GroupColumnID),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{GroupColumnResourceOwner})),
			handler.WithIndex(handler.NewIndex("state", []string{GroupColumnState})),
		),
	)
}

func (g *groupProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupAddedEventType,
					Reduce: g.reduceGroupAdded,
				},
				{
					Event:  group.GroupChangedEventType,
					Reduce: g.reduceGroupChanged,
				},
				{
					Event:  group.GroupRemovedEventType,
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
	e, err := assertEvent[*group.GroupAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupColumnID, e.Aggregate().ID),
			handler.NewCol(GroupColumnName, e.Name),
			handler.NewCol(GroupColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupColumnDescription, e.Description),
			handler.NewCol(GroupColumnCreationDate, e.CreationDate()),
			handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
			handler.NewCol(GroupColumnSequence, e.Sequence()),
			handler.NewCol(GroupColumnState, domain.GroupStateActive),
		},
	), nil
}

func (g *groupProjection) reduceGroupChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupChangedEvent](event)
	if err != nil {
		return nil, err
	}

	columns := make([]handler.Column, 0, 4)

	if e.Name != nil {
		columns = append(columns, handler.NewCol(GroupColumnName, *e.Name))
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(GroupColumnDescription, *e.Description))
	}
	if len(columns) == 0 {
		return handler.NewNoOpStatement(e), nil
	}

	columns = append(
		columns,
		handler.NewCol(GroupColumnChangeDate, e.CreationDate()),
		handler.NewCol(GroupColumnSequence, e.Sequence()),
	)

	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupColumnID, e.Aggregate().ID),
			handler.NewCond(GroupColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
