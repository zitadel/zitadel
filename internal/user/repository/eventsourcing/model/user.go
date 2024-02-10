package model

import (
	"encoding/json"
	"strings"

	"github.com/zitadel/logging"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
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

	switch event.Type() {
	case user.UserV1AddedType,
		user.HumanAddedType,
		user.MachineAddedEventType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.UserV1ProfileChangedType,
		user.UserDomainClaimedType,
		user.UserUserNameChangedType:
		err := u.setData(event)
		if err != nil {
			return err
		}
	case user.UserDeactivatedType:
		u.appendDeactivatedEvent()
	case user.UserReactivatedType:
		u.appendReactivatedEvent()
	case user.UserLockedType:
		u.appendLockedEvent()
	case user.UserUnlockedType:
		u.appendUnlockedEvent()
	case user.UserRemovedType:
		u.appendRemovedEvent()
	}

	if u.Human != nil {
		u.Human.user = u
		return u.Human.AppendEvent(event)
	} else if u.Machine != nil {
		u.Machine.user = u
		return u.Machine.AppendEvent(event)
	}
	if strings.HasPrefix(string(event.Typ), "user.human") || event.AggregateVersion == "v1" {
		u.Human = &Human{user: u}
		return u.Human.AppendEvent(event)
	}
	if strings.HasPrefix(string(event.Typ), "user.machine") {
		u.Machine = &Machine{user: u}
		return u.Machine.AppendEvent(event)
	}

	return zerrors.ThrowNotFound(nil, "MODEL-x9TaX", "Errors.UserType.Undefined")
}

func (u *User) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-ZDzQy").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-yGmhh", "could not unmarshal event")
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
