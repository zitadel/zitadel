package set

import (
	"net/url"
)

type Config struct {
	CallURL string
}

func (w *Config) Validate() error {
	_, err := url.Parse(w.CallURL)
	return err
}
