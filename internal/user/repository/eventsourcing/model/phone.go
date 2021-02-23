package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string `json:"phone,omitempty"`
	IsPhoneVerified bool   `json:"-"`
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (p *Phone) Changes(changed *Phone) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.PhoneNumber != "" && p.PhoneNumber != changed.PhoneNumber {
		changes["phone"] = changed.PhoneNumber
	}
	return changes
}

func PhoneFromModel(phone *model.Phone) *Phone {
	return &Phone{
		ObjectRoot:      phone.ObjectRoot,
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneToModel(phone *Phone) *model.Phone {
	return &model.Phone{
		ObjectRoot:      phone.ObjectRoot,
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneCodeFromModel(code *model.PhoneCode) *PhoneCode {
	if code == nil {
		return nil
	}
	return &PhoneCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func PhoneCodeToModel(code *PhoneCode) *model.PhoneCode {
	return &model.PhoneCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func (u *Human) appendUserPhoneChangedEvent(event *es_models.Event) error {
	u.Phone = new(Phone)
	return u.Phone.setData(event)
}

func (u *Human) appendUserPhoneCodeAddedEvent(event *es_models.Event) error {
	u.PhoneCode = new(PhoneCode)
	return u.PhoneCode.SetData(event)
}

func (u *Human) appendUserPhoneVerifiedEvent() {
	u.IsPhoneVerified = true
}

func (u *Human) appendUserPhoneRemovedEvent() {
	u.Phone = nil
	u.PhoneCode = nil
}

func (p *Phone) setData(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-lco9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lre56", "could not unmarshal event")
	}
	return nil
}

func (c *PhoneCode) SetData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	c.CreationDate = event.CreationDate
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-sk8ws").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-7hdj3", "could not unmarshal event")
	}
	return nil
}
