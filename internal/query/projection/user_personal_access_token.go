package projection

import (
	"context"

	"github.com/lib/pq"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type personalAccessTokenProjection struct {
	crdb.StatementHandler
}

const (
	PersonalAccessTokenProjectionTable = "zitadel.projections.personal_access_tokens"
)

func newPersonalAccessTokenProjection(ctx context.Context, config crdb.StatementHandlerConfig) *personalAccessTokenProjection {
	p := &personalAccessTokenProjection{}
	config.ProjectionName = PersonalAccessTokenProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *personalAccessTokenProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.PersonalAccessTokenAddedType,
					Reduce: p.reducePersonalAccessTokenAdded,
				},
				{
					Event:  user.PersonalAccessTokenRemovedType,
					Reduce: p.reducePersonalAccessTokenRemoved,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
	}
}

const (
	PersonalAccessTokenColumnID            = "id"
	PersonalAccessTokenColumnCreationDate  = "creation_date"
	PersonalAccessTokenColumnChangeDate    = "change_date"
	PersonalAccessTokenColumnResourceOwner = "resource_owner"
	PersonalAccessTokenColumnSequence      = "sequence"
	PersonalAccessTokenColumnUserID        = "user_id"
	PersonalAccessTokenColumnExpiration    = "expiration"
	PersonalAccessTokenColumnScopes        = "scopes"
)

func (p *personalAccessTokenProjection) reducePersonalAccessTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Dbfg2", "seq", event.Sequence(), "expectedType", user.PersonalAccessTokenAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-DVgf7", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(PersonalAccessTokenColumnID, e.TokenID),
			handler.NewCol(PersonalAccessTokenColumnCreationDate, e.CreationDate()),
			handler.NewCol(PersonalAccessTokenColumnChangeDate, e.CreationDate()),
			handler.NewCol(PersonalAccessTokenColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(PersonalAccessTokenColumnSequence, e.Sequence()),
			handler.NewCol(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
			handler.NewCol(PersonalAccessTokenColumnExpiration, e.Expiration),
			handler.NewCol(PersonalAccessTokenColumnScopes, pq.StringArray(e.Scopes)),
		},
	), nil
}

func (p *personalAccessTokenProjection) reducePersonalAccessTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Edf32", "seq", event.Sequence(), "expectedType", user.PersonalAccessTokenRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g7u3F", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnID, e.TokenID),
		},
	), nil
}

func (p *personalAccessTokenProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GEg43", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dff3h", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
		},
	), nil
}
