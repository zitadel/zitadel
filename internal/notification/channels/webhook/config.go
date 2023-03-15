package webhook

import (
	"net/url"
)

type Config struct {
	CallURL string
	Method  string
}

func (w *Config) IsValid() error {
	_, err := url.Parse(w.CallURL)
	return err
}
