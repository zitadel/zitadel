package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserMetadataProjectionTable = "projections.user_metadata5"

	UserMetadataColumnUserID        = "user_id"
	UserMetadataColumnCreationDate  = "creation_date"
	UserMetadataColumnChangeDate    = "change_date"
	UserMetadataColumnSequence      = "sequence"
	UserMetadataColumnResourceOwner = "resource_owner"
	UserMetadataColumnInstanceID    = "instance_id"
	UserMetadataColumnKey           = "key"
	UserMetadataColumnValue         = "value"
)

type userMetadataProjection struct{}

func newUserMetadataProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(userMetadataProjection))
}

func (*userMetadataProjection) Name() string {
	return UserMetadataProjectionTable
}

func (*userMetadataProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(UserMetadataColumnUserID, handler.ColumnTypeText),
			handler.NewColumn(UserMetadataColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserMetadataColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserMetadataColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(UserMetadataColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(UserMetadataColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(UserMetadataColumnKey, handler.ColumnTypeText),
			handler.NewColumn(UserMetadataColumnValue, handler.ColumnTypeBytes, handler.Nullable()),
		},
			handler.NewPrimaryKey(UserMetadataColumnInstanceID, UserMetadataColumnUserID, UserMetadataColumnKey),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{UserGrantResourceOwner})),
		),
	)
}

func (p *userMetadataProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Reduce: reduceInstanceRemovedHelper(UserMetadataColumnInstanceID),
				},
			},
		},
	}
}

func (p *userMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ghn52", "reduce.wrong.event.type %s", user.MetadataSetType)
	}
	return handler.NewUpsertStatement(
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
			handler.NewCol(UserMetadataColumnCreationDate, handler.OnlySetValueOnInsert(UserMetadataProjectionTable, e.CreationDate())),
			handler.NewCol(UserMetadataColumnChangeDate, e.CreationDate()),
			handler.NewCol(UserMetadataColumnSequence, e.Sequence()),
			handler.NewCol(UserMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *userMetadataProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bm542", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}
	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bmnf2", "reduce.wrong.event.type %v", []eventstore.EventType{user.MetadataRemovedAllType, user.UserRemovedType})
	}
	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-oqwul", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(UserMetadataColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
