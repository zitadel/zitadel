package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

type Password struct {
	es_models.ObjectRoot

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

type RequestPasswordSet struct {
	es_models.ObjectRoot

	Code             *crypto.CryptoValue `json:"code,omitempty"`
	Expiry           time.Duration       `json:"expiry,omitempty"`
	NotificationType int32               `json:"notificationType,omitempty"`
}

func PasswordFromModel(password *model.Password) *Password {
	return &Password{
		ObjectRoot: password.ObjectRoot,
		Secret:         password.SecretCrypto,
		ChangeRequired: password.ChangeRequired,
	}
}

func PasswordToModel(password *Password) *model.Password {
	return &model.Password{
		ObjectRoot: password.ObjectRoot,
		SecretCrypto:   password.Secret,
		ChangeRequired: password.ChangeRequired,
	}
}

func PasswordCodeToModel(code *RequestPasswordSet) *model.RequestPasswordSet {
	return &model.RequestPasswordSet{
		ObjectRoot: code.ObjectRoot,
		Expiry:           code.Expiry,
		Code:             code.Code,
		NotificationType: model.NotificationType(code.NotificationType),
	}
}

func (u *User) appendUserPasswordChangedEvent(event *es_models.Event) error {
	pw := new(Password)
	err := pw.setData(event)
	if err != nil {
		return err
	}
	pw.ObjectRoot.CreationDate = event.CreationDate
	u.Password = pw
	return nil
}

func (u *User) appendPasswordSetRequestedEvent(event *es_models.Event) error {
	u.PasswordCode = new(RequestPasswordSet)
	u.PasswordCode.setData(event)
	return nil
}

func (pw *Password) setData(event *es_models.Event) error {
	pw.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, pw); err != nil {
		logging.Log("EVEN-dks93").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (a *RequestPasswordSet) setData(event *es_models.Event) error {
	a.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo0y2").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
