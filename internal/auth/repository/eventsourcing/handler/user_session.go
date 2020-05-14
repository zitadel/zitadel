package handler

import (
	"context"
	"time"

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

func (u *UserSession) MinimumCycleDuration() time.Duration { return u.cycleDuration }

func (u *UserSession) ViewModel() string {
	return userSessionTable
}

func (u *UserSession) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSessionSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (u *UserSession) Process(event *models.Event) (err error) {
	eventData, err := view_model.UserSessionFromEvent(event)
	if err != nil {
		return err
	}
	session, err := u.view.UserSessionByIDs(eventData.UserAgentID, event.AggregateID)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		session = &view_model.UserSessionView{
			CreationDate:  event.CreationDate,
			ResourceOwner: event.ResourceOwner,
			UserAgentID:   eventData.UserAgentID,
			UserID:        event.AggregateID,
		}
	}
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.UserPasswordCheckFailed,
		es_model.UserPasswordChanged,
		es_model.MfaOtpCheckSucceeded,
		es_model.MfaOtpCheckFailed,
		es_model.MfaOtpRemoved:
		err = session.AppendEvent(event)
	default:
		return u.view.ProcessedUserSessionSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	if err := u.FillUserInfo(session, event.AggregateID); err != nil {
		return err
	}
	return u.view.PutUserSession(session)
}

func (u *UserSession) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-sdfw3s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user session handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *UserSession) FillUserInfo(session *view_model.UserSessionView, id string) error {
	user, err := u.userEvents.UserByID(context.Background(), id)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	return nil
}
