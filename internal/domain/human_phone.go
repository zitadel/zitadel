package domain

import (
	"time"

	"github.com/ttacon/libphonenumber"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const defaultRegion = "CH"

type PhoneNumber string

func (p PhoneNumber) Normalize() (PhoneNumber, error) {
	if p == "" {
		return p, caos_errs.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty")
	}
	phoneNr, err := libphonenumber.Parse(string(p), defaultRegion)
	if err != nil {
		return p, caos_errs.ThrowInvalidArgument(err, "PHONE-so0wa", "Errors.User.Phone.Invalid")
	}
	return PhoneNumber(libphonenumber.Format(phoneNr, libphonenumber.E164)), nil
}

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     PhoneNumber
	IsPhoneVerified bool
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (p *Phone) Normalize() error {
	if p == nil {
		return caos_errs.ThrowInvalidArgument(nil, "PHONE-YlbwO", "Errors.User.Phone.Empty")
	}
	normalizedNumber, err := p.PhoneNumber.Normalize()
	if err != nil {
		return err
	}
	// Issue for avoiding mutating state: https://github.com/zitadel/zitadel/issues/5412
	p.PhoneNumber = normalizedNumber
	return nil
}

func NewPhoneCode(phoneGenerator crypto.Generator) (*PhoneCode, error) {
	phoneCodeCrypto, _, err := crypto.NewCode(phoneGenerator)
	if err != nil {
		return nil, err
	}
	return &PhoneCode{
		Code:   phoneCodeCrypto,
		Expiry: phoneGenerator.Expiry(),
	}, nil
}

type PhoneState int32

const (
	PhoneStateUnspecified PhoneState = iota
	PhoneStateActive
	PhoneStateRemoved

	phoneStateCount
)

func (s PhoneState) Valid() bool {
	return s >= 0 && s < phoneStateCount
}

func (s PhoneState) Exists() bool {
	return s == PhoneStateActive
}
