package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string `json:"email,omitempty"`
	IsEmailVerified bool   `json:"-"`

	isEmailUnique bool `json:"-"`
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
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  email.ObjectRoot.AggregateID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func EmailToModel(email *Email) *model.Email {
	return &model.Email{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  email.ObjectRoot.AggregateID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func (u *User) appendUserEmailChangedEvent(event *es_models.Event) error {
	u.Email = new(Email)
	u.Email.setData(event)
	u.IsEmailVerified = false
	return nil
}

func (u *User) appendUserEmailCodeAddedEvent(event *es_models.Event) error {
	u.EmailCode = new(EmailCode)
	u.EmailCode.setData(event)
	return nil
}

func (u *User) appendUserEmailVerifiedEvent() error {
	u.IsEmailVerified = true
	return nil
}

func (a *Email) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-dlo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (a *EmailCode) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo9s").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
