package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type Password struct {
	es_models.ObjectRoot

	SecretString   string
	SecretCrypto   *crypto.CryptoValue
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

func (p *Password) HashPasswordIfExisting(policy *PasswordComplexityPolicy, passwordAlg crypto.HashAlgorithm) error {
	if p.SecretString == "" {
		return nil
	}
	if policy == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "DOMAIN-s8ifS", "Errors.User.PasswordComplexityPolicy.NotFound")
	}
	if err := policy.Check(p.SecretString); err != nil {
		return err
	}
	secret, err := crypto.Hash([]byte(p.SecretString), passwordAlg)
	if err != nil {
		return err
	}
	p.SecretCrypto = secret
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
