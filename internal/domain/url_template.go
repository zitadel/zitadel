package domain

import (
	"io"
	"text/template"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func renderURLTemplate(w io.Writer, tmpl string, data any) error {
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return caos_errs.ThrowInvalidArgument(err, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate")
	}
	if err = parsed.Execute(w, data); err != nil {
		return caos_errs.ThrowInvalidArgument(err, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate")
	}
	return nil
}
