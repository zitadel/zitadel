package model

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

const (
	UserLockedKeyUserID = "user_id"
)

type UserLockView struct {
	UserID                   string    `json:"-" gorm:"column:user_id;primary_key"`
	ChangeDate               time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner            string    `json:"-" gorm:"column:resource_owner"`
	Sequence                 uint64    `json:"-" gorm:"column:sequence"`
	State                    int32     `json:"-" gorm:"column:user_state"`
	PasswordCheckFailedCount uint64    `json:"-" gorm:"column:password_check_failed_count"`
}

func (u *UserLockView) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.MachineAdded,
		es_model.HumanAdded:
		u.UserID = event.AggregateID
		u.ResourceOwner = event.ResourceOwner
		u.State = int32(model.UserStateActive)
	case es_model.UserPasswordCheckFailed:
		u.PasswordCheckFailedCount += 1
	case es_model.UserPasswordCheckSucceeded:
		u.PasswordCheckFailedCount = 0
	case es_model.UserLocked:
		u.State = int32(model.UserStateLocked)
	case es_model.UserUnlocked:
		u.PasswordCheckFailedCount = 0
		u.State = int32(model.UserStateActive)
	}
	return err
}
