package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
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
	InitCode     *InitUserCode `json:"-"`
	EmailCode    *EmailCode    `json:"-"`
	PhoneCode    *PhoneCode    `json:"-"`
	PasswordCode *PasswordCode `json:"-"`
	OTP          *OTP          `json:"-"`
}

type InitUserCode struct {
	es_models.ObjectRoot
	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func UserFromModel(user *model.User) *User {
	converted := &User{
		ObjectRoot: user.ObjectRoot,
		State:      int32(user.State),
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
	if user.OTP != nil {
		converted.OTP = OTPFromModel(user.OTP)
	}
	return converted
}

func UserToModel(user *User) *model.User {
	converted := &model.User{
		ObjectRoot: user.ObjectRoot,
		State:      model.UserState(user.State),
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
	if user.OTP != nil {
		converted.OTP = OTPToModel(user.OTP)
	}
	return converted
}

func InitCodeFromModel(code *model.InitUserCode) *InitUserCode {
	if code == nil {
		return nil
	}
	return &InitUserCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func InitCodeToModel(code *InitUserCode) *model.InitUserCode {
	return &model.InitUserCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
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

func (u *User) AppendEvent(event *es_models.Event) (err error) {
	u.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case UserAdded,
		UserRegistered,
		UserProfileChanged:
		u.setData(event)
	case UserDeactivated:
		u.appendDeactivatedEvent()
	case UserReactivated:
		u.appendReactivatedEvent()
	case UserLocked:
		u.appendLockedEvent()
	case UserUnlocked:
		u.appendUnlockedEvent()
	case InitializedUserCodeAdded:
		u.appendInitUsercodeCreatedEvent(event)
	case UserPasswordChanged:
		err = u.appendUserPasswordChangedEvent(event)
	case UserPasswordCodeAdded:
		err = u.appendPasswordSetRequestedEvent(event)
	case UserEmailChanged:
		err = u.appendUserEmailChangedEvent(event)
	case UserEmailCodeAdded:
		err = u.appendUserEmailCodeAddedEvent(event)
	case UserEmailVerified:
		u.appendUserEmailVerifiedEvent()
	case UserPhoneChanged:
		err = u.appendUserPhoneChangedEvent(event)
	case UserPhoneCodeAdded:
		err = u.appendUserPhoneCodeAddedEvent(event)
	case UserPhoneVerified:
		u.appendUserPhoneVerifiedEvent()
	case UserAddressChanged:
		err = u.appendUserAddressChangedEvent(event)
	case MfaOtpAdded:
		err = u.appendOtpAddedEvent(event)
	case MfaOtpVerified:
		u.appendOtpVerifiedEvent()
	case MfaOtpRemoved:
		u.appendOtpRemovedEvent()
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
	if u.Password != nil && u.Password.ObjectRoot.IsZero() {
		u.Password.ObjectRoot = u.ObjectRoot
	}
	if u.Profile != nil && u.Profile.ObjectRoot.IsZero() {
		u.Profile.ObjectRoot = u.ObjectRoot
	}
	if u.Email != nil && u.Email.ObjectRoot.IsZero() {
		u.Email.ObjectRoot = u.ObjectRoot
	}
	if u.Phone != nil && u.Phone.ObjectRoot.IsZero() {
		u.Phone.ObjectRoot = u.ObjectRoot
	}
	if u.Address != nil && u.Address.ObjectRoot.IsZero() {
		u.Address.ObjectRoot = u.ObjectRoot
	}
}

func (u *User) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sj4jd", "could not unmarshal event")
	}
	return nil
}

func (u *User) appendDeactivatedEvent() {
	u.State = int32(model.USERSTATE_INACTIVE)
}

func (u *User) appendReactivatedEvent() {
	u.State = int32(model.USERSTATE_ACTIVE)
}

func (u *User) appendLockedEvent() {
	u.State = int32(model.USERSTATE_LOCKED)
}

func (u *User) appendUnlockedEvent() {
	u.State = int32(model.USERSTATE_ACTIVE)
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
		return caos_errs.ThrowInternal(err, "MODEL-lo34s", "could not unmarshal event")
	}
	return nil
}
