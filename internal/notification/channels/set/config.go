package set

import (
	"net/http"
	"net/url"
)

type Config struct {
	CallURL string
	Client  *http.Client
}

func (w *Config) Validate() error {
	_, err := url.Parse(w.CallURL)
	return err
}
