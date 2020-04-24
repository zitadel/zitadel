package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Password struct {
	es_models.ObjectRoot

	SecretString   string
	SecretCrypto   *crypto.CryptoValue
	ChangeRequired bool
}

type NotificationType int32

const (
	NOTIFICATIONTYPE_EMAIL NotificationType = iota
	NOTIFICATIONTYPE_SMS
)

func (p *Password) IsValid() bool {
	return p.AggregateID != "" && p.SecretString != ""
}
