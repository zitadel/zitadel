package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	defaultRegion = "CH"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string
	IsPhoneVerified bool
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (p *Phone) GeneratePhoneCodeIfNeeded(phoneGenerator crypto.Generator) (*PhoneCode, error) {
	var phoneCode *PhoneCode
	if p.IsPhoneVerified {
		return phoneCode, nil
	}
	phoneCode = new(PhoneCode)
	return phoneCode, phoneCode.GeneratePhoneCode(phoneGenerator)
}

func (code *PhoneCode) GeneratePhoneCode(phoneGenerator crypto.Generator) error {
	phoneCodeCrypto, _, err := crypto.NewCode(phoneGenerator)
	if err != nil {
		return err
	}
	code.Code = phoneCodeCrypto
	code.Expiry = phoneGenerator.Expiry()
	return nil
}
