package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type UserChanges struct {
	Changes      []*UserChange
	LastSequence uint64
}

type UserChange struct {
	ChangeDate        *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType         string               `json:"eventType,omitempty"`
	Sequence          uint64               `json:"sequence,omitempty"`
	ModifierID        string               `json:"modifierUser,omitempty"`
	ModifierName      string               `json:"-"`
	ModifierLoginName string               `json:"-"`
	Data              interface{}          `json:"data,omitempty"`
}
