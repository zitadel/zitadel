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
	InitCode     *InitUserCode
	EmailCode    *EmailCode
	PhoneCode    *PhoneCode
	PasswordCode *RequestPasswordSet
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
	if user.EmailCode != nil {
		converted.EmailCode = EmailCodeToModel(user.EmailCode)
	}
	if user.PhoneCode != nil {
		converted.PhoneCode = PhoneCodeToModel(user.PhoneCode)
	}
	if user.PasswordCode != nil {
		converted.PasswordCode = PasswordCodeToModel(user.PasswordCode)
	}
	return converted
}

func InitCodeToModel(code *InitUserCode) *model.InitUserCode {
	return &model.InitUserCode{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  code.ObjectRoot.AggregateID,
			Sequence:     code.Sequence,
			ChangeDate:   code.ChangeDate,
			CreationDate: code.CreationDate,
		},
		Expiry: code.Expiry,
		Code:   code.Code,
	}
}

func EmailCodeToModel(code *EmailCode) *model.EmailCode {
	return &model.EmailCode{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  code.ObjectRoot.AggregateID,
			Sequence:     code.Sequence,
			ChangeDate:   code.ChangeDate,
			CreationDate: code.CreationDate,
		},
		Expiry: code.Expiry,
		Code:   code.Code,
	}
}

func PhoneCodeToModel(code *PhoneCode) *model.PhoneCode {
	return &model.PhoneCode{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  code.ObjectRoot.AggregateID,
			Sequence:     code.Sequence,
			ChangeDate:   code.ChangeDate,
			CreationDate: code.CreationDate,
		},
		Expiry: code.Expiry,
		Code:   code.Code,
	}
}
func PasswordCodeToModel(code *RequestPasswordSet) *model.RequestPasswordSet {
	return &model.RequestPasswordSet{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  code.ObjectRoot.AggregateID,
			Sequence:     code.Sequence,
			ChangeDate:   code.ChangeDate,
			CreationDate: code.CreationDate,
		},
		Expiry:           code.Expiry,
		Code:             code.Code,
		NotificationType: model.NotificationType(code.NotificationType),
	}
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
	var err error
	switch event.Type {
	case model.UserAdded,
		model.UserRegistered,
		model.UserProfileChanged:
		if err := json.Unmarshal(event.Data, u); err != nil {
			logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
			return err
		}
	case model.UserDeactivated:
		err = u.appendDeactivatedEvent()
	case model.UserReactivated:
		err = u.appendReactivatedEvent()
	case model.UserLocked:
		err = u.appendLockedEvent()
	case model.UserUnlocked:
		err = u.appendUnlockedEvent()
	case model.InitializedUserCodeCreated:
		err = u.appendInitUsercodeCreatedEvent(event)
	case model.UserPasswordChanged:
		err = u.appendUserPasswordChangedEvent(event)
	case model.UserPasswordSetRequested:
		err = u.appendPasswordSetRequestedEvent(event)
	case model.UserEmailChanged:
		err = u.appendUserEmailChangedEvent(event)
	case model.UserEmailCodeAdded:
		err = u.appendUserEmailCodeAddedEvent(event)
	case model.UserEmailVerified:
		err = u.appendUserEmailVerifiedEvent()
	case model.UserPhoneChanged:
		err = u.appendUserPhoneChangedEvent(event)
	case model.UserPhoneCodeAdded:
		err = u.appendUserPhoneCodeAddedEvent(event)
	case model.UserPhoneVerified:
		err = u.appendUserPhoneVerifiedEvent()
	case model.UserAddressChanged:
		err = u.appendUserAddressChangedEvent(event)
	}
	if err != nil {
		return err
	}
	u.ComputeObject()
	return nil
}

func (u *User) ComputeObject() {
	if u.State == 0 {
		if u.Email != nil && u.IsEmailVerified {
			u.State = int32(model.USERSTATE_ACTIVE)
		} else {
			u.State = int32(model.USERSTATE_INITIAL)
		}
	}
	if u.Password != nil && u.Password.ObjectRoot.AggregateID == "" {
		u.Password.ObjectRoot = u.ObjectRoot
	}
	if u.Profile != nil && u.Profile.ObjectRoot.AggregateID == "" {
		u.Profile.ObjectRoot = u.ObjectRoot
	}
	if u.Email != nil && u.Email.ObjectRoot.AggregateID == "" {
		u.Email.ObjectRoot = u.ObjectRoot
	}
	if u.Phone != nil && u.Phone.ObjectRoot.AggregateID == "" {
		u.Phone.ObjectRoot = u.ObjectRoot
	}
	if u.Address != nil && u.Address.ObjectRoot.AggregateID == "" {
		u.Address.ObjectRoot = u.ObjectRoot
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
