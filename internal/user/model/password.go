package model

import (
	"time"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	es_models "github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type Password struct {
	es_models.ObjectRoot

	SecretString   string
	SecretCrypto   *crypto.CryptoValue
	ChangeRequired bool
}

type PasswordCode struct {
	es_models.ObjectRoot

	Code             *crypto.CryptoValue
	Expiry           time.Duration
	NotificationType NotificationType
}

type NotificationType int32

const (
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms
)

func (p *Password) IsValid() bool {
	return p.AggregateID != "" && p.SecretString != ""
}
