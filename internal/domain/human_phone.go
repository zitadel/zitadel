package domain

import (
	"time"

	"github.com/ttacon/libphonenumber"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const defaultRegion = "CH"

type PhoneNumber string

func (p PhoneNumber) Normalize() (PhoneNumber, error) {
	if p == "" {
		return p, zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty")
	}
	phoneNr, err := libphonenumber.Parse(string(p), defaultRegion)
	if err != nil {
		return p, zerrors.ThrowInvalidArgument(err, "PHONE-so0wa", "Errors.User.Phone.Invalid")
	}
	return PhoneNumber(libphonenumber.Format(phoneNr, libphonenumber.E164)), nil
}

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     PhoneNumber
	IsPhoneVerified bool
	// PlainCode is set by the command and can be used to return it to the caller (API)
	PlainCode *string
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (p *Phone) Normalize() error {
	if p == nil {
		return zerrors.ThrowInvalidArgument(nil, "PHONE-YlbwO", "Errors.User.Phone.Empty")
	}
	normalizedNumber, err := p.PhoneNumber.Normalize()
	if err != nil {
		return err
	}
	// Issue for avoiding mutating state: https://github.com/zitadel/zitadel/issues/5412
	p.PhoneNumber = normalizedNumber
	return nil
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
