package member

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type MemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *MemberRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *MemberRemovedEvent) Data() interface{} {
	return e
}
