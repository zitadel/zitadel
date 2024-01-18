package handler

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/zitadel/zitadel/internal/org/repository/view"
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
					Event:  user.UserV1ProfileChangedType,
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
					Event:  user.HumanProfileChangedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanAvatarAddedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.HumanAvatarRemovedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserDomainClaimedType,
					Reduce: s.Reduce,
				},
				{
					Event:  user.UserUserNameChangedType,
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
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: s.Reduce,
				},
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
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		var session *view_model.UserSessionView
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
		case user.UserV1PasswordChangedType,
			user.UserV1MFAOTPRemovedType,
			user.UserV1ProfileChangedType,
			user.UserLockedType,
			user.UserDeactivatedType,
			user.HumanPasswordChangedType,
			user.HumanMFAOTPRemovedType,
			user.HumanProfileChangedType,
			user.HumanAvatarAddedType,
			user.HumanAvatarRemovedType,
			user.UserDomainClaimedType,
			user.UserUserNameChangedType,
			user.UserIDPLinkRemovedType,
			user.UserIDPLinkCascadeRemovedType,
			user.HumanPasswordlessTokenRemovedType,
			user.HumanU2FTokenRemovedType:
			sessions, err := u.view.UserSessionsByUserID(event.Aggregate().ID, event.Aggregate().InstanceID)
			if err != nil || len(sessions) == 0 {
				return err
			}
			if err = u.appendEventOnSessions(sessions, event); err != nil {
				return err
			}
			if err = u.view.PutUserSessions(sessions); err != nil {
				return err
			}
			return nil
		case org.OrgDomainPrimarySetEventType:
			return u.fillLoginNamesOnOrgUsers(event)
		case user.UserRemovedType:
			return u.view.DeleteUserSessions(event.Aggregate().ID, event.Aggregate().InstanceID)
		case instance.InstanceRemovedEventType:
			return u.view.DeleteInstanceUserSessions(event.Aggregate().InstanceID)
		case org.OrgRemovedEventType:
			return u.view.DeleteOrgUserSessions(event)
		default:
			return nil
		}
	}), nil
}

func (u *UserSession) appendEventOnSessions(sessions []*view_model.UserSessionView, event eventstore.Event) error {
	users := make(map[string]*view_model.UserView)
	usersByID := func(userID, instanceID string) (user *view_model.UserView, err error) {
		user, ok := users[userID+"-"+instanceID]
		if ok {
			return user, nil
		}
		users[userID+"-"+instanceID], err = u.view.UserByID(userID, instanceID)
		if err != nil {
			return nil, err
		}

		return users[userID+"-"+instanceID], nil
	}
	for _, session := range sessions {
		if err := session.AppendEvent(event); err != nil {
			return err
		}
		if err := u.fillUserInfo(session, usersByID); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserSession) updateSession(session *view_model.UserSessionView, event eventstore.Event) error {
	if err := session.AppendEvent(event); err != nil {
		return err
	}
	if err := u.fillUserInfo(session, u.view.UserByID); err != nil {
		return err
	}
	if err := u.view.PutUserSession(session); err != nil {
		return err
	}
	return nil
}

func (u *UserSession) fillUserInfo(session *view_model.UserSessionView, getUserByID func(userID, instanceID string) (*view_model.UserView, error)) error {
	user, err := getUserByID(session.UserID, session.InstanceID)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	session.AvatarKey = user.AvatarKey
	return nil
}

func (u *UserSession) fillLoginNamesOnOrgUsers(event eventstore.Event) error {
	sessions, err := u.view.UserSessionsByOrgID(event.Aggregate().ResourceOwner, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		return nil
	}
	userLoginMustBeDomain, primaryDomain, err := u.loginNameInformation(context.Background(), event.Aggregate().ResourceOwner, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	for _, session := range sessions {
		session.LoginName = session.UserName + "@" + primaryDomain
	}
	return u.view.PutUserSessions(sessions)
}

func (u *UserSession) loginNameInformation(ctx context.Context, orgID string, instanceID string) (userLoginMustBeDomain bool, primaryDomain string, err error) {
	org, err := u.getOrgByID(ctx, orgID, instanceID)
	if err != nil {
		return false, "", err
	}
	primaryDomain, err = org.GetPrimaryDomain()
	if err != nil {
		return false, "", err
	}
	if org.DomainPolicy != nil {
		return org.DomainPolicy.UserLoginMustBeDomain, primaryDomain, nil
	}
	policy, err := u.queries.DefaultDomainPolicy(authz.WithInstanceID(ctx, org.InstanceID))
	if err != nil {
		return false, "", err
	}
	return policy.UserLoginMustBeDomain, primaryDomain, nil
}

func (u *UserSession) getOrgByID(ctx context.Context, orgID, instanceID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, instanceID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	events, err := u.es.Filter(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = esOrg.AppendEvents(events...); err != nil {
		return nil, err
	}

	if esOrg.Sequence == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}
	return org_es_model.OrgToModel(esOrg), nil
}
