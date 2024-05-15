package handler

import (
	"context"
	"time"

	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	userSessionTable = "auth.user_sessions"

	userAgentIDCol                  = "user_agent_id"
	userIDCol                       = "user_id"
	instanceIDCol                   = "instance_id"
	creationDateCol                 = "creation_date"
	changeDateCol                   = "change_date"
	resourceOwnerCol                = "resource_owner"
	sequenceCol                     = "sequence"
	passwordVerificationCol         = "password_verification"
	stateCol                        = "state"
	secondFactorVerificationCol     = "second_factor_verification"
	secondFactorVerificationTypeCol = "second_factor_verification_type"
	multiFactorVerificationCol      = "multi_factor_verification"
	multiFactorVerificationTypeCol  = "multi_factor_verification_type"
	passwordlessVerificationCol     = "passwordless_verification"
	externalLoginVerificationCol    = "external_login_verification"
	selectedIDPConfigIDCol          = "selected_idp_config_id"
)

type UserSession struct {
	queries *query2.Queries
	view    *auth_view.View
	es      handler.EventStore
}

var _ handler.Projection = (*UserSession)(nil)

func newUserSession(
	ctx context.Context,
	config handler.Config,
	view *auth_view.View,
	queries *query2.Queries,
) *handler.Handler {
	return handler.NewHandler(
		ctx,
		&config,
		&UserSession{
			queries: queries,
			view:    view,
			es:      config.Eventstore,
		},
	)
}

// Name implements [handler.Projection]
func (*UserSession) Name() string {
	return userSessionTable
}

// Reducers implements [handler.Projection]
func (s *UserSession) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserV1PasswordCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1PasswordCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1MFAOTPCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1MFAOTPCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1SignedOutType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserIDPLoginCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanMFAOTPCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanMFAOTPCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanU2FTokenCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanU2FTokenCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordlessTokenCheckSucceededType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordlessTokenCheckFailedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanSignedOutType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1PasswordChangedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserV1MFAOTPRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserLockedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserDeactivatedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordChangedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanMFAOTPRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserIDPLinkRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserIDPLinkCascadeRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanPasswordlessTokenRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanU2FTokenRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanRegisteredType,
					Reduce: s.Reduce,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: s.Reduce,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: s.Reduce,
				},
			},
		},
	}
}

func sessionColumns(event eventstore.Event, columns ...handler.Column) ([]handler.Column, error) {
	userAgent, err := agentIDFromSession(event)
	if err != nil {
		return nil, err
	}
	return append([]handler.Column{
		handler.NewCol(userAgentIDCol, userAgent),
		handler.NewCol(userIDCol, event.Aggregate().ID),
		handler.NewCol(instanceIDCol, event.Aggregate().InstanceID),
		handler.NewCol(creationDateCol, handler.OnlySetValueOnInsert(creationDateCol, event.CreatedAt())),
		handler.NewCol(changeDateCol, event.CreatedAt()),
		handler.NewCol(resourceOwnerCol, event.Aggregate().ResourceOwner),
		handler.NewCol(sequenceCol, event.Sequence()),
	}, columns...), nil
}

