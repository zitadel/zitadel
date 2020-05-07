package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
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
	NOTIFICATIONTYPE_EMAIL NotificationType = iota
	NOTIFICATIONTYPE_SMS
)

func (p *Password) IsValid() bool {
	return p.AggregateID != "" && p.SecretString != ""
}

func (p *Password) HashPasswordIfExisting(passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if p.SecretString == "" {
		return nil
	}
	secret, err := crypto.Hash([]byte(p.SecretString), passwordAlg)
	if err != nil {
		return err
	}
	p.SecretCrypto = secret
	p.ChangeRequired = onetime
	return nil
}
