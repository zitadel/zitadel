package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	PersonalAccessTokenProjectionTable = "projections.personal_access_tokens3"

	PersonalAccessTokenColumnID            = "id"
	PersonalAccessTokenColumnCreationDate  = "creation_date"
	PersonalAccessTokenColumnChangeDate    = "change_date"
	PersonalAccessTokenColumnSequence      = "sequence"
	PersonalAccessTokenColumnResourceOwner = "resource_owner"
	PersonalAccessTokenColumnInstanceID    = "instance_id"
	PersonalAccessTokenColumnUserID        = "user_id"
	PersonalAccessTokenColumnExpiration    = "expiration"
	PersonalAccessTokenColumnScopes        = "scopes"
	PersonalAccessTokenColumnOwnerRemoved  = "owner_removed"
)

type personalAccessTokenProjection struct{}

func newPersonalAccessTokenProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(personalAccessTokenProjection))
}

func (*personalAccessTokenProjection) Name() string {
	return PersonalAccessTokenProjectionTable
}

func (*personalAccessTokenProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(PersonalAccessTokenColumnID, handler.ColumnTypeText),
			handler.NewColumn(PersonalAccessTokenColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(PersonalAccessTokenColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(PersonalAccessTokenColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(PersonalAccessTokenColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(PersonalAccessTokenColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(PersonalAccessTokenColumnUserID, handler.ColumnTypeText),
			handler.NewColumn(PersonalAccessTokenColumnExpiration, handler.ColumnTypeTimestamp),
			handler.NewColumn(PersonalAccessTokenColumnScopes, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(PersonalAccessTokenColumnOwnerRemoved, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(PersonalAccessTokenColumnInstanceID, PersonalAccessTokenColumnID),
			handler.WithIndex(handler.NewIndex("user_id", []string{PersonalAccessTokenColumnUserID})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{PersonalAccessTokenColumnResourceOwner})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{PersonalAccessTokenColumnOwnerRemoved})),
		),
	)
}

func (p *personalAccessTokenProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Reduce: reduceInstanceRemovedHelper(PersonalAccessTokenColumnInstanceID),
				},
			},
		},
	}
}

func (p *personalAccessTokenProjection) reducePersonalAccessTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-DVgf7", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(PersonalAccessTokenColumnID, e.TokenID),
			handler.NewCol(PersonalAccessTokenColumnCreationDate, e.CreationDate()),
			handler.NewCol(PersonalAccessTokenColumnChangeDate, e.CreationDate()),
			handler.NewCol(PersonalAccessTokenColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(PersonalAccessTokenColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(PersonalAccessTokenColumnSequence, e.Sequence()),
			handler.NewCol(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
			handler.NewCol(PersonalAccessTokenColumnExpiration, e.Expiration),
			handler.NewCol(PersonalAccessTokenColumnScopes, database.TextArray[string](e.Scopes)),
		},
	), nil
}

func (p *personalAccessTokenProjection) reducePersonalAccessTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g7u3F", "reduce.wrong.event.type %s", user.PersonalAccessTokenRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnID, e.TokenID),
			handler.NewCond(PersonalAccessTokenColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *personalAccessTokenProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dff3h", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
			handler.NewCond(PersonalAccessTokenColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *personalAccessTokenProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-zQVhl", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(PersonalAccessTokenColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(PersonalAccessTokenColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
