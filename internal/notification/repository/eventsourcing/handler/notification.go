package handler

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/api/auth"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
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
	AesCrypto      crypto.EncryptionAlgorithm
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
		err = p.handleInitUserCode(event)
	case es_model.UserEmailCodeAdded:
		err = p.handleEmailVerificationCode(event)
	case es_model.UserPhoneCodeAdded:
		err = p.handlePhoneVerificationCode(event)
	case es_model.UserPasswordCodeAdded:
		err = p.handlePasswordCode(event)
	default:
		return p.view.ProcessedNotificationSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return p.view.ProcessedNotificationSequence(event.Sequence)
}

func (p *Notification) handleInitUserCode(event *models.Event) (err error) {
	initCode := new(es_model.InitUserCode)
	initCode.SetData(event)
	user, err := p.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendUserInitCode(user, initCode, p.systemDefaults, p.AesCrypto)
	if err != nil {
		return err
	}
	return p.userEvents.InitCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (p *Notification) handlePasswordCode(event *models.Event) (err error) {
	pwCode := new(es_model.PasswordCode)
	pwCode.SetData(event)
	user, err := p.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPasswordCodeCode(user, pwCode, p.systemDefaults, p.AesCrypto)
	if err != nil {
		return err
	}
	return p.userEvents.PasswordCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (p *Notification) handleEmailVerificationCode(event *models.Event) (err error) {
	emailCode := new(es_model.EmailCode)
	emailCode.SetData(event)
	user, err := p.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendEmailVerificationCode(user, emailCode, p.systemDefaults, p.AesCrypto)
	if err != nil {
		return err
	}
	return p.userEvents.EmailVerificationCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (p *Notification) handlePhoneVerificationCode(event *models.Event) (err error) {
	phoneCode := new(es_model.PhoneCode)
	phoneCode.SetData(event)
	user, err := p.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPhoneVerificationCode(user, phoneCode, p.systemDefaults, p.AesCrypto)
	if err != nil {
		return err
	}
	return p.userEvents.PhoneVerificationCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (p *Notification) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-s9opc", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, p.view.GetLatestNotificationFailedEvent, p.view.ProcessedNotificationFailedEvent, p.view.ProcessedNotificationSequence, p.errorCountUntilSkip)
}

func getSetNotifyContextData(orgID string) context.Context {
	return context.WithValue(context.Background(), auth.GetCtxDataKey(), auth.CtxData{UserID: NOTIFY_USER, OrgID: orgID})
}
