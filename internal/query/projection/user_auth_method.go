package projection

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	v3 "github.com/zitadel/zitadel/internal/eventstore/handler/v3"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserAuthMethodTable = "projections.user_auth_methods3"

	UserAuthMethodUserIDCol        = "user_id"
	UserAuthMethodTypeCol          = "method_type"
	UserAuthMethodTokenIDCol       = "token_id"
	UserAuthMethodCreationDateCol  = "creation_date"
	UserAuthMethodChangeDateCol    = "change_date"
	UserAuthMethodResourceOwnerCol = "resource_owner"
	UserAuthMethodInstanceIDCol    = "instance_id"
	UserAuthMethodStateCol         = "state"
	UserAuthMethodNameCol          = "name"
)

type userAuthMethodProjection struct {
	crdb.StatementHandler
}

func newUserAuthMethodProjection(ctx context.Context, config v3.Config) *v3.IDProjection {
	p := new(userAuthMethodProjection)
	config.Check = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserAuthMethodUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserAuthMethodTokenIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserAuthMethodChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserAuthMethodStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserAuthMethodResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodNameCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(UserAuthMethodInstanceIDCol, UserAuthMethodUserIDCol, UserAuthMethodTypeCol, UserAuthMethodTokenIDCol),
			crdb.WithIndex(crdb.NewIndex("auth_meth_ro_idx", []string{UserAuthMethodResourceOwnerCol})),
		),
	)

	config.Reduces = map[eventstore.AggregateType][]v3.Reducer{
		user.AggregateType: {
			{
				Event:  user.HumanPasswordlessTokenAddedType,
				Reduce: p.reduceInitAuthMethod,
			},
			{
				Event:  user.HumanU2FTokenAddedType,
				Reduce: p.reduceInitAuthMethod,
			},
			{
				Event:  user.HumanMFAOTPAddedType,
				Reduce: p.reduceInitAuthMethod,
			},
			{
				Event:  user.HumanPasswordlessTokenVerifiedType,
				Reduce: p.reduceActivateEvent,
			},
			{
				Event:  user.HumanU2FTokenVerifiedType,
				Reduce: p.reduceActivateEvent,
			},
			{
				Event:  user.HumanMFAOTPVerifiedType,
				Reduce: p.reduceActivateEvent,
			},
			{
				Event:  user.HumanPasswordlessTokenRemovedType,
				Reduce: p.reduceRemoveAuthMethod,
			},
			{
				Event:  user.HumanU2FTokenRemovedType,
				Reduce: p.reduceRemoveAuthMethod,
			},
			{
				Event:  user.HumanMFAOTPRemovedType,
				Reduce: p.reduceRemoveAuthMethod,
			},
		},
		instance.AggregateType: {
			{
				Event:  instance.InstanceRemovedEventType,
				Reduce: reduceInstanceRemovedHelper(UserAuthMethodInstanceIDCol),
			},
		},
	}

	return v3.New(UserAuthMethodTable, config)
}

func (p *userAuthMethodProjection) reduceInitAuthMethod(event eventstore.Event) (*handler.Statement, error) {
	tokenID := ""
	var methodType domain.UserAuthMethodType
	switch e := event.(type) {
	case *user.HumanPasswordlessAddedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
	case *user.HumanU2FAddedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
	case *user.HumanOTPAddedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	return crdb.NewUpsertStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodInstanceIDCol, nil),
			handler.NewCol(UserAuthMethodUserIDCol, nil),
			handler.NewCol(UserAuthMethodTypeCol, nil),
			handler.NewCol(UserAuthMethodTokenIDCol, nil),
		},
		[]handler.Column{
			handler.NewCol(UserAuthMethodTokenIDCol, tokenID),
			handler.NewCol(UserAuthMethodCreationDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCol(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCol(UserAuthMethodUserIDCol, event.Aggregate().ID),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateNotReady),
			handler.NewCol(UserAuthMethodTypeCol, methodType),
			handler.NewCol(UserAuthMethodNameCol, ""),
		},
	), nil
}

