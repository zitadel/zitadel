package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string
	IsPhoneVerified bool
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (p *Phone) IsValid() bool {
	return p.PhoneNumber != ""
}

func (u *User) appendUserPhoneChangedEvent(event *es_models.Event) error {
	u.Phone = new(Phone)
	u.Phone.setData(event)
	u.IsPhoneVerified = false
	return nil
}

func (u *User) appendUserPhoneVerifiedEvent() error {
	u.IsPhoneVerified = true
	return nil
}

func (p *Phone) setData(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *Phone) GeneratePhoneCodeIfNeeded(phoneGenerator crypto.Generator) (*PhoneCode, error) {
	var phoneCode *PhoneCode
	if p.IsPhoneVerified {
		return phoneCode, nil
	}
	phoneCode = new(PhoneCode)
	return phoneCode, phoneCode.GeneratePhoneCode(phoneGenerator)
}

func (code *PhoneCode) GeneratePhoneCode(phoneGenerator crypto.Generator) error {
	phoneCodeCrypto, _, err := crypto.NewCode(phoneGenerator)
	if err != nil {
		return err
	}
	code.Code = phoneCodeCrypto
	code.Expiry = phoneGenerator.Expiry()
	return nil
}
