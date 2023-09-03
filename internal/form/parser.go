package form

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/zitadel/zitadel/internal/errors"
)

type Parser struct {
	decoder *schema.Decoder
}

func NewParser() *Parser {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return &Parser{d}
}

func (p *Parser) Parse(r *http.Request, data interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return errors.ThrowInternal(err, "FORM-lCC9zI", "Errors.Internal")
	}

	return p.decoder.Decode(data, r.Form)
}