func (p *userAuthMethodProjection) previousEventsInit(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	tokenID := ""
	var methodType domain.UserAuthMethodType
	switch e := event.(type) {
	case *user.HumanPasswordlessAddedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
	case *user.HumanU2FAddedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
	case *user.HumanOTPAddedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	row := tx.QueryRow("SELECT "+UserAuthMethodChangeDateCol+
		" FROM "+UserAuthMethodTable+" WHERE "+
		UserAuthMethodUserIDCol+" = $1 AND "+
		UserAuthMethodInstanceIDCol+" = $2 AND "+
		UserAuthMethodTypeCol+" = $3 AND "+
		UserAuthMethodTokenIDCol+" = $4 FOR UPDATE",
		event.Aggregate().ID,
		event.Aggregate().InstanceID,
		methodType,
		tokenID)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if changeDate.IsZero() {
		return nil, nil
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(event.Aggregate().InstanceID).
		SystemTime(event.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(event.Aggregate().ID).
		EventTypes(
			user.HumanPasswordlessTokenAddedType,
			user.HumanU2FTokenAddedType,
			user.HumanMFAOTPAddedType,
			user.HumanPasswordlessTokenVerifiedType,
			user.HumanU2FTokenVerifiedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanPasswordlessTokenRemovedType,
			user.HumanU2FTokenRemovedType,
			user.HumanMFAOTPRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
}

func (p *userAuthMethodProjection) reduceActivateEvent(event eventstore.Event) (*handler.Statement, error) {
	tokenID := ""
	name := ""
	var methodType domain.UserAuthMethodType

	switch e := event.(type) {
	case *user.HumanPasswordlessVerifiedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
		name = e.WebAuthNTokenName
	case *user.HumanU2FVerifiedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
		name = e.WebAuthNTokenName
	case *user.HumanOTPVerifiedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodNameCol, name),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateReady),
		},
		[]handler.Condition{
			handler.NewCond(UserAuthMethodUserIDCol, event.Aggregate().ID),
			handler.NewCond(UserAuthMethodTypeCol, methodType),
			handler.NewCond(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(UserAuthMethodTokenIDCol, tokenID),
		},
	), nil
}

func (p *userAuthMethodProjection) previousEventsActivate(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	tokenID := ""
	var methodType domain.UserAuthMethodType

	switch e := event.(type) {
	case *user.HumanPasswordlessVerifiedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
	case *user.HumanU2FVerifiedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
	case *user.HumanOTPVerifiedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-U9gg5", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	row := tx.QueryRow("SELECT "+UserAuthMethodChangeDateCol+
		" FROM "+UserAuthMethodTable+" WHERE "+
		UserAuthMethodUserIDCol+" = $1 AND "+
		UserAuthMethodInstanceIDCol+" = $2 AND "+
		UserAuthMethodTypeCol+" = $3 AND "+
		UserAuthMethodTokenIDCol+" = $4 FOR UPDATE",
		event.Aggregate().ID,
		event.Aggregate().InstanceID,
		methodType,
		tokenID)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(event.Aggregate().InstanceID).
		SystemTime(event.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(event.Aggregate().ID).
		EventTypes(
			user.HumanPasswordlessTokenAddedType,
			user.HumanU2FTokenAddedType,
			user.HumanMFAOTPAddedType,
			user.HumanPasswordlessTokenVerifiedType,
			user.HumanU2FTokenVerifiedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanPasswordlessTokenRemovedType,
			user.HumanU2FTokenRemovedType,
			user.HumanMFAOTPRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
}

func (p *userAuthMethodProjection) reduceRemoveAuthMethod(event eventstore.Event) (*handler.Statement, error) {
	var tokenID string
	var methodType domain.UserAuthMethodType
	switch e := event.(type) {
	case *user.HumanPasswordlessRemovedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
	case *user.HumanU2FRemovedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
	case *user.HumanOTPRemovedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}
	conditions := []handler.Condition{
		handler.NewCond(UserAuthMethodUserIDCol, event.Aggregate().ID),
		handler.NewCond(UserAuthMethodTypeCol, methodType),
		handler.NewCond(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
	}
	if tokenID != "" {
		conditions = append(conditions, handler.NewCond(UserAuthMethodTokenIDCol, tokenID))
	}
	return crdb.NewDeleteStatement(
		event,
		conditions,
	), nil
}

func (p *userAuthMethodProjection) previousEventsRemove(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	var tokenID string
	var methodType domain.UserAuthMethodType
	switch e := event.(type) {
	case *user.HumanPasswordlessRemovedEvent:
		methodType = domain.UserAuthMethodTypePasswordless
		tokenID = e.WebAuthNTokenID
	case *user.HumanU2FRemovedEvent:
		methodType = domain.UserAuthMethodTypeU2F
		tokenID = e.WebAuthNTokenID
	case *user.HumanOTPRemovedEvent:
		methodType = domain.UserAuthMethodTypeOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-7z95W", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	row := tx.QueryRow("SELECT "+UserAuthMethodChangeDateCol+
		" FROM "+UserAuthMethodTable+" WHERE "+
		UserAuthMethodUserIDCol+" = $1 AND "+
		UserAuthMethodInstanceIDCol+" = $2 AND "+
		UserAuthMethodTypeCol+" = $3 AND "+
		UserAuthMethodTokenIDCol+" = $4 FOR UPDATE",
		event.Aggregate().ID,
		event.Aggregate().InstanceID,
		methodType,
		tokenID)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		InstanceID(event.Aggregate().InstanceID).
		SystemTime(event.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(event.Aggregate().ID).
		EventTypes(
			user.HumanPasswordlessTokenAddedType,
			user.HumanU2FTokenAddedType,
			user.HumanMFAOTPAddedType,
			user.HumanPasswordlessTokenVerifiedType,
			user.HumanU2FTokenVerifiedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanPasswordlessTokenRemovedType,
			user.HumanU2FTokenRemovedType,
			user.HumanMFAOTPRemovedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
}
