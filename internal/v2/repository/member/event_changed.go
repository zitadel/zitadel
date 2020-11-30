package member

import (
	"encoding/json"
	"reflect"
	"sort"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	ChangedEventType = "member.changed"
)

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles,omitempty"`
	UserID string   `json:"userId,omitempty"`
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func ChangeEventFromExisting(
	base *eventstore.BaseEvent,
	current *WriteModel,
	roles ...string,
) (*ChangedEvent, error) {

	change := NewChangedEvent(base, current.UserID)
	hasChanged := false

	sort.Strings(current.Roles)
	sort.Strings(roles)
	if !reflect.DeepEqual(current.Roles, roles) {
		change.Roles = roles
		hasChanged = true
	}

	if !hasChanged {
		return nil, errors.ThrowPreconditionFailed(nil, "MEMBE-SeKlD", "Errors.NoChanges")
	}

	return change, nil
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	userID string,
	roles ...string,
) *ChangedEvent {

	return &ChangedEvent{
		BaseEvent: *base,
		UserID:    userID,
		Roles:     roles,
	}
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}
