package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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
	return e.EmailAddress != "" && emailRegex.MatchString(e.EmailAddress)
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
