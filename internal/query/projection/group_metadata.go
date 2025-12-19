package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupMetadataProjectionTable = "projections.group_metadata"

	GroupMetadataColumnGroupID       = "group_id"
	GroupMetadataColumnCreationDate  = "creation_date"
	GroupMetadataColumnChangeDate    = "change_date"
	GroupMetadataColumnSequence      = "sequence"
	GroupMetadataColumnResourceOwner = "resource_owner"
	GroupMetadataColumnInstanceID    = "instance_id"
	GroupMetadataColumnKey           = "key"
	GroupMetadataColumnValue         = "value"
)

type groupMetadataProjection struct{}

func newGroupMetadataProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupMetadataProjection))
}

func (*groupMetadataProjection) Name() string {
	return GroupMetadataProjectionTable
}

func (*groupMetadataProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupMetadataColumnGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupMetadataColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupMetadataColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(GroupMetadataColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(GroupMetadataColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupMetadataColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupMetadataColumnKey, handler.ColumnTypeText),
			handler.NewColumn(GroupMetadataColumnValue, handler.ColumnTypeBytes, handler.Nullable()),
		},
			handler.NewPrimaryKey(GroupMetadataColumnInstanceID, GroupMetadataColumnGroupID, GroupMetadataColumnKey),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{GroupGrantResourceOwner})),
		),
	)
}

func (p *groupMetadataProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.MetadataSetType,
					Reduce: p.reduceMetadataSet,
				},
				{
					Event:  group.MetadataRemovedType,
					Reduce: p.reduceMetadataRemoved,
				},
				{
					Event:  group.MetadataRemovedAllType,
					Reduce: p.reduceMetadataRemovedAll,
				},
				{
					Event:  group.GroupRemovedType,
					Reduce: p.reduceMetadataRemovedAll,
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
					Reduce: reduceInstanceRemovedHelper(GroupMetadataColumnInstanceID),
				},
			},
		},
	}
}

func (p *groupMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KOn12", "reduce.wrong.event.type %s", group.MetadataSetType)
	}
	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupMetadataColumnGroupID, e.Aggregate().ID),
			handler.NewCol(GroupMetadataColumnKey, e.Key),
		},
		[]handler.Column{
			handler.NewCol(GroupMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupMetadataColumnGroupID, e.Aggregate().ID),
			handler.NewCol(GroupMetadataColumnKey, e.Key),
			handler.NewCol(GroupMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupMetadataColumnCreationDate, handler.OnlySetValueOnInsert(GroupMetadataProjectionTable, e.CreationDate())),
			handler.NewCol(GroupMetadataColumnChangeDate, e.CreationDate()),
			handler.NewCol(GroupMetadataColumnSequence, e.Sequence()),
			handler.NewCol(GroupMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *groupMetadataProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Cp522", "reduce.wrong.event.type %s", group.MetadataRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupMetadataColumnGroupID, e.Aggregate().ID),
			handler.NewCond(GroupMetadataColumnKey, e.Key),
			handler.NewCond(GroupMetadataColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupMetadataProjection) reduceMetadataRemovedAll(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *group.MetadataRemovedAllEvent,
		*group.GroupRemovedEvent:
		//ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Cmnf2", "reduce.wrong.event.type %v", []eventstore.EventType{group.MetadataRemovedAllType, group.GroupRemovedType})
	}
	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(GroupMetadataColumnGroupID, event.Aggregate().ID),
			handler.NewCond(GroupMetadataColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *groupMetadataProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pswul", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupMetadataColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
