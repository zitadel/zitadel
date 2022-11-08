package projection

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	v3 "github.com/zitadel/zitadel/internal/eventstore/handler/v3"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserMetadataProjectionTable = "projections.user_metadata3"

	UserMetadataColumnUserID        = "user_id"
	UserMetadataColumnCreationDate  = "creation_date"
	UserMetadataColumnChangeDate    = "change_date"
	UserMetadataColumnResourceOwner = "resource_owner"
	UserMetadataColumnInstanceID    = "instance_id"
	UserMetadataColumnKey           = "key"
	UserMetadataColumnValue         = "value"
)

type userMetadataProjection struct {
	crdb.StatementHandler
}

func newUserMetadataProjection(ctx context.Context, config v3.Config) *v3.IDProjection {
	p := new(userMetadataProjection)

	config.Check = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserMetadataColumnUserID, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserMetadataColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserMetadataColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnKey, crdb.ColumnTypeText),
			crdb.NewColumn(UserMetadataColumnValue, crdb.ColumnTypeBytes, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(UserMetadataColumnInstanceID, UserMetadataColumnUserID, UserMetadataColumnKey),
			crdb.WithIndex(crdb.NewIndex("usr_md_ro_idx", []string{UserGrantResourceOwner})),
		),
	)

	config.Reduces = map[eventstore.AggregateType][]v3.Reducer{
		user.AggregateType: {
			{
				Event:          user.MetadataSetType,
				Reduce:         p.reduceMetadataSet,
				PreviousEvents: p.previousEventsSet,
			},
			{
				Event:          user.MetadataRemovedType,
				Reduce:         p.reduceMetadataRemoved,
				PreviousEvents: p.previousEventsRemoved,
			},
			{
				Event:          user.MetadataRemovedAllType,
				Reduce:         p.reduceMetadataRemovedAll,
				PreviousEvents: p.previousEventsRemovedAll,
			},
			{
				Event:          user.UserRemovedType,
				Reduce:         p.reduceMetadataRemovedAll,
				PreviousEvents: p.previousEventsRemovedAll,
			},
		},
	}

	return v3.StartSubscriptionIDProjection(ctx, UserMetadataProjectionTable, config)
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
			handler.NewCol(UserMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *userMetadataProjection) previousEventsSet(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-u7mzn", "reduce.wrong.event.type %s", user.MetadataSetType)
	}

	row := tx.QueryRow("SELECT "+UserMetadataColumnChangeDate+" FROM "+UserMetadataProjectionTable+" WHERE "+UserMetadataColumnUserID+" = $1 AND "+UserMetadataColumnInstanceID+" = $2 AND "+UserMetadataColumnKey+" = $3 FOR UPDATE", e.Aggregate().ID, e.Aggregate().InstanceID, e.Key)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(e.Aggregate().InstanceID).
		SystemTime(e.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(e.Aggregate().ID).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType,
			user.MetadataRemovedAllType,
			user.UserRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
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
			handler.NewCond(UserMetadataColumnInstanceID, event.Aggregate().InstanceID),
			handler.NewCond(UserMetadataColumnKey, e.Key),
		},
	), nil
}

func (p *userMetadataProjection) previousEventsRemoved(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-TAFtA", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}

	row := tx.QueryRow("SELECT "+UserMetadataColumnChangeDate+" FROM "+UserMetadataProjectionTable+" WHERE "+UserMetadataColumnUserID+" = $1 AND "+UserMetadataColumnInstanceID+" = $2 AND "+UserMetadataColumnKey+" = $3 FOR UPDATE", e.Aggregate().ID, e.Aggregate().InstanceID, e.Key)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(e.Aggregate().InstanceID).
		SystemTime(e.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(e.Aggregate().ID).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType,
			user.MetadataRemovedAllType,
			user.UserRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
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
			handler.NewCond(UserMetadataColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userMetadataProjection) previousEventsRemovedAll(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	_, err := tx.Exec("SELECT 1 FROM "+UserMetadataProjectionTable+" WHERE "+UserMetadataColumnUserID+" = $1 AND "+UserMetadataColumnInstanceID+" = $2 FOR UPDATE", event.Aggregate().ID, event.Aggregate().InstanceID)
	return nil, err
}
