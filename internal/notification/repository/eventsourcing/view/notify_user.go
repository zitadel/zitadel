package view

import (
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	notifyUserTable = "notification.notify_users"
)

func (v *View) NotifyUserByID(userID string) (*model.NotifyUser, error) {
	return view.NotifyUserByID(v.Db, notifyUserTable, userID)
}

func (v *View) PutNotifyUser(user *model.NotifyUser, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutNotifyUser(v.Db, notifyUserTable, user)
	if err != nil {
		return err
	}
	if sequence != 0 {
		return v.ProcessedNotifyUserSequence(sequence, eventTimestamp)
	}
	return nil
}

func (v *View) NotifyUsersByOrgID(orgID string) ([]*model.NotifyUser, error) {
	return view.NotifyUsersByOrgID(v.Db, notifyUserTable, orgID)
}

func (v *View) DeleteNotifyUser(userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteNotifyUser(v.Db, notifyUserTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedNotifyUserSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestNotifyUserSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(notifyUserTable)
}

func (v *View) ProcessedNotifyUserSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(notifyUserTable, eventSequence, eventTimestamp)
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
