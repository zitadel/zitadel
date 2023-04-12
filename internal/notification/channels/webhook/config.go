package webhook

import (
	"net/url"
)

type Config struct {
	CallURL string
	Method  string
}

func (w *Config) Validate() error {
	_, err := url.Parse(w.CallURL)
	return err
}
