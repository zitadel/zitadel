package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserMetadataProjectionTable = "projections.user_metadata4"

	UserMetadataColumnUserID        = "user_id"
	UserMetadataColumnCreationDate  = "creation_date"
	UserMetadataColumnChangeDate    = "change_date"
	UserMetadataColumnSequence      = "sequence"
	UserMetadataColumnResourceOwner = "resource_owner"
	UserMetadataColumnInstanceID    = "instance_id"
	UserMetadataColumnKey           = "key"
	UserMetadataColumnValue         = "value"
	UserMetadataColumnOwnerRemoved  = "owner_removed"
)

type userMetadataProjection struct {
	crdb.StatementHandler
}

func newUserMetadataProjection(ctx context.Context, config crdb.StatementHandlerConfig) *userMetadataProjection {
	p := new(userMetadataProjection)
	config.ProjectionName = UserMetadataProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserMetadataColumnUserID, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserMetadataColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserMetadataColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(UserMetadataColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnKey, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnValue, crdb.ColumnTypeBytes, crdb.Nullable()),
			crdb.NewColumn(UserMetadataColumnOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(UserMetadataColumnInstanceID, UserMetadataColumnUserID, UserMetadataColumnKey),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{UserGrantResourceOwner})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{UserMetadataColumnOwnerRemoved})),
		),
	)

	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *userMetadataProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.MetadataSetType,
					Reduce: p.reduceMetadataSet,
				},
				{
					Event:  user.MetadataRemovedType,
					Reduce: p.reduceMetadataRemoved,
				},
				{
					Event:  user.MetadataRemovedAllType,
					Reduce: p.reduceMetadataRemovedAll,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceMetadataRemovedAll,
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
					Reduce: reduceInstanceRemovedHelper(UserMetadataColumnInstanceID),
				},
			},
		},
	}
}

func (p *userMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Ghn52", "reduce.wrong.event.type %s", user.MetadataSetType)
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserMetadataColumnInstanceID, nil),
			handler.NewCol(UserMetadataColumnUserID, nil),
			handler.NewCol(UserMetadataColumnKey, e.Key),
		},
		[]handler.Column{
			handler.NewCol(UserMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(UserMetadataColumnUserID, e.Aggregate().ID),
			handler.NewCol(UserMetadataColumnKey, e.Key),
			handler.NewCol(UserMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserMetadataColumnCreationDate, e.CreationDate()),
			handler.NewCol(UserMetadataColumnChangeDate, e.CreationDate()),
			handler.NewCol(UserMetadataColumnSequence, e.Sequence()),
			handler.NewCol(UserMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *userMetadataProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bm542", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnUserID, e.Aggregate().ID),
			handler.NewCond(UserMetadataColumnKey, e.Key),
			handler.NewCond(UserAuthMethodInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userMetadataProjection) reduceMetadataRemovedAll(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *user.MetadataRemovedAllEvent,
		*user.UserRemovedEvent:
		//ok
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bmnf2", "reduce.wrong.event.type %v", []eventstore.EventType{user.MetadataRemovedAllType, user.UserRemovedType})
	}
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnUserID, event.Aggregate().ID),
			handler.NewCond(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userMetadataProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-oqwul", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(UserMetadataColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
