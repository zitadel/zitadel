package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	UserAuthMethodTable = "projections.user_auth_methods4"

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
	UserAuthMethodOwnerRemovedCol  = "owner_removed"
)

type userAuthMethodProjection struct{}

func newUserAuthMethodProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(userAuthMethodProjection))
}

func (*userAuthMethodProjection) Name() string {
	return UserAuthMethodTable
}

func (*userAuthMethodProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(UserAuthMethodUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserAuthMethodTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(UserAuthMethodTokenIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserAuthMethodCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserAuthMethodChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserAuthMethodSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(UserAuthMethodStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(UserAuthMethodResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(UserAuthMethodInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserAuthMethodNameCol, handler.ColumnTypeText),
			handler.NewColumn(UserAuthMethodOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(UserAuthMethodInstanceIDCol, UserAuthMethodUserIDCol, UserAuthMethodTypeCol, UserAuthMethodTokenIDCol),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{UserAuthMethodResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{UserAuthMethodOwnerRemovedCol})),
		),
	)
}

func (p *userAuthMethodProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  user.HumanOTPSMSAddedType,
					Reduce: p.reduceAddAuthMethod,
				},
				{
					Event:  user.HumanOTPEmailAddedType,
					Reduce: p.reduceAddAuthMethod,
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
				{
					Event:  user.HumanOTPSMSRemovedType,
					Reduce: p.reduceRemoveAuthMethod,
				},
				{
					Event:  user.HumanPhoneRemovedType,
					Reduce: p.reduceRemoveAuthMethod,
				},
				{
					Event:  user.UserV1PhoneRemovedType,
					Reduce: p.reduceRemoveAuthMethod,
				},
				{
					Event:  user.HumanOTPEmailRemovedType,
					Reduce: p.reduceRemoveAuthMethod,
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
					Reduce: reduceInstanceRemovedHelper(UserAuthMethodInstanceIDCol),
				},
			},
		},
	}
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
		methodType = domain.UserAuthMethodTypeTOTP
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	return handler.NewUpsertStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodInstanceIDCol, nil),
			handler.NewCol(UserAuthMethodUserIDCol, nil),
			handler.NewCol(UserAuthMethodTypeCol, nil),
			handler.NewCol(UserAuthMethodTokenIDCol, nil),
		},
		[]handler.Column{
			handler.NewCol(UserAuthMethodTokenIDCol, tokenID),
			handler.NewCol(UserAuthMethodCreationDateCol, event.CreatedAt()),
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreatedAt()),
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
		methodType = domain.UserAuthMethodTypeTOTP

	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType})
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreatedAt()),
			handler.NewCol(UserAuthMethodSequenceCol, event.Sequence()),
			handler.NewCol(UserAuthMethodNameCol, name),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateReady),
		},
		[]handler.Condition{
			handler.NewCond(UserAuthMethodUserIDCol, event.Aggregate().ID),
			handler.NewCond(UserAuthMethodTypeCol, methodType),
			handler.NewCond(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCond(UserAuthMethodTokenIDCol, tokenID),
			handler.NewCond(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *userAuthMethodProjection) reduceAddAuthMethod(event eventstore.Event) (*handler.Statement, error) {
	var methodType domain.UserAuthMethodType
	switch event.(type) {
	case *user.HumanOTPSMSAddedEvent:
		methodType = domain.UserAuthMethodTypeOTPSMS
	case *user.HumanOTPEmailAddedEvent:
		methodType = domain.UserAuthMethodTypeOTPEmail
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-DS4g3", "reduce.wrong.event.type %v", []eventstore.EventType{user.HumanOTPSMSAddedType, user.HumanOTPEmailAddedType})
	}

	return handler.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(UserAuthMethodTokenIDCol, ""),
			handler.NewCol(UserAuthMethodCreationDateCol, event.CreatedAt()),
			handler.NewCol(UserAuthMethodChangeDateCol, event.CreatedAt()),
			handler.NewCol(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCol(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
			handler.NewCol(UserAuthMethodUserIDCol, event.Aggregate().ID),
			handler.NewCol(UserAuthMethodSequenceCol, event.Sequence()),
			handler.NewCol(UserAuthMethodStateCol, domain.MFAStateReady),
			handler.NewCol(UserAuthMethodTypeCol, methodType),
			handler.NewCol(UserAuthMethodNameCol, ""),
		},
	), nil
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
		methodType = domain.UserAuthMethodTypeTOTP
	case *user.HumanOTPSMSRemovedEvent,
		*user.HumanPhoneRemovedEvent:
		methodType = domain.UserAuthMethodTypeOTPSMS
	case *user.HumanOTPEmailRemovedEvent:
		methodType = domain.UserAuthMethodTypeOTPEmail

	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-f92f", "reduce.wrong.event.type %v",
			[]eventstore.EventType{user.HumanPasswordlessTokenAddedType, user.HumanU2FTokenAddedType, user.HumanMFAOTPRemovedType,
				user.HumanOTPSMSRemovedType, user.HumanPhoneRemovedType, user.HumanOTPEmailRemovedType})
	}
	conditions := []handler.Condition{
		handler.NewCond(UserAuthMethodUserIDCol, event.Aggregate().ID),
		handler.NewCond(UserAuthMethodTypeCol, methodType),
		handler.NewCond(UserAuthMethodResourceOwnerCol, event.Aggregate().ResourceOwner),
		handler.NewCond(UserAuthMethodInstanceIDCol, event.Aggregate().InstanceID),
	}
	if tokenID != "" {
		conditions = append(conditions, handler.NewCond(UserAuthMethodTokenIDCol, tokenID))
	}
	return handler.NewDeleteStatement(
		event,
		conditions,
	), nil
}

func (p *userAuthMethodProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-FwDZ8", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserAuthMethodInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(UserAuthMethodResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
