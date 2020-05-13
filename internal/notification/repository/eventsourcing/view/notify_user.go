package view

import (
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	notifyUserTable = "notification.notify_users"
)

func (v *View) NotifyUserByID(userID string) (*model.NotifyUser, error) {
	return view.NotifyUserByID(v.Db, notifyUserTable, userID)
}

func (v *View) PutNotifyUser(user *model.NotifyUser) error {
	err := view.PutNotifyUser(v.Db, notifyUserTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(user.Sequence)
}

func (v *View) DeleteNotifyUser(userID string, eventSequence uint64) error {
	err := view.DeleteNotifyUser(v.Db, notifyUserTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserSequence(eventSequence)
}

func (v *View) GetLatestUserSequence() (uint64, error) {
	return v.latestSequence(notifyUserTable)
}

func (v *View) ProcessedUserSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(notifyUserTable, eventSequence)
}

func (v *View) GetLatestUserFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(notifyUserTable, sequence)
}

func (v *View) ProcessedUserFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
