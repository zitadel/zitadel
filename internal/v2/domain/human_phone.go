package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/ttacon/libphonenumber"
	"time"
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

func (p *Phone) IsValid() bool {
	err := p.formatPhone()
	return p.PhoneNumber != "" && err == nil
}

func (p *Phone) formatPhone() error {
	phoneNr, err := libphonenumber.Parse(p.PhoneNumber, defaultRegion)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-so0wa", "Errors.User.Phone.Invalid")
	}
	p.PhoneNumber = libphonenumber.Format(phoneNr, libphonenumber.E164)
	return nil
}
