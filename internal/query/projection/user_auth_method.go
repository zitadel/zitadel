package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserAuthMethodTable = "projections.user_auth_methods"

	UserAuthMethodUserIDCol        = "user_id"
	UserAuthMethodTypeCol          = "method_type"
	UserAuthMethodTokenIDCol       = "token_id"
	UserAuthMethodCreationDateCol  = "creation_date"
	UserAuthMethodChangeDateCol    = "change_date"
	UserAuthMethodSequenceCol      = "sequence"
	UserAuthMethodResourceOwnerCol = "resource_owner"
	UserAuthMethodInstanceIDCol    = "instance_id"
	UserAuthMethodStateCol         = "state"
	UserAuthMethodNameCol          = "name"
)

type UserAuthMethodProjection struct {
	crdb.StatementHandler
}

func NewUserAuthMethodProjection(ctx context.Context, config crdb.StatementHandlerConfig) *UserAuthMethodProjection {
	p := new(UserAuthMethodProjection)
	config.ProjectionName = UserAuthMethodTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserAuthMethodUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodTypeCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodTokenIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserAuthMethodChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserAuthMethodSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(UserAuthMethodStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserAuthMethodResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserAuthMethodNameCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(UserAuthMethodInstanceIDCol, UserAuthMethodUserIDCol, UserAuthMethodTypeCol, UserAuthMethodTokenIDCol),
			crdb.WithIndex(crdb.NewIndex("ro_idx", []string{UserAuthMethodResourceOwnerCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *UserAuthMethodProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		},
	}
}

func (p *UserAuthMethodProjection) reduceInitAuthMethod(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCol(UserAuthMethodTokenIDCol, tokenID),
			handler.NewCol(UserAuthMethodCreationDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCol(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCol(UserAuthMethodUserIDCol, event.Aggregate().ID),
			handler.NewCol(UserAuthMethodSequenceCol, event.Sequence()),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateNotReady),
			handler.NewCol(UserAuthMethodTypeCol, methodType),
			handler.NewCol(UserAuthMethodNameCol, ""),
		},
	), nil
}

func (p *UserAuthMethodProjection) reduceActivateEvent(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCol(UserAuthMethodSequenceCol, event.Sequence()),
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

func (p *UserAuthMethodProjection) reduceRemoveAuthMethod(event eventstore.Event) (*handler.Statement, error) {
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
