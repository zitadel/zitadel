package model

import (
	"encoding/json"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

const (
	UserVersion = "v2"
)

type User struct {
	es_models.ObjectRoot
	State    int32  `json:"-"`
	UserName string `json:"userName"`

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
	u.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case UserAdded,
		HumanAdded,
		MachineAdded,
		UserRegistered,
		HumanRegistered,
		UserProfileChanged,
		DomainClaimed,
		UserUserNameChanged:
		err := u.setData(event)
		if err != nil {
			return err
		}
	case UserDeactivated:
		u.appendDeactivatedEvent()
	case UserReactivated:
		u.appendReactivatedEvent()
	case UserLocked:
		u.appendLockedEvent()
	case UserUnlocked:
		u.appendUnlockedEvent()
	case UserRemoved:
		u.appendRemovedEvent()
	}

	if u.Human != nil {
		u.Human.user = u
		return u.Human.AppendEvent(event)
	} else if u.Machine != nil {
		u.Machine.user = u
		return u.Machine.AppendEvent(event)
	}
	if strings.HasPrefix(string(event.Type), "user.human") || event.AggregateVersion == "v1" {
		u.Human = &Human{user: u}
		return u.Human.AppendEvent(event)
	}
	if strings.HasPrefix(string(event.Type), "user.machine") {
		u.Machine = &Machine{user: u}
		return u.Machine.AppendEvent(event)
	}

	return errors.ThrowNotFound(nil, "MODEL-x9TaX", "Errors.UserType.Undefined")
}

func (u *User) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-ZDzQy").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-yGmhh", "could not unmarshal event")
	}
	return nil
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

func (u *User) appendRemovedEvent() {
	u.State = int32(model.UserStateDeleted)
}
