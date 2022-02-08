package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserMetadataProjection struct {
	crdb.StatementHandler
}

const UserMetadataProjectionTable = "zitadel.projections.user_metadata"

func NewUserMetadataProjection(ctx context.Context, config crdb.StatementHandlerConfig) *UserMetadataProjection {
	p := &UserMetadataProjection{}
	config.ProjectionName = UserMetadataProjectionTable
	config.Reducers = p.reducers()
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

const (
	UserMetadataColumnUserID        = "user_id"
	UserMetadataColumnResourceOwner = "resource_owner"
	UserMetadataColumnCreationDate  = "creation_date"
	UserMetadataColumnChangeDate    = "change_date"
	UserMetadataColumnSequence      = "sequence"
	UserMetadataColumnKey           = "key"
	UserMetadataColumnValue         = "value"
)

func (p *UserMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		logging.LogWithFields("HANDL-Sgn5w", "seq", event.Sequence(), "expectedType", user.MetadataSetType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Ghn52", "reduce.wrong.event.type")
	}
	return crdb.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserMetadataColumnUserID, e.Aggregate().ID),
			handler.NewCol(UserMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
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
		logging.LogWithFields("HANDL-Dbfg2", "seq", event.Sequence(), "expectedType", user.MetadataRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bm542", "reduce.wrong.event.type")
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
		logging.LogWithFields("HANDL-Dfbh2", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{user.MetadataRemovedAllType, user.UserRemovedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Bmnf2", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(UserMetadataColumnUserID, event.Aggregate().ID),
		},
	), nil
}
