package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	UserMetadataProjectionTable = "projections.user_metadata"

	UserMetadataColumnUserID        = "user_id"
	UserMetadataColumnCreationDate  = "creation_date"
	UserMetadataColumnChangeDate    = "change_date"
	UserMetadataColumnSequence      = "sequence"
	UserMetadataColumnResourceOwner = "resource_owner"
	UserMetadataColumnInstanceID    = "instance_id"
	UserMetadataColumnKey           = "key"
	UserMetadataColumnValue         = "value"
)

type UserMetadataProjection struct {
	crdb.StatementHandler
}

func NewUserMetadataProjection(ctx context.Context, config crdb.StatementHandlerConfig) *UserMetadataProjection {
	p := new(UserMetadataProjection)
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
		},
			crdb.NewPrimaryKey(UserMetadataColumnUserID),
			crdb.NewIndex("ro_idx", []string{UserGrantResourceOwner}),
		),
	)

	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *UserMetadataProjection) reducers() []handler.AggregateReducer {
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
	}
}

func (p *UserMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Ghn52", "reduce.wrong.event.type %s", user.MetadataSetType)
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserMetadataColumnUserID, e.Aggregate().ID),
			handler.NewCol(UserMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(UserMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(UserMetadataColumnCreationDate, e.CreationDate()),
			handler.NewCol(UserMetadataColumnChangeDate, e.CreationDate()),
			handler.NewCol(UserMetadataColumnSequence, e.Sequence()),
			handler.NewCol(UserMetadataColumnKey, e.Key),
			handler.NewCol(UserMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *UserMetadataProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Bm542", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnUserID, e.Aggregate().ID),
			handler.NewCond(UserMetadataColumnKey, e.Key),
		},
	), nil
}

func (p *UserMetadataProjection) reduceMetadataRemovedAll(event eventstore.Event) (*handler.Statement, error) {
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
		},
	), nil
}
