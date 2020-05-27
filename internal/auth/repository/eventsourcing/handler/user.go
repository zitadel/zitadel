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

type User struct {
	handler
	eventstore eventstore.Eventstore
}

const (
	userTable = "auth.users"
)

func (p *User) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *User) ViewModel() string {
	return userTable
}

func (p *User) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestUserSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (p *User) Process(event *models.Event) (err error) {
	user := new(view_model.UserView)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		user.AppendEvent(event)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserEmailVerified,
		es_model.UserPhoneChanged,
		es_model.UserPhoneVerified,
		es_model.UserAddressChanged,
		es_model.UserDeactivated,
		es_model.UserReactivated,
		es_model.UserLocked,
		es_model.UserUnlocked,
		es_model.MfaOtpAdded,
		es_model.MfaOtpVerified,
		es_model.MfaOtpRemoved:
		user, err = p.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.UserRemoved:
		err = p.view.DeleteUser(event.AggregateID, event.Sequence)
	default:
		return p.view.ProcessedUserSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.PutUser(user)
}

func (p *User) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-is8wa", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, p.view.GetLatestUserFailedEvent, p.view.ProcessedUserFailedEvent, p.view.ProcessedUserSequence, p.errorCountUntilSkip)
}
