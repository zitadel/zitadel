package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserAuthMethodProjection struct {
	crdb.StatementHandler
}

const (
	UserAuthMethodTable = "zitadel.projections.user_auth_method"
)

func NewUserAuthMethodProjection(ctx context.Context, config crdb.StatementHandlerConfig) *UserAuthMethodProjection {
	p := &UserAuthMethodProjection{}
	config.ProjectionName = UserAuthMethodTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

const (
	UserAuthMethodTokenIDCol       = "token_id"
	UserAuthMethodCreationDateCol  = "creation_date"
	UserAuthMethodChangeUseCol     = "change_date"
	UserAuthMethodResourceOwnerCol = "resource_owner"
	UserAuthMethodUserIDCol        = "user_id"
	UserAuthMethodSequenceCol      = "sequence"
	UserAuthMethodNameCol          = "name"
	UserAuthMethodStateCol         = "state"
	UserAuthMethodTypeCol          = "type"
)

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
					Event:  user.HumanU2FTokenVerifiedType,
					Reduce: p.reduceRemoveAuthMethod,
				},
				{
					Event:  user.HumanMFAOTPVerifiedType,
					Reduce: p.reduceRemoveAuthMethod,
				},
			},
		},
	}
}

func (p *UserAuthMethodProjection) reduceInitAuthMethod(event eventstore.Event) (*handler.Statement, error) {
	var tokenID string
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
		logging.LogWithFields("PROJE-9j3f", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-f92f", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodTokenIDCol, tokenID),
			handler.NewCol(UserAuthMethodCreationDateCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodChangeUseCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
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
	case *user.HumanOTPAddedEvent:
		methodType = domain.UserAuthMethodTypeOTP

	default:
		logging.LogWithFields("PROJE-9j3f", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-f92f", "reduce.wrong.event.type")
	}
	conditions := []handler.Condition{
		handler.NewCond(UserAuthMethodUserIDCol, event.Aggregate().ID),
		handler.NewCond(UserAuthMethodTypeCol, methodType),
		handler.NewCond(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
	}
	if tokenID != "" {
		conditions = append(conditions, handler.NewCond(UserAuthMethodTokenIDCol, tokenID))
	}
	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodChangeUseCol, event.CreationDate()),
			handler.NewCol(UserAuthMethodSequenceCol, event.Sequence()),
			handler.NewCol(UserAuthMethodNameCol, name),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateReady),
		},
		conditions,
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
		logging.LogWithFields("PROJE-9j3f", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-f92f", "reduce.wrong.event.type")
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
