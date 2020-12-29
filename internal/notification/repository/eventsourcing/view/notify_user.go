package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	notifyUserTable = "notification.notify_users"
)

func (v *View) NotifyUserByID(userID string) (*model.NotifyUser, error) {
	return view.NotifyUserByID(v.Db, notifyUserTable, userID)
}

func (v *View) PutNotifyUser(user *model.NotifyUser, event *models.Event) error {
	err := view.PutNotifyUser(v.Db, notifyUserTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedNotifyUserSequence(event)
}

func (v *View) NotifyUsersByOrgID(orgID string) ([]*model.NotifyUser, error) {
	return view.NotifyUsersByOrgID(v.Db, notifyUserTable, orgID)
}

func (v *View) DeleteNotifyUser(userID string, event *models.Event) error {
	err := view.DeleteNotifyUser(v.Db, notifyUserTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedNotifyUserSequence(event)
}

func (v *View) GetLatestNotifyUserSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(notifyUserTable, aggregateType)
}

func (v *View) ProcessedNotifyUserSequence(event *models.Event) error {
	return v.saveCurrentSequence(notifyUserTable, event)
}

func (v *View) UpdateNotifyUserSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(notifyUserTable)
}

func (v *View) GetLatestNotifyUserFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notifyUserTable, sequence)
}

func (v *View) ProcessedNotifyUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
