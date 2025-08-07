package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		return zerrors.ThrowInternal(err, "MODEL-sl9xw", "could not unmarshal event")
	}
	return nil
}

func (a *EmailCode) SetData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	a.CreationDate = event.CreationDate
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo9s").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-s8uws", "could not unmarshal event")
	}
	return nil
}
