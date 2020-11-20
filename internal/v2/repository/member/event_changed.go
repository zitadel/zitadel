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

	hasChanged bool
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current,
	changed *Aggregate,
) (*ChangedEvent, error) {

	change := &ChangedEvent{
		BaseEvent: *base,
	}

	if current.UserID != changed.UserID {
		change.UserID = changed.UserID
		change.hasChanged = true
	}

	sort.Strings(current.Roles)
	sort.Strings(changed.Roles)
	if !reflect.DeepEqual(current.Roles, changed.Roles) {
		change.Roles = changed.Roles
		change.hasChanged = true
	}

	if !change.hasChanged {
		return nil, errors.ThrowPreconditionFailed(nil, "MEMBE-SeKlD", "Errors.NoChanges")
	}

	return change, nil

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
