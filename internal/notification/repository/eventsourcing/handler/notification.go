package handler

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/types"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"net/http"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Notification struct {
	handler
	eventstore     eventstore.Eventstore
	userEvents     *usr_event.UserEventstore
	systemDefaults sd.SystemDefaults
	AesCrypto      crypto.EncryptionAlgorithm
	i18n           *i18n.Translator
	statikDir      http.FileSystem
}

const (
	notificationTable = "notification.notifications"
	NotifyUserID      = "NOTIFICATION"
)

func (n *Notification) MinimumCycleDuration() time.Duration { return n.cycleDuration }

func (n *Notification) ViewModel() string {
	return notificationTable
}

func (n *Notification) EventQuery() (*models.SearchQuery, error) {
	sequence, err := n.view.GetLatestNotificationSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (n *Notification) Reduce(event *models.Event) (err error) {
	switch event.Type {
	case es_model.InitializedUserCodeAdded:
		err = n.handleInitUserCode(event)
	case es_model.UserEmailCodeAdded:
		err = n.handleEmailVerificationCode(event)
	case es_model.UserPhoneCodeAdded:
		err = n.handlePhoneVerificationCode(event)
	case es_model.UserPasswordCodeAdded:
		err = n.handlePasswordCode(event)
	default:
		return n.view.ProcessedNotificationSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return n.view.ProcessedNotificationSequence(event.Sequence)
}

func (n *Notification) handleInitUserCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.InitializedUserCodeAdded, es_model.InitializedUserCodeSent)
	if err != nil || alreadyHandled {
		return err
	}
	initCode := new(es_model.InitUserCode)
	initCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendUserInitCode(n.statikDir, n.i18n, user, initCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.InitCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handlePasswordCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserPasswordCodeAdded, es_model.UserPasswordCodeSent)
	if err != nil || alreadyHandled {
		return err
	}
	pwCode := new(es_model.PasswordCode)
	pwCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPasswordCode(n.statikDir, n.i18n, user, pwCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.PasswordCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handleEmailVerificationCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserEmailCodeAdded, es_model.UserEmailCodeSent)
	if err != nil || alreadyHandled {
		return nil
	}
	emailCode := new(es_model.EmailCode)
	emailCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendEmailVerificationCode(n.statikDir, n.i18n, user, emailCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.EmailVerificationCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handlePhoneVerificationCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserPhoneCodeAdded, es_model.UserPhoneCodeSent)
	if err != nil || alreadyHandled {
		return nil
	}
	phoneCode := new(es_model.PhoneCode)
	phoneCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPhoneVerificationCode(n.i18n, user, phoneCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.PhoneVerificationCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) checkIfCodeAlreadyHandled(userID string, sequence uint64, addedType, sentType models.EventType) (bool, error) {
	events, err := n.getUserEvents(userID, sequence)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		if event.Type == addedType || event.Type == sentType {
			return true, nil
		}
	}
	return false, nil
}

func (n *Notification) getUserEvents(userID string, sequence uint64) ([]*models.Event, error) {
	query, err := eventsourcing.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return n.eventstore.FilterEvents(context.Background(), query)
}

func (n *Notification) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-s9opc", "id", event.AggregateID, "sequence", event.Sequence).WithError(err).Warn("something went wrong in notification handler")
	return spooler.HandleError(event, err, n.view.GetLatestNotificationFailedEvent, n.view.ProcessedNotificationFailedEvent, n.view.ProcessedNotificationSequence, n.errorCountUntilSkip)
}

func getSetNotifyContextData(orgID string) context.Context {
	return auth.SetCtxData(context.Background(), auth.CtxData{UserID: NotifyUserID, OrgID: orgID})
}
