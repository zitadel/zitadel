package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/repository/user"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
)

type UserSession struct {
	handler
	subscription *v1.Subscription
}

func newUserSession(
	handler handler,
) *UserSession {
	h := &UserSession{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (k *UserSession) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
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
	return []models.AggregateType{user.AggregateType}
}

func (u *UserSession) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := u.view.GetLatestUserSessionSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserSession) EventQuery() (*models.SearchQuery, error) {
	sequences, err := u.view.GetLatestUserSessionSequences()
	if err != nil {
		return nil, err
	}
	query := models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(u.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(0).
		IgnoredInstanceIDsFilter(instances...).
		SearchQuery(), nil
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
			if err := u.fillUserInfo(session, event.AggregateID); err != nil {
				return err
			}
		}
		return u.view.PutUserSessions(sessions, event)
	case user.UserRemovedType:
		return u.view.DeleteUserSessions(event.AggregateID, event.InstanceID, event)
	default:
		return u.view.ProcessedUserSessionSequence(event)
	}
}

func (u *UserSession) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-sdfw3s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user session handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserSessionFailedEvent, u.view.ProcessedUserSessionFailedEvent, u.view.ProcessedUserSessionSequence, u.errorCountUntilSkip)
}

func (u *UserSession) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserSessionSpoolerRunTimestamp)
}

func (u *UserSession) updateSession(session *view_model.UserSessionView, event *models.Event) error {
	if err := session.AppendEvent(event); err != nil {
		return err
	}
	if err := u.fillUserInfo(session, event.AggregateID); err != nil {
		return err
	}
	return u.view.PutUserSession(session, event)
}

func (u *UserSession) fillUserInfo(session *view_model.UserSessionView, id string) error {
	user, err := u.view.UserByID(id, session.InstanceID)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	session.AvatarKey = user.AvatarKey
	return nil
}
