package handler

import (
	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserSession struct {
	handler
	userEvents *user_events.UserEventstore
}

const (
	userSessionTable = "auth.user_sessions"
)

func (u *UserSession) ViewModel() string {
	return userSessionTable
}

func (u *UserSession) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSessionSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence.CurrentSequence), nil
}

func (u *UserSession) Reduce(event *models.Event) (err error) {
	var session *view_model.UserSessionView
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.UserPasswordCheckFailed,
		es_model.MFAOTPCheckSucceeded,
		es_model.MFAOTPCheckFailed,
		es_model.SignedOut,
		es_model.HumanPasswordCheckSucceeded,
		es_model.HumanPasswordCheckFailed,
		es_model.HumanExternalLoginCheckSucceeded,
		es_model.HumanMFAOTPCheckSucceeded,
		es_model.HumanMFAOTPCheckFailed,
		es_model.HumanSignedOut:
		eventData, err := view_model.UserSessionFromEvent(event)
		if err != nil {
			return err
		}
		session, err = u.view.UserSessionByIDs(eventData.UserAgentID, event.AggregateID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			session = &view_model.UserSessionView{
				CreationDate:  event.CreationDate,
				ResourceOwner: event.ResourceOwner,
				UserAgentID:   eventData.UserAgentID,
				UserID:        event.AggregateID,
				State:         int32(req_model.UserSessionStateActive),
			}
		}
		return u.updateSession(session, event)
	case es_model.UserPasswordChanged,
		es_model.MFAOTPRemoved,
		es_model.UserProfileChanged,
		es_model.UserLocked,
		es_model.UserDeactivated,
		es_model.HumanPasswordChanged,
		es_model.HumanMFAOTPRemoved,
		es_model.HumanProfileChanged,
		es_model.DomainClaimed,
		es_model.UserUserNameChanged,
		es_model.HumanExternalIDPRemoved,
		es_model.HumanExternalIDPCascadeRemoved:
		sessions, err := u.view.UserSessionsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(sessions) == 0 {
			return u.view.ProcessedUserSessionSequence(event.Sequence, event.CreationDate)
		}
		for _, session := range sessions {
			if err := session.AppendEvent(event); err != nil {
				return err
			}
			if err := u.fillUserInfo(session, event.AggregateID); err != nil {
				return err
			}
		}
		return u.view.PutUserSessions(sessions, event.Sequence, event.CreationDate)
	case es_model.UserRemoved:
		return u.view.DeleteUserSessions(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedUserSessionSequence(event.Sequence, event.CreationDate)
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
	return u.view.PutUserSession(session, event.CreationDate)
}

func (u *UserSession) fillUserInfo(session *view_model.UserSessionView, id string) error {
	user, err := u.view.UserByID(id)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	return nil
}
