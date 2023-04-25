package domain

import (
	"io"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type EmailAddress string

func (e EmailAddress) Validate() error {
	if e == "" {
		return errors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
	}
	if !emailRegex.MatchString(string(e)) {
		return errors.ThrowInvalidArgument(nil, "EMAIL-599BI", "Errors.User.Email.Invalid")
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
	PlainCode       *string
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (e *Email) Validate() error {
	if e == nil {
		return errors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
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

// RenderConfirmURLTemplate parses and renders tmplStr.
// userID, code and orgID are passed into the [ConfirmURLData].
// "%s%s?userID=%s&code=%s&orgID=%s"
func RenderConfirmURLTemplate(w io.Writer, tmplStr, userID, code, orgID string) error {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return caos_errs.ThrowInvalidArgument(err, "USERv2-ooD8p", "Errors.User.Email.InvalidURLTemplate")
	}

	data := &ConfirmURLData{userID, code, orgID}
	if err = tmpl.Execute(w, data); err != nil {
		return caos_errs.ThrowInvalidArgument(err, "USERv2-ohSi5", "Errors.User.Email.InvalidURLTemplate")
	}

	return nil
}
