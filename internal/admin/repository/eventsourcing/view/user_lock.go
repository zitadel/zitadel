package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userLockTable = "adminapi.user_lock"
)

func (v *View) UserLockByID(userID string) (*model.UserLockView, error) {
	return view.UserLockByID(v.Db, userLockTable, userID)
}

func (v *View) PutUserLock(user *model.UserLockView, event *models.Event) error {
	err := view.PutUserLock(v.Db, userLockTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserLockSequence(event)
}

func (v *View) DeleteUserLock(userID string, event *models.Event) error {
	err := view.DeleteUserLock(v.Db, userLockTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserLockSequence(event)
}

func (v *View) GetLatestUserLockSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(userLockTable)
}

func (v *View) ProcessedUserLockSequence(event *models.Event) error {
	return v.saveCurrentSequence(userLockTable, event)
}

func (v *View) UpdateUserLockSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userLockTable)
}

func (v *View) GetLatestUserLockFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userLockTable, sequence)
}

func (v *View) ProcessedUserLockFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
