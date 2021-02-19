package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type User struct {
	es_models.ObjectRoot
	State    UserState
	UserName string

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
	ModifierID   string               `json:"modifierUser,omitempty"`
	ModifierName string               `json:"-"`
	Data         interface{}          `json:"data,omitempty"`
}
