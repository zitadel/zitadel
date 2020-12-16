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

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string `json:"email,omitempty"`
	IsEmailVerified bool   `json:"-"`
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *Email) Changes(changed *Email) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.EmailAddress != "" && e.EmailAddress != changed.EmailAddress {
		changes["email"] = changed.EmailAddress
	}
	return changes
}

func EmailFromModel(email *model.Email) *Email {
	return &Email{
		ObjectRoot:      email.ObjectRoot,
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func EmailToModel(email *Email) *model.Email {
	return &model.Email{
		ObjectRoot:      email.ObjectRoot,
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func EmailCodeFromModel(code *model.EmailCode) *EmailCode {
	if code == nil {
		return nil
	}
	return &EmailCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func EmailCodeToModel(code *EmailCode) *model.EmailCode {
	return &model.EmailCode{
		ObjectRoot: code.ObjectRoot,
		Expiry:     code.Expiry,
		Code:       code.Code,
	}
}

func (u *Human) appendUserEmailChangedEvent(event *es_models.Event) error {
	u.Email = new(Email)
	return u.Email.setData(event)
}

func (u *Human) appendUserEmailCodeAddedEvent(event *es_models.Event) error {
	u.EmailCode = new(EmailCode)
	return u.EmailCode.SetData(event)
}

func (u *Human) appendUserEmailVerifiedEvent() {
	u.IsEmailVerified = true
}

func (a *Email) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-dlo9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sl9xw", "could not unmarshal event")
	}
	return nil
}

func (a *EmailCode) SetData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	a.CreationDate = event.CreationDate
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-s8uws", "could not unmarshal event")
	}
	return nil
}
