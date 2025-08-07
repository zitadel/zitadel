package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserChanges struct {
	Changes      []*UserChange
	LastSequence uint64
}

type UserChange struct {
	ChangeDate        *timestamppb.Timestamp `json:"changeDate,omitempty"`
	EventType         string                 `json:"eventType,omitempty"`
	Sequence          uint64                 `json:"sequence,omitempty"`
	ModifierID        string                 `json:"modifierUser,omitempty"`
	ModifierName      string                 `json:"-"`
	ModifierLoginName string                 `json:"-"`
	ModifierAvatarURL string                 `json:"-"`
	Data              interface{}            `json:"data,omitempty"`
}
