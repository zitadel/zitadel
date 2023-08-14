package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Password struct {
	es_models.ObjectRoot

	SecretString   string
	EncodedSecret  string
	ChangeRequired bool
}

func NewPassword(password string) *Password {
	return &Password{
		SecretString: password,
	}
}

type PasswordCode struct {
	es_models.ObjectRoot

	Code             *crypto.CryptoValue
	Expiry           time.Duration
	NotificationType NotificationType
}

func (p *Password) HashPasswordIfExisting(policy *PasswordComplexityPolicy, hasher *crypto.PasswordHasher) error {
	if p.SecretString == "" {
		return nil
	}
	if policy == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "DOMAIN-s8ifS", "Errors.User.PasswordComplexityPolicy.NotFound")
	}
	if err := policy.Check(p.SecretString); err != nil {
		return err
	}
	encoded, err := hasher.Hash(p.SecretString)
	if err != nil {
		return err
	}
	p.EncodedSecret = encoded
	return nil
}

func NewPasswordCode(passwordGenerator crypto.Generator) (*PasswordCode, error) {
	passwordCodeCrypto, _, err := crypto.NewCode(passwordGenerator)
	if err != nil {
		return nil, err
	}
	return &PasswordCode{
		Code:   passwordCodeCrypto,
		Expiry: passwordGenerator.Expiry(),
	}, nil
}
