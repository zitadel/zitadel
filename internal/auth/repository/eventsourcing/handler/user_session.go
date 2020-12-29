package handler

import (
	"github.com/caos/logging"
	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userSessionTable = "auth.user_sessions"
)

type UserSession struct {
	handler
	userEvents   *user_events.UserEventstore
	subscription *eventstore.Subscription
}

func newUserSession(
	handler handler,
	userEvents *user_events.UserEventstore,
) *UserSession {
	h := &UserSession{
		handler:    handler,
		userEvents: userEvents,
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

func (_ *UserSession) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.UserAggregate}
}

func (u *UserSession) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := u.view.GetLatestUserSessionSequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *UserSession) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSessionSequence("")
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence.CurrentSequence), nil
}

func (u *UserSession) Reduce(event *models.Event) (err error) {
	var sessions []*view_model.UserSessionView
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.HumanPasswordCheckSucceeded,
		es_model.HumanExternalLoginCheckSucceeded,
		es_model.MFAOTPCheckSucceeded,
		es_model.HumanMFAOTPCheckSucceeded,
		es_model.HumanMFAU2FTokenCheckSucceeded,
		es_model.HumanPasswordlessTokenCheckSucceeded,
		es_model.UserPasswordCheckFailed,
		es_model.HumanPasswordCheckFailed,
		es_model.MFAOTPCheckFailed,
		es_model.HumanMFAOTPCheckFailed,
		es_model.HumanMFAU2FTokenCheckFailed,
		es_model.HumanPasswordlessTokenCheckFailed,
		es_model.SignedOut,
		es_model.HumanSignedOut,
		es_model.UserLocked,
		es_model.UserDeactivated:

		var session *view_model.UserSessionView
		session, err = u.sessionFromEvent(event)
		sessions = append(sessions, session)
	case es_model.UserProfileChanged,
		es_model.HumanProfileChanged,
		es_model.UserUserNameChanged,
		es_model.DomainClaimed:

		sessions, err = u.UserDataChanged(event)
	case es_model.UserRemoved:
		return u.view.DeleteUserSessions(event.AggregateID, event)
	default:
		return u.view.ProcessedUserSessionSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUserSessions(sessions, event)
}

func (u *UserSession) UserDataChanged(event *models.Event) ([]*view_model.UserSessionView, error) {
	sessions, err := u.view.UserSessionsByUserID(event.AggregateID)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, u.view.ProcessedUserSessionSequence(event)
	}
	for i := len(sessions) - 1; i > 0; i-- {
		if sessions[i].State != int32(req_model.UserSessionStateActive) {
			copy(sessions[i:], sessions[i+1:])
			sessions[len(sessions)-1] = nil
			sessions = sessions[:len(sessions)-1]
			continue
		}
		if err := sessions[i].AppendEvent(event); err != nil {
			return nil, err
		}
		if err := u.fillUserInfo(sessions[i], event.AggregateID); err != nil {
			return nil, err
		}
	}
	return sessions, nil
}

func (u *UserSession) sessionFromEvent(event *models.Event) (*view_model.UserSessionView, error) {
	eventData, err := view_model.UserSessionFromEvent(event)
	if err != nil {
		return nil, err
	}
	session, err := u.view.UserSessionByIDs(eventData.UserAgentID, event.AggregateID)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		session = &view_model.UserSessionView{
			CreationDate:  event.CreationDate,
			ResourceOwner: event.ResourceOwner,
			UserAgentID:   eventData.UserAgentID,
			UserID:        event.AggregateID,
			State:         int32(req_model.UserSessionStateInitiated),
		}
	}

	if err = session.AppendEvent(event); err != nil {
		return nil, err
	}

	if err = u.fillUserInfo(session, event.AggregateID); err != nil {
		return nil, err
	}

	return session, nil
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
	if session.State != int32(req_model.UserSessionStateActive) {
		return nil
	}
	user, err := u.view.UserByID(id)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	return nil
}
