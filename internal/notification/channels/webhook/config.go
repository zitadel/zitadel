package webhook

import (
	"net/http"
	"net/url"
)

type Config struct {
	CallURL    string
	Method     string
	Headers    http.Header
	SigningKey string
}

func (w *Config) Validate() error {
	_, err := url.Parse(w.CallURL)
	return err
}
