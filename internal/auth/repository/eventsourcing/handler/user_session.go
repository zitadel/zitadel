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
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (u *UserSession) Reduce(event eventstore.Event) (_ *handler.Statement, err error) {
	switch event.Type() {
	case user.UserV1PasswordCheckSucceededType,
		user.UserV1PasswordCheckFailedType,
		user.UserV1MFAOTPCheckSucceededType,
		user.UserV1MFAOTPCheckFailedType,
		user.UserV1SignedOutType,
		user.HumanPasswordCheckSucceededType,
		user.HumanPasswordCheckFailedType,
		user.UserIDPLoginCheckSucceededType,
		user.HumanMFAOTPCheckSucceededType,
		user.HumanMFAOTPCheckFailedType,
		user.HumanU2FTokenCheckSucceededType,
		user.HumanU2FTokenCheckFailedType,
		user.HumanPasswordlessTokenCheckSucceededType,
		user.HumanPasswordlessTokenCheckFailedType,
		user.HumanSignedOutType:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			var session *view_model.UserSessionView
			eventData, err := view_model.UserSessionFromEvent(event)
			if err != nil {
				return err
			}
			session, err = u.view.UserSessionByIDs(eventData.UserAgentID, event.Aggregate().ID, event.Aggregate().InstanceID)
			if err != nil {
				if !zerrors.IsNotFound(err) {
					return err
				}
				session = &view_model.UserSessionView{
					CreationDate:  event.CreatedAt(),
					ResourceOwner: event.Aggregate().ResourceOwner,
					UserAgentID:   eventData.UserAgentID,
					UserID:        event.Aggregate().ID,
					State:         int32(domain.UserSessionStateActive),
					InstanceID:    event.Aggregate().InstanceID,
				}
			}
			return u.updateSession(session, event)
		}), nil
	case user.UserLockedType,
		user.UserDeactivatedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("passwordless_verification", time.Time{}),
				handler.NewCol("password_verification", time.Time{}),
				handler.NewCol("second_factor_verification", time.Time{}),
				handler.NewCol("second_factor_verification_type", domain.MFALevelNotSetUp),
				handler.NewCol("multi_factor_verification", time.Time{}),
				handler.NewCol("multi_factor_verification_type", domain.MFALevelNotSetUp),
				handler.NewCol("external_login_verification", time.Time{}),
				handler.NewCol("state", domain.UserSessionStateTerminated),
				handler.NewCol("change_date", event.CreatedAt()),
				handler.NewCol("sequence", event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond("instance_id", event.Aggregate().InstanceID),
				handler.NewCond("user_id", event.Aggregate().ID),
				handler.Not(handler.NewCond("state", domain.UserSessionStateTerminated)),
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
					handler.NewCol("password_verification", event.CreatedAt()),
					handler.NewCol("change_date", event.CreatedAt()),
					handler.NewCol("sequence", event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond("instance_id", event.Aggregate().InstanceID),
					handler.NewCond("user_id", event.Aggregate().ID),
					handler.NewCond("user_agent_id", userAgent),
				}),
			handler.AddUpdateStatement(
				[]handler.Column{
					handler.NewCol("password_verification", time.Time{}),
					handler.NewCol("change_date", event.CreatedAt()),
					handler.NewCol("sequence", event.Sequence()),
				},
				[]handler.Condition{
					handler.NewCond("instance_id", event.Aggregate().InstanceID),
					handler.NewCond("user_id", event.Aggregate().ID),
					handler.Not(handler.NewCond("user_agent_id", userAgent)),
					handler.Not(handler.NewCond("state", domain.UserSessionStateTerminated)),
				}),
		), nil
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType,
		user.HumanU2FTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("second_factor_verification", time.Time{}),
				handler.NewCol("change_date", event.CreatedAt()),
				handler.NewCol("sequence", event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond("instance_id", event.Aggregate().InstanceID),
				handler.NewCond("user_id", event.Aggregate().ID),
				handler.Not(handler.NewCond("state", domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserIDPLinkRemovedType,
		user.UserIDPLinkCascadeRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("external_login_verification", time.Time{}),
				handler.NewCol("selected_idp_config_id", ""),
				handler.NewCol("change_date", event.CreatedAt()),
				handler.NewCol("sequence", event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond("instance_id", event.Aggregate().InstanceID),
				handler.NewCond("user_id", event.Aggregate().ID),
				handler.Not(handler.NewCond("selected_idp_config_id", "")),
			},
		), nil
	case user.HumanPasswordlessTokenRemovedType:
		return handler.NewUpdateStatement(event,
			[]handler.Column{
				handler.NewCol("passwordless_verification", time.Time{}),
				handler.NewCol("multi_factor_verification", time.Time{}),
				handler.NewCol("change_date", event.CreatedAt()),
				handler.NewCol("sequence", event.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond("instance_id", event.Aggregate().InstanceID),
				handler.NewCond("user_id", event.Aggregate().ID),
				handler.Not(handler.NewCond("state", domain.UserSessionStateTerminated)),
			},
		), nil
	case user.UserRemovedType:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return u.view.DeleteUserSessions(event.Aggregate().ID, event.Aggregate().InstanceID)
		}), nil
	case user.HumanRegisteredType:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			eventData, err := view_model.UserSessionFromEvent(event)
			if err != nil {
				return err
			}
			session := &view_model.UserSessionView{
				CreationDate:         event.CreatedAt(),
				ResourceOwner:        event.Aggregate().ResourceOwner,
				UserAgentID:          eventData.UserAgentID,
				UserID:               event.Aggregate().ID,
				State:                int32(domain.UserSessionStateActive),
				InstanceID:           event.Aggregate().InstanceID,
				PasswordVerification: event.CreatedAt(),
			}
			return u.updateSession(session, event)
		}), nil
	case instance.InstanceRemovedEventType:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return u.view.DeleteInstanceUserSessions(event.Aggregate().InstanceID)
		}), nil
	case org.OrgRemovedEventType:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return u.view.DeleteOrgUserSessions(event)
		}), nil
	default:
		return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
			return nil
		}), nil
	}
}

func (u *UserSession) updateSession(session *view_model.UserSessionView, event eventstore.Event) error {
	if err := session.AppendEvent(event); err != nil {
		return err
	}
	return u.view.PutUserSession(session)
}
