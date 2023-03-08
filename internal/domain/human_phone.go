package domain

import (
	"time"

	"github.com/ttacon/libphonenumber"

	"github.com/zitadel/zitadel/internal/crypto"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
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

func (p *Phone) Validate() error {
	if p == nil || p.PhoneNumber == "" {
		return caos_errs.ThrowInvalidArgument(nil, "PHONE-YlbwO", "Errors.User.Phone.Empty")
	}
	return p.formatPhone()
}

func (p *Phone) formatPhone() error {
	phoneNr, err := libphonenumber.Parse(p.PhoneNumber, defaultRegion)
	if err != nil {
		return caos_errs.ThrowInvalidArgument(err, "PHONE-so0wa", "Errors.User.Phone.Invalid")
	}
	p.PhoneNumber = libphonenumber.Format(phoneNr, libphonenumber.E164)
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
