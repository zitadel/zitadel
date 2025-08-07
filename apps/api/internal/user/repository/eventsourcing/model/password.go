package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Password struct {
	es_models.ObjectRoot

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	EncodedHash    string              `json:"encodedHash,omitempty"`
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

func (u *Human) appendUserPasswordChangedEvent(event eventstore.Event) error {
	u.Password = new(Password)
	err := u.Password.setData(event)
	if err != nil {
		return err
	}
	u.Password.ObjectRoot.CreationDate = event.CreatedAt()
	return nil
}

func (u *Human) appendPasswordSetRequestedEvent(event eventstore.Event) error {
	u.PasswordCode = new(PasswordCode)
	return u.PasswordCode.SetData(event)
}

func (pw *Password) setData(event eventstore.Event) error {
	pw.ObjectRoot.AppendEvent(event)
	if err := event.Unmarshal(pw); err != nil {
		logging.Log("EVEN-dks93").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-sl9xlo2rsw", "could not unmarshal event")
	}
	return nil
}

func (c *PasswordCode) SetData(event eventstore.Event) error {
	c.ObjectRoot.AppendEvent(event)
	c.CreationDate = event.CreatedAt()
	if err := event.Unmarshal(c); err != nil {
		logging.Log("EVEN-lo0y2").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-q21dr", "could not unmarshal event")
	}
	return nil
}

func (pw *PasswordChange) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(pw); err != nil {
		logging.Log("EVEN-ADs31").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-BDd32", "could not unmarshal event")
	}
	pw.ObjectRoot.AppendEvent(event)
	return nil
}
