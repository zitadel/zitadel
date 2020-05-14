package handler

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/types"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Notification struct {
	handler
	userEvents     *usr_event.UserEventstore
	systemDefaults sd.SystemDefaults
}

const (
	notificationTable = "notification.notifications"
	NOTIFY_USER       = "NOTIFICATION"
)

func (p *Notification) MinimumCycleDuration() time.Duration { return p.cycleDuration }

func (p *Notification) ViewModel() string {
	return notificationTable
}

func (p *Notification) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestNotificationSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (p *Notification) Process(event *models.Event) (err error) {
	switch event.Type {
	case es_model.InitializedUserCodeAdded:
		err := p.handleInitUserCode(event)
		if err != nil {
			return err
		}
		err = p.userEvents.InitCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
		if err != nil {
			return err
		}
		return p.view.ProcessedNotificationSequence(event.Sequence)
	case es_model.UserEmailCodeAdded,
		es_model.UserPhoneCodeAdded,
		es_model.UserPasswordCodeAdded:

	default:
		return p.view.ProcessedNotificationSequence(event.Sequence)
	}
	return nil
}

func (p *Notification) handleInitUserCode(event *models.Event) (err error) {
	initCode := new(es_model.InitUserCode)
	initCode.SetData(event)
	user, err := p.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendUserInitCode(user, initCode, p.systemDefaults)
	return err
}

func (p *Notification) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-s9opc", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, p.view.GetLatestNotificationFailedEvent, p.view.ProcessedNotificationFailedEvent, p.view.ProcessedNotificationSequence, p.errorCountUntilSkip)
}

func getSetNotifyContextData(orgID string) context.Context {
	return context.WithValue(context.Background(), auth.GetCtxDataKey(), auth.CtxData{UserID: NOTIFY_USER, OrgID: orgID})
}
