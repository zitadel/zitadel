package projection

import (
	"context"
	"crypto/sha256"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	GroupExternalIDsProjectionTable = "projections.group_external_ids2"

	GroupExternalIDsColumnGroupID       = "group_id"
	GroupExternalIDsColumnScope         = "scope"
	GroupExternalIDsColumnResourceOwner = "resource_owner"
	GroupExternalIDsColumnInstanceID    = "instance_id"
	GroupExternalIDsColumnValue         = "external_id"
	GroupExternalIDsColumnDigest        = "external_id_sha256"
)

type groupExternalIDsProjection struct{}

func (g *groupExternalIDsProjection) Name() string {
	return GroupExternalIDsProjectionTable
}

func newGroupExternalIDsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupExternalIDsProjection))
}

func (*groupExternalIDsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupExternalIDsColumnGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupExternalIDsColumnScope, handler.ColumnTypeText),
			handler.NewColumn(GroupExternalIDsColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupExternalIDsColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupExternalIDsColumnValue, handler.ColumnTypeText),
			handler.NewColumn(GroupExternalIDsColumnDigest, handler.ColumnTypeBytes),
		},
			handler.NewPrimaryKey(GroupExternalIDsColumnInstanceID, GroupExternalIDsColumnGroupID, GroupExternalIDsColumnScope),
			handler.WithIndex(handler.NewIndex("lookup", []string{GroupExternalIDsColumnInstanceID, GroupExternalIDsColumnResourceOwner, GroupExternalIDsColumnScope, GroupExternalIDsColumnDigest})),
			handler.WithIndex(handler.NewIndex("group_id", []string{GroupExternalIDsColumnGroupID})),
		),
	)
}

func (g *groupExternalIDsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupExternalIDSetEventType,
					Reduce: g.reduceExternalIDSet,
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

func (g *groupExternalIDsProjection) reduceExternalIDSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupExternalIDSetEvent](event)
	if err != nil {
		return nil, err
	}
	if e.ExternalID == "" {
		return handler.NewDeleteStatement(e, []handler.Condition{
			handler.NewCond(GroupExternalIDsColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupExternalIDsColumnGroupID, e.Aggregate().ID),
			handler.NewCond(GroupExternalIDsColumnScope, e.Scope),
		}), nil
	}
	digest := sha256.Sum256([]byte(e.ExternalID))
	return handler.NewUpsertStatement(e, []handler.Column{
		handler.NewCol(GroupExternalIDsColumnInstanceID, nil), handler.NewCol(GroupExternalIDsColumnGroupID, nil), handler.NewCol(GroupExternalIDsColumnScope, nil),
	}, []handler.Column{
		handler.NewCol(GroupExternalIDsColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(GroupExternalIDsColumnGroupID, e.Aggregate().ID),
		handler.NewCol(GroupExternalIDsColumnResourceOwner, e.Aggregate().ResourceOwner),
		handler.NewCol(GroupExternalIDsColumnScope, e.Scope),
		handler.NewCol(GroupExternalIDsColumnValue, e.ExternalID),
		handler.NewCol(GroupExternalIDsColumnDigest, digest[:]),
	}), nil
}

func (g *groupExternalIDsProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupExternalIDsColumnGroupID, e.Aggregate().ID),
			handler.NewCond(GroupExternalIDsColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCond(GroupExternalIDsColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupExternalIDsProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupExternalIDsColumnResourceOwner, e.Aggregate().ID),
			handler.NewCond(GroupExternalIDsColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
