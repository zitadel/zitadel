package domain

import (
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type EmailAddress string

func (e EmailAddress) Validate() error {
	if e == "" {
		return zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
	}
	if !emailRegex.MatchString(string(e)) {
		return zerrors.ThrowInvalidArgument(nil, "EMAIL-599BI", "Errors.User.Email.Invalid")
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
