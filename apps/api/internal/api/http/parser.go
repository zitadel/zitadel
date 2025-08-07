package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/zitadel/zitadel/internal/zerrors"
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
		return zerrors.ThrowInternal(err, "FORM-lCC9zI", "error parsing http form")
	}

	return p.decoder.Decode(data, r.Form)
}

func (p *Parser) UnwrapParserError(err error) error {
	if err == nil {
		return nil
	}

	// try to unwrap the error
	var multiErr schema.MultiError
	if errors.As(err, &multiErr) && len(multiErr) == 1 {
		for _, v := range multiErr {
			var schemaErr schema.ConversionError
			if errors.As(v, &schemaErr) {
				return schemaErr.Err
			}

			return v
		}
	}

	return err
}
