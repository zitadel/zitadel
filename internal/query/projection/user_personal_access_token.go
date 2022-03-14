package projection

import (
	"context"

	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	PersonalAccessTokenProjectionTable = "projections.personal_access_tokens"

	PersonalAccessTokenColumnID            = "id"
	PersonalAccessTokenColumnCreationDate  = "creation_date"
	PersonalAccessTokenColumnChangeDate    = "change_date"
	PersonalAccessTokenColumnSequence      = "sequence"
	PersonalAccessTokenColumnResourceOwner = "resource_owner"
	PersonalAccessTokenColumnUserID        = "user_id"
	PersonalAccessTokenColumnExpiration    = "expiration"
	PersonalAccessTokenColumnScopes        = "scopes"
)

type PersonalAccessTokenProjection struct {
	crdb.StatementHandler
}

func NewPersonalAccessTokenProjection(ctx context.Context, config crdb.StatementHandlerConfig) *PersonalAccessTokenProjection {
	p := new(PersonalAccessTokenProjection)
	config.ProjectionName = PersonalAccessTokenProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(PersonalAccessTokenColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(PersonalAccessTokenColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnUserID, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnExpiration, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnScopes, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(PersonalAccessTokenColumnID),
			crdb.NewIndex("user_idx", []string{PersonalAccessTokenColumnUserID}),
			crdb.NewIndex("ro_idx", []string{PersonalAccessTokenColumnResourceOwner}),
		),
	)

	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *PersonalAccessTokenProjection) reducers() []handler.AggregateReducer {
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

func (p *PersonalAccessTokenProjection) reducePersonalAccessTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DVgf7", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
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

func (p *PersonalAccessTokenProjection) reducePersonalAccessTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g7u3F", "reduce.wrong.event.type %s", user.PersonalAccessTokenRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnID, e.TokenID),
		},
	), nil
}

func (p *PersonalAccessTokenProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff3h", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
		},
	), nil
}
