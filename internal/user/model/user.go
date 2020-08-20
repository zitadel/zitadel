package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type User struct {
	es_models.ObjectRoot
	State UserState

	*Human
	*Machine
}

type UserState int32

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateDeleted
	UserStateLocked
	UserStateSuspend
	UserStateInitial
)

type UserChanges struct {
	Changes      []*UserChange
	LastSequence uint64
}

type UserChange struct {
	ChangeDate   *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType    string               `json:"eventType,omitempty"`
	Sequence     uint64               `json:"sequence,omitempty"`
	ModifierId   string               `json:"modifierUser,omitempty"`
	ModifierName string               `json:"-"`
	Data         interface{}          `json:"data,omitempty"`
}

func (u *User) IsActive() bool {
	return u.State == UserStateActive
}

func (u *User) IsInitial() bool {
	return u.State == UserStateInitial
}

func (u *User) IsInactive() bool {
	return u.State == UserStateInactive
}

func (u *User) IsLocked() bool {
	return u.State == UserStateLocked
}

func (u *User) IsValid() bool {
	if u.Human == nil && u.Machine == nil {
		return false
	}
	if u.Human != nil {
		return u.Human.IsValid()
	}
	return u.Machine.IsValid()
}
