package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

const (
	UserVersion = "v1"
)

type User struct {
	es_models.ObjectRoot
	State int32 `json:"-"`
	*Password
	*Profile
	*Email
	*Phone
	*Address
	InitCode *InitUserCode
}

type InitUserCode struct {
	es_models.ObjectRoot
	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func UserFromEvents(user *User, events ...*es_models.Event) (*User, error) {
	if user == nil {
		user = &User{}
	}

	return user, user.AppendEvents(events...)
}

func UserFromModel(user *model.User) *User {
	converted := &User{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  user.ObjectRoot.AggregateID,
			Sequence:     user.Sequence,
			ChangeDate:   user.ChangeDate,
			CreationDate: user.CreationDate,
		},
		State: int32(user.State),
	}
	if user.Password != nil {
		converted.Password = PasswordFromModel(user.Password)
	}
	if user.Profile != nil {
		converted.Profile = ProfileFromModel(user.Profile)
	}
	if user.Email != nil {
		converted.Email = EmailFromModel(user.Email)
	}
	if user.Phone != nil {
		converted.Phone = PhoneFromModel(user.Phone)
	}
	if user.Address != nil {
		converted.Address = AddressFromModel(user.Address)
	}
	return converted
}

func UserToModel(user *User) *model.User {
	converted := &model.User{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  user.ObjectRoot.AggregateID,
			Sequence:     user.Sequence,
			ChangeDate:   user.ChangeDate,
			CreationDate: user.CreationDate,
		},
		State: model.UserState(user.State),
	}
	if user.Password != nil {
		converted.Password = PasswordToModel(user.Password)
	}
	if user.Profile != nil {
		converted.Profile = ProfileToModel(user.Profile)
	}
	if user.Email != nil {
		converted.Email = EmailToModel(user.Email)
	}
	if user.Phone != nil {
		converted.Phone = PhoneToModel(user.Phone)
	}
	if user.Address != nil {
		converted.Address = AddressToModel(user.Address)
	}
	if user.InitCode != nil {
		converted.InitCode = InitCodeToModel(user.InitCode)
	}
	return converted
}

func (p *User) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) AppendEvent(event *es_models.Event) error {
	u.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.UserAdded, model.UserRegistered:
		if err := json.Unmarshal(event.Data, u); err != nil {
			logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
			return err
		}
	case model.UserDeactivated:
		u.appendDeactivatedEvent()
	case model.UserReactivated:
		u.appendReactivatedEvent()
	case model.UserLocked:
		u.appendLockedEvent()
	case model.UserUnlocked:
		u.appendUnlockedEvent()
	case model.InitializedUserCodeCreated:
		u.appendInitUsercodeCreatedEvent(event)
	}
	u.ComputeState()
	return nil
}

func (u *User) ComputeState() {
	if u.State == 0 {
		if u.Email != nil && u.IsEmailVerified {
			u.State = int32(model.USERSTATE_ACTIVE)
		} else {
			u.State = int32(model.USERSTATE_INITIAL)
		}
	}
}

func (u *User) appendDeactivatedEvent() error {
	u.State = int32(model.USERSTATE_INACTIVE)
	return nil
}

func (u *User) appendReactivatedEvent() error {
	u.State = int32(model.USERSTATE_ACTIVE)
	return nil
}

func (u *User) appendLockedEvent() error {
	u.State = int32(model.USERSTATE_LOCKED)
	return nil
}

func (u *User) appendUnlockedEvent() error {
	u.State = int32(model.USERSTATE_ACTIVE)
	return nil
}

func (u *User) appendInitUsercodeCreatedEvent(event *es_models.Event) error {
	initCode := new(InitUserCode)
	err := initCode.setData(event)
	if err != nil {
		return err
	}
	initCode.ObjectRoot.CreationDate = event.CreationDate
	u.InitCode = initCode
	return nil
}

func (c *InitUserCode) setData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-7duwe").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