func (u *UserSession) Reduce(event eventstore.Event) (_ *handler.Statement, err error) {
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		columns, err := sessionColumns(event,
			handler.NewCol(passwordVerificationCol, event.CreatedAt()),
			handler.NewCol(stateCol, domain.UserSchemaStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType:
		columns, err := sessionColumns(event,
			handler.NewCol(passwordVerificationCol, time.Time{}),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserV1MFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType:
		columns, err := sessionColumns(event,
			handler.NewCol(secondFactorVerificationCol, event.CreatedAt()),
			handler.NewCol(secondFactorVerificationTypeCol, domain.MFATypeTOTP),
			handler.NewCol(stateCol, domain.UserSchemaStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserV1MFAOTPCheckFailedType,
		user.HumanMFAOTPCheckFailedType,
		user.HumanU2FTokenCheckFailedType:
		columns, err := sessionColumns(event,
			handler.NewCol(secondFactorVerificationCol, time.Time{}),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		columns, err := sessionColumns(event,
			handler.NewCol(passwordlessVerificationCol, time.Time{}),
			handler.NewCol(passwordVerificationCol, time.Time{}),
			handler.NewCol(secondFactorVerificationCol, time.Time{}),
			handler.NewCol(secondFactorVerificationTypeCol, domain.MFALevelNotSetUp),
			handler.NewCol(multiFactorVerificationCol, time.Time{}),
			handler.NewCol(multiFactorVerificationTypeCol, domain.MFALevelNotSetUp),
			handler.NewCol(externalLoginVerificationCol, time.Time{}),
			handler.NewCol(stateCol, domain.UserSessionStateTerminated),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(externalLoginVerificationCol, event.CreatedAt()),
			handler.NewCol(selectedIDPConfigIDCol, data.SelectedIDPConfigID),
			handler.NewCol(stateCol, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.HumanU2FTokenCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(secondFactorVerificationCol, event.CreatedAt()),
			handler.NewCol(secondFactorVerificationTypeCol, domain.MFATypeU2F),
			handler.NewCol(stateCol, domain.UserSchemaStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.HumanPasswordlessTokenCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(passwordlessVerificationCol, event.CreatedAt()),
			handler.NewCol(multiFactorVerificationCol, event.CreatedAt()),
			handler.NewCol(multiFactorVerificationTypeCol, domain.MFATypeU2FUserVerification),
			handler.NewCol(stateCol, domain.UserSchemaStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.HumanPasswordlessTokenCheckFailedType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(passwordlessVerificationCol, time.Time{}),
			handler.NewCol(multiFactorVerificationCol, time.Time{}),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:2], columns), nil
	case user.UserLockedType,
		user.UserDeactivatedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(passwordlessVerificationCol, time.Time{}),
				handler.NewCol(passwordVerificationCol, time.Time{}),
				handler.NewCol(secondFactorVerificationCol, time.Time{}),
				handler.NewCol(secondFactorVerificationTypeCol, domain.MFALevelNotSetUp),
				handler.NewCol(multiFactorVerificationCol, time.Time{}),
				handler.NewCol(multiFactorVerificationTypeCol, domain.MFALevelNotSetUp),
				handler.NewCol(externalLoginVerificationCol, time.Time{}),
				handler.NewCol(stateCol, domain.UserSessionStateTerminated),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol(sequenceCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
				handler.Not(handler.NewCond(stateCol, domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		userAgent, err := agentIDFromSession(event)
		if err != nil {
			return nil, err
		}
		return handler.NewMultiStatement(event,
			handler.AddUpdateStatement(
				[]handler.Column{
					handler.NewCol(passwordVerificationCol, event.CreatedAt()),
					handler.NewCol(changeDateCol, event.CreatedAt()),
					handler.NewCol(sequenceCol, event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
					handler.NewCond(userIDCol, event.Aggregate().ID),
					handler.NewCond(userAgentIDCol, userAgent),
				}),
			handler.AddUpdateStatement(
				[]handler.Column{
					handler.NewCol(passwordVerificationCol, time.Time{}),
					handler.NewCol(changeDateCol, event.CreatedAt()),
					handler.NewCol(sequenceCol, event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
					handler.NewCond(userIDCol, event.Aggregate().ID),
					handler.Not(handler.NewCond(userAgentIDCol, userAgent)),
					handler.Not(handler.NewCond(stateCol, domain.UserSessionStateTerminated)),
				}),
		), nil
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType,
		user.HumanU2FTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(secondFactorVerificationCol, time.Time{}),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol(sequenceCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
				handler.Not(handler.NewCond(stateCol, domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserIDPLinkRemovedType,
		user.UserIDPLinkCascadeRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(externalLoginVerificationCol, time.Time{}),
				handler.NewCol(selectedIDPConfigIDCol, ""),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol(sequenceCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
				handler.Not(handler.NewCond(selectedIDPConfigIDCol, "")),
			},
		), nil
	case user.HumanPasswordlessTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(passwordlessVerificationCol, time.Time{}),
				handler.NewCol(multiFactorVerificationCol, time.Time{}),
				handler.NewCol(changeDateCol, event.CreatedAt()),
				handler.NewCol(sequenceCol, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
				handler.Not(handler.NewCond(stateCol, domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(userIDCol, event.Aggregate().ID),
			},
		), nil
	case user.HumanRegisteredType:
		columns, err := sessionColumns(event,
			handler.NewCol(stateCol, domain.UserSessionStateActive),
			handler.NewCol(passwordVerificationCol, event.CreatedAt()),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewCreateStatement(event,
			columns,
		), nil
	case instance.InstanceRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(instanceIDCol, event.Aggregate().InstanceID),
				handler.NewCond(resourceOwnerCol, event.Aggregate().ResourceOwner),
			},
		), nil
	default:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return nil
		}), nil
	}
}
