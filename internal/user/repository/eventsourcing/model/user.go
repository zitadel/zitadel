package model

import (
	"strings"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

const (
	UserVersion = "v2"
)

type User struct {
	es_models.ObjectRoot
	State int32

	*Human
	*Machine
}

func (u *User) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := u.AppendEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (u *User) AppendEvent(event *es_models.Event) error {
	switch event.Type {
	case UserDeactivated:
		u.appendDeactivatedEvent()
	case UserReactivated:
		u.appendReactivatedEvent()
	case UserLocked:
		u.appendLockedEvent()
	case UserUnlocked:
		u.appendUnlockedEvent()
	}

	if u.Human != nil {
		u.Human.objectRoot = u.ObjectRoot
		u.Human.state = u.State
		err := u.Human.AppendEvent(event)
		u.State = u.Human.state
		return err
	} else if u.Machine != nil {
		u.Machine.objectRoot = u.ObjectRoot
		u.Machine.state = u.State
		err := u.Machine.AppendEvent(event)
		u.State = u.Machine.state
		return err
	}
	if strings.HasPrefix(string(event.Type), "user.human") || event.AggregateVersion == "v1" {
		u.Human = &Human{
			objectRoot: u.ObjectRoot,
			state:      u.State,
		}
		err := u.Human.AppendEvent(event)
		u.State = u.Human.state
		return err
	}
	if strings.HasPrefix(string(event.Type), "user.machine") {
		u.Machine = &Machine{
			objectRoot: u.ObjectRoot,
			state:      u.State,
		}
		err := u.Machine.AppendEvent(event)
		u.State = u.Machine.state
		return err
	}

	return errors.ThrowNotFound(nil, "MODEL-x9TaX", "Errors.UserType.Undefined")
}

func (u *User) appendDeactivatedEvent() {
	u.State = int32(model.UserStateInactive)
}

func (u *User) appendReactivatedEvent() {
	u.State = int32(model.UserStateActive)
}

func (u *User) appendLockedEvent() {
	u.State = int32(model.UserStateLocked)
}

func (u *User) appendUnlockedEvent() {
	u.State = int32(model.UserStateActive)
}
