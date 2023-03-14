package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
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

func (e *Email) GenerateEmailCodeIfNeeded(emailGenerator crypto.Generator) (*EmailCode, error) {
	var emailCode *EmailCode
	if e.IsEmailVerified {
		return emailCode, nil
	}
	emailCode = new(EmailCode)
	return emailCode, emailCode.GenerateEmailCode(emailGenerator)
}

func (code *EmailCode) GenerateEmailCode(emailGenerator crypto.Generator) error {
	emailCodeCrypto, _, err := crypto.NewCode(emailGenerator)
	if err != nil {
		return err
	}
	code.Code = emailCodeCrypto
	code.Expiry = emailGenerator.Expiry()
	return nil
}
