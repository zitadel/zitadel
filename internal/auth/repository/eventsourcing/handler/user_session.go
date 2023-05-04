package handler

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/zitadel/zitadel/internal/eventstore/v1/sdk"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/org/repository/view"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
)

type UserSession struct {
	handler
	subscription *v1.Subscription
	queries      *query2.Queries
}

func newUserSession(ctx context.Context, handler handler, queries *query2.Queries) *UserSession {
	h := &UserSession{
		handler: handler,
		queries: queries,
	}

	h.subscribe(ctx)

	return h
}

func (u *UserSession) subscribe(ctx context.Context) {
	u.subscription = u.es.Subscribe(u.AggregateTypes()...)
	go func() {
		for event := range u.subscription.Events {
			query.ReduceEvent(ctx, u, event)
		}
	}()
}

func (u *UserSession) ViewModel() string {
	return userSessionTable
}

func (u *UserSession) Subscription() *v1.Subscription {
	return u.subscription
}

func (_ *UserSession) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{user.AggregateType, org.AggregateType, instance.AggregateType}
}

func (u *UserSession) CurrentSequence(ctx context.Context, instanceID string) (uint64, error) {
	sequence, err := u.view.GetLatestUserSessionSequence(ctx, instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserSession) EventQuery(ctx context.Context, instanceIDs []string) (*models.SearchQuery, error) {
	sequences, err := u.view.GetLatestUserSessionSequences(ctx, instanceIDs)
	if err != nil {
		return nil, err
	}
	return newSearchQuery(sequences, u.AggregateTypes(), instanceIDs), nil
}

func (u *UserSession) Reduce(event *models.Event) (err error) {
	var session *view_model.UserSessionView
	switch eventstore.EventType(event.Type) {
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
		session, err = u.view.UserSessionByIDs(eventData.UserAgentID, event.AggregateID, event.InstanceID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			session = &view_model.UserSessionView{
				CreationDate:  event.CreationDate,
				ResourceOwner: event.ResourceOwner,
				UserAgentID:   eventData.UserAgentID,
				UserID:        event.AggregateID,
				State:         int32(domain.UserSessionStateActive),
				InstanceID:    event.InstanceID,
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
		sessions, err := u.view.UserSessionsByUserID(event.AggregateID, event.InstanceID)
		if err != nil {
			return err
		}
		if len(sessions) == 0 {
			return u.view.ProcessedUserSessionSequence(event)
		}
		for _, session := range sessions {
			if err := session.AppendEvent(event); err != nil {
				return err
			}
			if err := u.fillUserInfo(session); err != nil {
				return err
			}
		}
		return u.view.PutUserSessions(sessions, event)
	case org.OrgDomainPrimarySetEventType:
		return u.fillLoginNamesOnOrgUsers(event)
	case user.UserRemovedType:
		return u.view.DeleteUserSessions(event.AggregateID, event.InstanceID, event)
	case instance.InstanceRemovedEventType:
		return u.view.DeleteInstanceUserSessions(event)
	case org.OrgRemovedEventType:
		return u.view.DeleteOrgUserSessions(event)
	default:
		return u.view.ProcessedUserSessionSequence(event)
	}
}

func (u *UserSession) OnError(event *models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in user session handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserSessionFailedEvent, u.view.ProcessedUserSessionFailedEvent, u.view.ProcessedUserSessionSequence, u.errorCountUntilSkip)
}

func (u *UserSession) OnSuccess(instanceIDs []string) error {
	return spooler.HandleSuccess(u.view.UpdateUserSessionSpoolerRunTimestamp, instanceIDs)
}

func (u *UserSession) updateSession(session *view_model.UserSessionView, event *models.Event) error {
	if err := session.AppendEvent(event); err != nil {
		return err
	}
	if err := u.fillUserInfo(session); err != nil {
		return err
	}
	return u.view.PutUserSession(session, event)
}

func (u *UserSession) fillUserInfo(session *view_model.UserSessionView) error {
	user, err := u.view.UserByID(session.UserID, session.InstanceID)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	session.AvatarKey = user.AvatarKey
	return nil
}

func (u *UserSession) fillLoginNamesOnOrgUsers(event *models.Event) error {
	sessions, err := u.view.UserSessionsByOrgID(event.ResourceOwner, event.InstanceID)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		return u.view.ProcessedUserSessionSequence(event)
	}
	userLoginMustBeDomain, primaryDomain, err := u.loginNameInformation(context.Background(), event.ResourceOwner, event.InstanceID)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	for _, session := range sessions {
		session.LoginName = session.UserName + "@" + primaryDomain
	}
	return u.view.PutUserSessions(sessions, event)
}

func (u *UserSession) loginNameInformation(ctx context.Context, orgID string, instanceID string) (userLoginMustBeDomain bool, primaryDomain string, err error) {
	org, err := u.getOrgByID(ctx, orgID, instanceID)
	if err != nil {
		return false, "", err
	}
	if org.DomainPolicy != nil {
		return org.DomainPolicy.UserLoginMustBeDomain, org.GetPrimaryDomain().Domain, nil
	}
	policy, err := u.queries.DefaultDomainPolicy(withInstanceID(ctx, org.InstanceID))
	if err != nil {
		return false, "", err
	}
	return policy.UserLoginMustBeDomain, org.GetPrimaryDomain().Domain, nil
}

func (u *UserSession) getOrgByID(ctx context.Context, orgID, instanceID string) (*org_model.Org, error) {
	orgQuery, err := view.OrgByIDQuery(orgID, instanceID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, orgQuery)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}
