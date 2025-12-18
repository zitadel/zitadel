package domain

import (
	"errors"
	"io"
	"net/mail"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type EmailAddress string

var errEmailNameProvided = errors.New("email with name provided")

func (e EmailAddress) Validate() error {
	if e == "" {
		return zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
	}
	p, err := mail.ParseAddress(string(e))
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "EMAIL-599BI", "Errors.User.Email.Invalid")
	}
	if p.Name != "" {
		return zerrors.ThrowInvalidArgument(errEmailNameProvided, "EMAIL-599GV", "Errors.User.Email.Invalid")
	}
	return nil
}

func (e EmailAddress) Normalize() EmailAddress {
	return EmailAddress(strings.TrimSpace(string(e)))
}

type Email struct {
	es_models.ObjectRoot

	EmailAddress    EmailAddress
	IsEmailVerified bool
	// PlainCode is set by the command and can be used to return it to the caller (API)
	PlainCode *string
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (e *Email) Validate() error {
	if e == nil {
		return zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
	}
	return e.EmailAddress.Validate()
}

func NewEmailCode(emailGenerator crypto.Generator) (*EmailCode, string, error) {
	emailCodeCrypto, code, err := crypto.NewCode(emailGenerator)
	if err != nil {
		return nil, "", err
	}
	return &EmailCode{
		Code:   emailCodeCrypto,
		Expiry: emailGenerator.Expiry(),
	}, code, nil
}

type ConfirmURLData struct {
	UserID string
	Code   string
	OrgID  string
}

// RenderConfirmURLTemplate parses and renders tmpl.
// userID, code and orgID are passed into the [ConfirmURLData].
func RenderConfirmURLTemplate(w io.Writer, tmpl, userID, code, orgID string) error {
	return renderURLTemplate(w, tmpl, &ConfirmURLData{userID, code, orgID})
}
