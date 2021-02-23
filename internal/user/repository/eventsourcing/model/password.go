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

type Password struct {
	es_models.ObjectRoot

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

type PasswordCode struct {
	es_models.ObjectRoot

	Code             *crypto.CryptoValue `json:"code,omitempty"`
	Expiry           time.Duration       `json:"expiry,omitempty"`
	NotificationType int32               `json:"notificationType,omitempty"`
}

type PasswordChange struct {
	Password
	UserAgentID string `json:"userAgentID,omitempty"`
}

func PasswordFromModel(password *model.Password) *Password {
	return &Password{
		ObjectRoot:     password.ObjectRoot,
		Secret:         password.SecretCrypto,
		ChangeRequired: password.ChangeRequired,
	}
}

func PasswordToModel(password *Password) *model.Password {
	return &model.Password{
		ObjectRoot:     password.ObjectRoot,
		SecretCrypto:   password.Secret,
		ChangeRequired: password.ChangeRequired,
	}
}

func PasswordCodeToModel(code *PasswordCode) *model.PasswordCode {
	return &model.PasswordCode{
		ObjectRoot:       code.ObjectRoot,
		Expiry:           code.Expiry,
		Code:             code.Code,
		NotificationType: model.NotificationType(code.NotificationType),
	}
}

func PasswordChangeFromModel(password *model.Password, userAgentID string) *PasswordChange {
	return &PasswordChange{
		Password: Password{
			ObjectRoot:     password.ObjectRoot,
			Secret:         password.SecretCrypto,
			ChangeRequired: password.ChangeRequired,
		},
		UserAgentID: userAgentID,
	}
}

func (u *Human) appendUserPasswordChangedEvent(event *es_models.Event) error {
	u.Password = new(Password)
	err := u.Password.setData(event)
	if err != nil {
		return err
	}
	u.Password.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (u *Human) appendPasswordSetRequestedEvent(event *es_models.Event) error {
	u.PasswordCode = new(PasswordCode)
	return u.PasswordCode.SetData(event)
}

func (pw *Password) setData(event *es_models.Event) error {
	pw.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, pw); err != nil {
		logging.Log("EVEN-dks93").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-sl9xlo2rsw", "could not unmarshal event")
	}
	return nil
}

func (c *PasswordCode) SetData(event *es_models.Event) error {
	c.ObjectRoot.AppendEvent(event)
	c.CreationDate = event.CreationDate
	if err := json.Unmarshal(event.Data, c); err != nil {
		logging.Log("EVEN-lo0y2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-q21dr", "could not unmarshal event")
	}
	return nil
}

func (pw *PasswordChange) SetData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, pw); err != nil {
		logging.Log("EVEN-ADs31").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-BDd32", "could not unmarshal event")
	}
	pw.ObjectRoot.AppendEvent(event)
	return nil
}
