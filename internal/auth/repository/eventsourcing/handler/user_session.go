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
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
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
		handler.NewCol(view_model.UserSessionKeyUserAgentID, userAgent),
		handler.NewCol(view_model.UserSessionKeyUserID, event.Aggregate().ID),
		handler.NewCol(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
		handler.NewCol(view_model.UserSessionKeyCreationDate, handler.OnlySetValueOnInsert(userSessionTable, event.CreatedAt())),
		handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
		handler.NewCol(view_model.UserSessionKeyResourceOwner, event.Aggregate().ResourceOwner),
		handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
	}, columns...), nil
}

func (u *UserSession) Reduce(event eventstore.Event) (_ *handler.Statement, err error) {
	// in case anything needs to be change here check if appendEvent function needs the change as well
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.HumanPasswordCheckSucceededType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyPasswordVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserV1PasswordCheckFailedType,
		user.HumanPasswordCheckFailedType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyPasswordVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserV1MFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeySecondFactorVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeySecondFactorVerificationType, domain.MFATypeTOTP),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserV1MFAOTPCheckFailedType,
		user.HumanMFAOTPCheckFailedType,
		user.HumanU2FTokenCheckFailedType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeySecondFactorVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserV1SignedOutType,
		user.HumanSignedOutType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyPasswordlessVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyPasswordVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeySecondFactorVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeySecondFactorVerificationType, domain.MFALevelNotSetUp),
			handler.NewCol(view_model.UserSessionKeyMultiFactorVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyMultiFactorVerificationType, domain.MFALevelNotSetUp),
			handler.NewCol(view_model.UserSessionKeyExternalLoginVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateTerminated),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserIDPLoginCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyExternalLoginVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeySelectedIDPConfigID, data.SelectedIDPConfigID),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.HumanU2FTokenCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeySecondFactorVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeySecondFactorVerificationType, domain.MFATypeU2F),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.HumanPasswordlessTokenCheckSucceededType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyPasswordlessVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeyMultiFactorVerification, event.CreatedAt()),
			handler.NewCol(view_model.UserSessionKeyMultiFactorVerificationType, domain.MFATypeU2FUserVerification),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.HumanPasswordlessTokenCheckFailedType:
		data := new(es_model.AuthRequest)
		err := data.SetData(event)
		if err != nil {
			return nil, err
		}
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyPasswordlessVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyMultiFactorVerification, time.Time{}),
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
		)
		if err != nil {
			return nil, err
		}
		return handler.NewUpsertStatement(event, columns[0:3], columns), nil
	case user.UserLockedType,
		user.UserDeactivatedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(view_model.UserSessionKeyPasswordlessVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyPasswordVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeySecondFactorVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeySecondFactorVerificationType, domain.MFALevelNotSetUp),
				handler.NewCol(view_model.UserSessionKeyMultiFactorVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyMultiFactorVerificationType, domain.MFALevelNotSetUp),
				handler.NewCol(view_model.UserSessionKeyExternalLoginVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateTerminated),
				handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
				handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
				handler.Not(handler.NewCond(view_model.UserSessionKeyState, domain.UserSessionStateTerminated)),
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
					handler.NewCol(view_model.UserSessionKeyPasswordVerification, event.CreatedAt()),
					handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
					handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
					handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
					handler.NewCond(view_model.UserSessionKeyUserAgentID, userAgent),
				}),
			handler.AddUpdateStatement(
				[]handler.Column{
					handler.NewCol(view_model.UserSessionKeyPasswordVerification, time.Time{}),
					handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
					handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
					handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
					handler.Not(handler.NewCond(view_model.UserSessionKeyUserAgentID, userAgent)),
					handler.Not(handler.NewCond(view_model.UserSessionKeyState, domain.UserSessionStateTerminated)),
				}),
		), nil
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType,
		user.HumanU2FTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(view_model.UserSessionKeySecondFactorVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
				handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
				handler.Not(handler.NewCond(view_model.UserSessionKeyState, domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserIDPLinkRemovedType,
		user.UserIDPLinkCascadeRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(view_model.UserSessionKeyExternalLoginVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeySelectedIDPConfigID, ""),
				handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
				handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
				handler.Not(handler.NewCond(view_model.UserSessionKeySelectedIDPConfigID, "")),
			},
		), nil
	case user.HumanPasswordlessTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol(view_model.UserSessionKeyPasswordlessVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyMultiFactorVerification, time.Time{}),
				handler.NewCol(view_model.UserSessionKeyChangeDate, event.CreatedAt()),
				handler.NewCol(view_model.UserSessionKeySequence, event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
				handler.Not(handler.NewCond(view_model.UserSessionKeyState, domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserRemovedType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyUserID, event.Aggregate().ID),
			},
		), nil
	case user.HumanRegisteredType:
		columns, err := sessionColumns(event,
			handler.NewCol(view_model.UserSessionKeyState, domain.UserSessionStateActive),
			handler.NewCol(view_model.UserSessionKeyPasswordVerification, event.CreatedAt()),
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
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
			},
		), nil
	case org.OrgRemovedEventType:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(view_model.UserSessionKeyInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(view_model.UserSessionKeyResourceOwner, event.Aggregate().ResourceOwner),
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}
