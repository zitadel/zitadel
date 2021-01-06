package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string
	IsEmailVerified bool
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (e *Email) IsValid() bool {
	return e.EmailAddress != ""
}

func NewEmailCode(emailGenerator crypto.Generator) (*EmailCode, error) {
	emailCodeCrypto, _, err := crypto.NewCode(emailGenerator)
	if err != nil {
		return nil, err
	}
	return &EmailCode{
		Code:   emailCodeCrypto,
		Expiry: emailGenerator.Expiry(),
	}, nil
}
