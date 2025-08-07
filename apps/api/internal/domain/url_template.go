package domain

import (
	"io"
	"text/template"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func RenderURLTemplate(w io.Writer, tmpl string, data any) error {
	return renderURLTemplate(w, tmpl, data)
}

func renderURLTemplate(w io.Writer, tmpl string, data any) error {
	parsed, err := template.New("").Parse(tmpl)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate")
	}
	if err = parsed.Execute(w, data); err != nil {
		return zerrors.ThrowInvalidArgument(err, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate")
	}
	return nil
}
