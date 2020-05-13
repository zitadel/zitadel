package handler

import (
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type NotifyUser struct {
	handler
	eventstore eventstore.Eventstore
}

const (
	userTable = "notification.notify_users"
)

func (p *NotifyUser) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *NotifyUser) ViewModel() string {
	return userTable
}

func (p *NotifyUser) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestNotifyUserSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (p *NotifyUser) Process(event *models.Event) (err error) {
	user := new(view_model.NotifyUser)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		user.AppendEvent(event)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserEmailVerified,
		es_model.UserPhoneChanged,
		es_model.UserPhoneVerified:
		user, err = p.view.NotifyUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.UserDeleted:
		err = p.view.DeleteNotifyUser(event.AggregateID, event.Sequence)
	default:
		return p.view.ProcessedNotifyUserSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutNotifyUser(user)
}

func (p *NotifyUser) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-s9opc", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, p.view.GetLatestNotifyUserFailedEvent, p.view.ProcessedNotifyUserFailedEvent, p.view.ProcessedNotifyUserSequence, p.errorCountUntilSkip)
}
