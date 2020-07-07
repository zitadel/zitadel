package view

import (
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

func (v *View) PutNotifyUser(user *model.NotifyUser, sequence uint64) error {
	err := view.PutNotifyUser(v.Db, notifyUserTable, user)
	if err != nil {
		return err
	}
	if sequence != 0 {
		return v.ProcessedNotifyUserSequence(sequence)
	}
	return nil
}

func (v *View) NotifyUsersByOrgID(orgID string) ([]*model.NotifyUser, error) {
	return view.NotifyUsersByOrgID(v.Db, notifyUserTable, orgID)
}

func (v *View) DeleteNotifyUser(userID string, eventSequence uint64) error {
	err := view.DeleteNotifyUser(v.Db, notifyUserTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedNotifyUserSequence(eventSequence)
}

func (v *View) GetLatestNotifyUserSequence() (uint64, error) {
	return v.latestSequence(notifyUserTable)
}

func (v *View) ProcessedNotifyUserSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(notifyUserTable, eventSequence)
}

func (v *View) GetLatestNotifyUserFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(notifyUserTable, sequence)
}

func (v *View) ProcessedNotifyUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
