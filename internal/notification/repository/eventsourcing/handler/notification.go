package handler

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/api/auth"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
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
	eventstore     eventstore.Eventstore
	userEvents     *usr_event.UserEventstore
	systemDefaults sd.SystemDefaults
	AesCrypto      crypto.EncryptionAlgorithm
}

const (
	notificationTable = "notification.notifications"
	NOTIFY_USER       = "NOTIFICATION"
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

func (n *Notification) Process(event *models.Event) (err error) {
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
	if err != nil {
		return err
	}
	if alreadyHandled {
		return nil
	}
	initCode := new(es_model.InitUserCode)
	initCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendUserInitCode(user, initCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.InitCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handlePasswordCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserPasswordCodeAdded, es_model.UserPasswordCodeSent)
	if err != nil {
		return err
	}
	if alreadyHandled {
		return nil
	}
	pwCode := new(es_model.PasswordCode)
	pwCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPasswordCodeCode(user, pwCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.PasswordCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handleEmailVerificationCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserEmailCodeAdded, es_model.UserEmailCodeSent)
	if err != nil {
		return err
	}
	if alreadyHandled {
		return nil
	}
	emailCode := new(es_model.EmailCode)
	emailCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendEmailVerificationCode(user, emailCode, n.systemDefaults, n.AesCrypto)
	if err != nil {
		return err
	}
	return n.userEvents.EmailVerificationCodeSent(getSetNotifyContextData(event.ResourceOwner), event.AggregateID)
}

func (n *Notification) handlePhoneVerificationCode(event *models.Event) (err error) {
	alreadyHandled, err := n.checkIfCodeAlreadyHandled(event.AggregateID, event.Sequence, es_model.UserPhoneCodeAdded, es_model.UserPhoneCodeSent)
	if err != nil {
		return err
	}
	if alreadyHandled {
		return nil
	}
	phoneCode := new(es_model.PhoneCode)
	phoneCode.SetData(event)
	user, err := n.view.NotifyUserByID(event.AggregateID)
	if err != nil {
		return err
	}
	err = types.SendPhoneVerificationCode(user, phoneCode, n.systemDefaults, n.AesCrypto)
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
	fmt.Printf("ONError ", event)
	return spooler.HandleError(event, err, n.view.GetLatestNotificationFailedEvent, n.view.ProcessedNotificationFailedEvent, n.view.ProcessedNotificationSequence, n.errorCountUntilSkip)
}

func getSetNotifyContextData(orgID string) context.Context {
	return context.WithValue(context.Background(), auth.GetCtxDataKey(), auth.CtxData{UserID: NOTIFY_USER, OrgID: orgID})
}
