package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Human struct {
	user *User `json:"-"`

	*Password
	*Profile
	*Email
	*Phone
	*Address
	ExternalIDPs       []*ExternalIDP   `json:"-"`
	InitCode           *InitUserCode    `json:"-"`
	EmailCode          *EmailCode       `json:"-"`
	PhoneCode          *PhoneCode       `json:"-"`
	PasswordCode       *PasswordCode    `json:"-"`
	OTP                *OTP             `json:"-"`
	U2FTokens          []*WebAuthNToken `json:"-"`
	PasswordlessTokens []*WebAuthNToken `json:"-"`
	U2FLogins          []*WebAuthNLogin `json:"-"`
	PasswordlessLogins []*WebAuthNLogin `json:"-"`
}

type InitUserCode struct {
	es_models.ObjectRoot
	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
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
	switch event.Type() {
	case user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.UserV1ProfileChangedType,
		user.HumanAddedType,
		user.HumanRegisteredType,
		user.HumanProfileChangedType:
		err = h.setData(event)
	case user.UserV1InitialCodeAddedType,
		user.HumanInitialCodeAddedType:
		err = h.appendInitUsercodeCreatedEvent(event)
	case user.UserV1PasswordChangedType,
		user.HumanPasswordChangedType:
		err = h.appendUserPasswordChangedEvent(event)
	case user.UserV1PasswordCodeAddedType,
		user.HumanPasswordCodeAddedType:
		err = h.appendPasswordSetRequestedEvent(event)
	case user.UserV1EmailChangedType,
		user.HumanEmailChangedType:
		err = h.appendUserEmailChangedEvent(event)
	case user.UserV1EmailCodeAddedType,
		user.HumanEmailCodeAddedType:
		err = h.appendUserEmailCodeAddedEvent(event)
	case user.UserV1EmailVerifiedType,
		user.HumanEmailVerifiedType:
		h.appendUserEmailVerifiedEvent()
	case user.UserV1PhoneChangedType,
		user.HumanPhoneChangedType:
		err = h.appendUserPhoneChangedEvent(event)
	case user.UserV1PhoneCodeAddedType,
		user.HumanPhoneCodeAddedType:
		err = h.appendUserPhoneCodeAddedEvent(event)
	case user.UserV1PhoneVerifiedType,
		user.HumanPhoneVerifiedType:
		h.appendUserPhoneVerifiedEvent()
	case user.UserV1PhoneRemovedType,
		user.HumanPhoneRemovedType:
		h.appendUserPhoneRemovedEvent()
	case user.UserV1AddressChangedType,
		user.HumanAddressChangedType:
		err = h.appendUserAddressChangedEvent(event)
	case user.UserV1MFAOTPAddedType,
		user.HumanMFAOTPAddedType:
		err = h.appendOTPAddedEvent(event)
	case user.UserV1MFAOTPVerifiedType,
		user.HumanMFAOTPVerifiedType:
		h.appendOTPVerifiedEvent()
	case user.UserV1MFAOTPRemovedType,
		user.HumanMFAOTPRemovedType:
		h.appendOTPRemovedEvent()
	case user.UserIDPLinkAddedType:
		err = h.appendExternalIDPAddedEvent(event)
	case user.UserIDPLinkRemovedType, user.UserIDPLinkCascadeRemovedType:
		err = h.appendExternalIDPRemovedEvent(event)
	case user.HumanU2FTokenAddedType:
		err = h.appendU2FAddedEvent(event)
	case user.HumanU2FTokenVerifiedType:
		err = h.appendU2FVerifiedEvent(event)
	case user.HumanU2FTokenSignCountChangedType:
		err = h.appendU2FChangeSignCountEvent(event)
	case user.HumanU2FTokenRemovedType:
		err = h.appendU2FRemovedEvent(event)
	case user.HumanPasswordlessTokenAddedType:
		err = h.appendPasswordlessAddedEvent(event)
	case user.HumanPasswordlessTokenVerifiedType:
		err = h.appendPasswordlessVerifiedEvent(event)
	case user.HumanPasswordlessTokenSignCountChangedType:
		err = h.appendPasswordlessChangeSignCountEvent(event)
	case user.HumanPasswordlessTokenRemovedType:
		err = h.appendPasswordlessRemovedEvent(event)
	case user.HumanU2FTokenBeginLoginType:
		err = h.appendU2FLoginEvent(event)
	case user.HumanPasswordlessTokenBeginLoginType:
		err = h.appendPasswordlessLoginEvent(event)
	}
	if err != nil {
		return err
	}
	h.ComputeObject()
	return nil
}

func (h *Human) ComputeObject() {
	if h.user.State == int32(model.UserStateUnspecified) || h.user.State == int32(model.UserStateInitial) {
		if h.Email != nil && h.IsEmailVerified {
			h.user.State = int32(model.UserStateActive)
		} else {
			h.user.State = int32(model.UserStateInitial)
		}
	}
	if h.Password != nil && h.Password.ObjectRoot.IsZero() {
		h.Password.ObjectRoot = h.user.ObjectRoot
	}
	if h.Profile != nil && h.Profile.ObjectRoot.IsZero() {
		h.Profile.ObjectRoot = h.user.ObjectRoot
	}
	if h.Email != nil && h.Email.ObjectRoot.IsZero() {
		h.Email.ObjectRoot = h.user.ObjectRoot
	}
	if h.Phone != nil && h.Phone.ObjectRoot.IsZero() {
		h.Phone.ObjectRoot = h.user.ObjectRoot
	}
	if h.Address != nil && h.Address.ObjectRoot.IsZero() {
		h.Address.ObjectRoot = h.user.ObjectRoot
	}
}

func (u *Human) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-8ujgd").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-sj4jd", "could not unmarshal event")
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
		return zerrors.ThrowInternal(err, "MODEL-lo34s", "could not unmarshal event")
	}
	return nil
}
