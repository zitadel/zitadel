package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		return zerrors.ThrowInternal(err, "MODEL-lre56", "could not unmarshal event")
	}
	return nil
}

func (c *PhoneCode) SetData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	c.CreationDate = event.CreationDate
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-sk8ws").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-7hdj3", "could not unmarshal event")
	}
	return nil
}
