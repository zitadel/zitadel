package projection

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	v3 "github.com/zitadel/zitadel/internal/eventstore/handler/v3"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	PersonalAccessTokenProjectionTable = "projections.personal_access_tokens2"

	PersonalAccessTokenColumnID            = "id"
	PersonalAccessTokenColumnCreationDate  = "creation_date"
	PersonalAccessTokenColumnChangeDate    = "change_date"
	PersonalAccessTokenColumnResourceOwner = "resource_owner"
	PersonalAccessTokenColumnInstanceID    = "instance_id"
	PersonalAccessTokenColumnUserID        = "user_id"
	PersonalAccessTokenColumnExpiration    = "expiration"
	PersonalAccessTokenColumnScopes        = "scopes"
)

type personalAccessTokenProjection struct {
	crdb.StatementHandler
}

func newPersonalAccessTokenProjection(ctx context.Context, config v3.Config) *v3.IDProjection {
	p := new(personalAccessTokenProjection)
	config.Check = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(PersonalAccessTokenColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnUserID, crdb.ColumnTypeText),
			crdb.NewColumn(PersonalAccessTokenColumnExpiration, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PersonalAccessTokenColumnScopes, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(PersonalAccessTokenColumnInstanceID, PersonalAccessTokenColumnID),
			crdb.WithIndex(crdb.NewIndex("pat_user_idx", []string{PersonalAccessTokenColumnUserID})),
			crdb.WithIndex(crdb.NewIndex("pat_ro_idx", []string{PersonalAccessTokenColumnResourceOwner})),
		),
	)

	config.Reduces = map[eventstore.AggregateType][]v3.Reducer{
		user.AggregateType: {
			{
				Event:  user.PersonalAccessTokenAddedType,
				Reduce: p.reducePersonalAccessTokenAdded,
			},
			{
				Event:          user.PersonalAccessTokenRemovedType,
				Reduce:         p.reducePersonalAccessTokenRemoved,
				PreviousEvents: p.previousEventsRemoved,
			},
			{
				Event:          user.UserRemovedType,
				Reduce:         p.reduceUserRemoved,
				PreviousEvents: p.previousEventsUserRemoved,
			},
		},
		instance.AggregateType: {
			{
				Event:  instance.InstanceRemovedEventType,
				Reduce: reduceInstanceRemovedHelper(PersonalAccessTokenColumnInstanceID),
			},
		},
	}

	return v3.New(PersonalAccessTokenProjectionTable, config)
}

func (p *personalAccessTokenProjection) reducePersonalAccessTokenAdded(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCol(PersonalAccessTokenColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(PersonalAccessTokenColumnUserID, e.Aggregate().ID),
			handler.NewCol(PersonalAccessTokenColumnExpiration, e.Expiration),
			handler.NewCol(PersonalAccessTokenColumnScopes, database.StringArray(e.Scopes)),
		},
	), nil
}

func (p *personalAccessTokenProjection) reducePersonalAccessTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *personalAccessTokenProjection) previousEventsRemoved(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*user.PersonalAccessTokenRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Ca7Yv", "reduce.wrong.event.type %s", user.PersonalAccessTokenRemovedType)
	}

	row := tx.QueryRow("SELECT "+PersonalAccessTokenColumnChangeDate+" FROM "+PersonalAccessTokenProjectionTable+" WHERE "+PersonalAccessTokenColumnUserID+" = $1 AND "+PersonalAccessTokenColumnInstanceID+" = $2 AND "+PersonalAccessTokenColumnID+" = $3 FOR UPDATE", e.Aggregate().ID, e.Aggregate().InstanceID, e.Id)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		// SetTx(tx).
		InstanceID(e.Aggregate().InstanceID).
		SystemTime(e.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(e.Aggregate().ID).
		EventTypes(
			user.PersonalAccessTokenAddedType,
			user.PersonalAccessTokenRemovedType,
			user.UserRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
}

func (p *personalAccessTokenProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *personalAccessTokenProjection) previousEventsUserRemoved(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Ca7Yv", "reduce.wrong.event.type %s", user.PersonalAccessTokenRemovedType)
	}

	_, err := tx.Exec("SELECT 1 FROM "+PersonalAccessTokenProjectionTable+" WHERE "+PersonalAccessTokenColumnUserID+" = $1 AND "+PersonalAccessTokenColumnInstanceID+" = $2 FOR UPDATE", e.Aggregate().ID, e.Aggregate().InstanceID)
	return nil, err
}
