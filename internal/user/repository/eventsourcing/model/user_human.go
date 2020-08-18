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

type Human struct {
	objectRoot es_models.ObjectRoot
	state      int32

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

func HumanFromModel(user *model.Human) *Human {
	human := new(Human)
	if user.Password != nil {
		human.Password = PasswordFromModel(user.Password)
	}
	if user.Profile != nil {
		human.Profile = ProfileFromModel(user.Profile)
	}
	if user.Email != nil {
		human.Email = EmailFromModel(user.Email)
	}
	if user.Phone != nil {
		human.Phone = PhoneFromModel(user.Phone)
	}
	if user.Address != nil {
		human.Address = AddressFromModel(user.Address)
	}
	if user.OTP != nil {
		human.OTP = OTPFromModel(user.OTP)
	}
	return human
}

func HumanToModel(user *Human) *model.Human {
	human := new(model.Human)
	if user.Password != nil {
		human.Password = PasswordToModel(user.Password)
	}
	if user.Profile != nil {
		human.Profile = ProfileToModel(user.Profile)
	}
	if user.Email != nil {
		human.Email = EmailToModel(user.Email)
	}
	if user.Phone != nil {
		human.Phone = PhoneToModel(user.Phone)
	}
	if user.Address != nil {
		human.Address = AddressToModel(user.Address)
	}
	if user.InitCode != nil {
		human.InitCode = InitCodeToModel(user.InitCode)
	}
	if user.EmailCode != nil {
		human.EmailCode = EmailCodeToModel(user.EmailCode)
	}
	if user.PhoneCode != nil {
		human.PhoneCode = PhoneCodeToModel(user.PhoneCode)
	}
	if user.PasswordCode != nil {
		human.PasswordCode = PasswordCodeToModel(user.PasswordCode)
	}
	if user.OTP != nil {
		human.OTP = OTPToModel(user.OTP)
	}
	return human
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

func (p *Human) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (h *Human) AppendEvent(event *es_models.Event) (err error) {
	switch event.Type {
	case UserAdded,
		UserRegistered,
		UserProfileChanged,
		HumanAdded,
		HumanRegistered,
		HumanProfileChanged:
		h.setData(event)
	case InitializedUserCodeAdded,
		InitializedHumanCodeAdded:
		h.appendInitUsercodeCreatedEvent(event)
	case UserPasswordChanged,
		HumanPasswordChanged:
		err = h.appendUserPasswordChangedEvent(event)
	case UserPasswordCodeAdded,
		HumanPasswordCodeAdded:
		err = h.appendPasswordSetRequestedEvent(event)
	case UserEmailChanged,
		HumanEmailChanged:
		err = h.appendUserEmailChangedEvent(event)
	case UserEmailCodeAdded,
		HumanEmailCodeAdded:
		err = h.appendUserEmailCodeAddedEvent(event)
	case UserEmailVerified,
		HumanEmailVerified:
		h.appendUserEmailVerifiedEvent()
	case UserPhoneChanged,
		HumanPhoneChanged:
		err = h.appendUserPhoneChangedEvent(event)
	case UserPhoneCodeAdded,
		HumanPhoneCodeAdded:
		err = h.appendUserPhoneCodeAddedEvent(event)
	case UserPhoneVerified,
		HumanPhoneVerified:
		h.appendUserPhoneVerifiedEvent()
	case UserPhoneRemoved,
		HumanPhoneRemoved:
		h.appendUserPhoneRemovedEvent()
	case UserAddressChanged,
		HumanAddressChanged:
		err = h.appendUserAddressChangedEvent(event)
	case MfaOtpAdded,
		HumanMfaOtpAdded:
		err = h.appendOtpAddedEvent(event)
	case MfaOtpVerified,
		HumanMfaOtpVerified:
		h.appendOtpVerifiedEvent()
	case MfaOtpRemoved,
		HumanMfaOtpRemoved:
		h.appendOtpRemovedEvent()
	}
	if err != nil {
		return err
	}
	h.ComputeObject()
	return nil
}

func (h *Human) ComputeObject() {
	if h.state == int32(model.UserStateUnspecified) {
		if h.Email != nil && h.IsEmailVerified {
			h.state = int32(model.UserStateActive)
		} else {
			h.state = int32(model.UserStateInitial)
		}
	}
	if h.Password != nil && h.Password.ObjectRoot.IsZero() {
		h.Password.ObjectRoot = h.objectRoot
	}
	if h.Profile != nil && h.Profile.ObjectRoot.IsZero() {
		h.Profile.ObjectRoot = h.objectRoot
	}
	if h.Email != nil && h.Email.ObjectRoot.IsZero() {
		h.Email.ObjectRoot = h.objectRoot
	}
	if h.Phone != nil && h.Phone.ObjectRoot.IsZero() {
		h.Phone.ObjectRoot = h.objectRoot
	}
	if h.Address != nil && h.Address.ObjectRoot.IsZero() {
		h.Address.ObjectRoot = h.objectRoot
	}
}

func (u *Human) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sj4jd", "could not unmarshal event")
	}
	return nil
}

func (u *Human) appendInitUsercodeCreatedEvent(event *es_models.Event) error {
	initCode := new(InitUserCode)
	err := initCode.SetData(event)
	if err != nil {
		return err
	}
	initCode.ObjectRoot.CreationDate = event.CreationDate
	u.InitCode = initCode
	return nil
}

func (c *InitUserCode) SetData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-7duwe").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lo34s", "could not unmarshal event")
	}
	return nil
}
