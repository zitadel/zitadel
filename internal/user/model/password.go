package model

import (
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
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
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms
)

func (p *Password) IsValid() bool {
	return p.AggregateID != "" && p.SecretString != ""
}

func (p *Password) HashPasswordIfExisting(policy *iam_model.PasswordComplexityPolicyView, passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if p.SecretString == "" {
		return nil
	}
	if policy == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "MODEL-s8ifS", "Errors.User.PasswordComplexityPolicy.NotFound")
	}
	if err := policy.Check(p.SecretString); err != nil {
		return err
	}
	secret, err := crypto.Hash([]byte(p.SecretString), passwordAlg)
	if err != nil {
		return err
	}
	p.SecretCrypto = secret
	p.ChangeRequired = onetime
	return nil
}
