package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	notifyUserTable = "notification.notify_users"
)

func (v *View) NotifyUserByID(userID, instanceID string) (*model.NotifyUser, error) {
	return view.NotifyUserByID(v.Db, notifyUserTable, userID, instanceID)
}

func (v *View) PutNotifyUser(user *model.NotifyUser, event *models.Event) error {
	err := view.PutNotifyUser(v.Db, notifyUserTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedNotifyUserSequence(event)
}

func (v *View) NotifyUsersByOrgID(orgID, instanceID string) ([]*model.NotifyUser, error) {
	return view.NotifyUsersByOrgID(v.Db, notifyUserTable, orgID, instanceID)
}

func (v *View) DeleteNotifyUser(userID, instanceID string, event *models.Event) error {
	err := view.DeleteNotifyUser(v.Db, notifyUserTable, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedNotifyUserSequence(event)
}

func (v *View) GetLatestNotifyUserSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(notifyUserTable, instanceID)
}

func (v *View) GetLatestNotifyUserSequences() ([]*repository.CurrentSequence, error) {
	return v.latestSequences(notifyUserTable)
}

func (v *View) ProcessedNotifyUserSequence(event *models.Event) error {
	return v.saveCurrentSequence(notifyUserTable, event)
}

func (v *View) UpdateNotifyUserSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(notifyUserTable)
}

func (v *View) GetLatestNotifyUserFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notifyUserTable, instanceID, sequence)
}

func (v *View) ProcessedNotifyUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
