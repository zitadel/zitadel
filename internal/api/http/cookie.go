package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"

	"github.com/zitadel/zitadel/internal/errors"
)

const (
	prefixSecure = "__Secure-"
	prefixHost   = "__Host-"
)

type CookieHandler struct {
	securecookie *securecookie.SecureCookie
	secureOnly   bool
	httpOnly     bool
	sameSite     http.SameSite
	path         string
	maxAge       int
}

func NewCookieHandler(opts ...CookieHandlerOpt) *CookieHandler {
	c := &CookieHandler{
		secureOnly: true,
		httpOnly:   true,
		sameSite:   http.SameSiteLaxMode,
		path:       "/",
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

type CookieHandlerOpt func(*CookieHandler)

func WithEncryption(hashKey, encryptKey []byte) CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.securecookie = securecookie.New(hashKey, encryptKey)
	}
}

func WithUnsecure() CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.secureOnly = false
	}
}

func WithNonHttpOnly() CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.httpOnly = false
	}
}

func WithSameSite(sameSite http.SameSite) CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.sameSite = sameSite
	}
}

func WithPath(path string) CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.path = path
	}
}

func WithMaxAge(maxAge int) CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.maxAge = maxAge
		if c.securecookie != nil {
			c.securecookie.MaxAge(maxAge)
		}
	}
}

func SetCookiePrefix(name, domain, path string, secureOnly bool) string {
	if !secureOnly {
		return name
	}
	if domain != "" || path != "/" {
		return prefixSecure + name
	}
	return prefixHost + name
}

func (c *CookieHandler) GetCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (c *CookieHandler) GetEncryptedCookieValue(r *http.Request, name string, value interface{}) error {
	cookie, err := r.Cookie(SetCookiePrefix(name, r.Host, c.path, c.secureOnly))
	if err != nil {
		return err
	}
	if c.securecookie == nil {
		return errors.ThrowInternal(nil, "HTTP-X6XpnL", "securecookie not configured")
	}
	return c.securecookie.Decode(name, cookie.Value, value)
}

func (c *CookieHandler) SetCookie(w http.ResponseWriter, name, domain, value string) {
	c.httpSet(w, name, domain, value, c.maxAge)
}

func (c *CookieHandler) SetEncryptedCookie(w http.ResponseWriter, name, domain string, value interface{}) error {
	if c.securecookie == nil {
		return errors.ThrowInternal(nil, "HTTP-s2HUtx", "securecookie not configured")
	}
	encoded, err := c.securecookie.Encode(name, value)
	if err != nil {
		return err
	}
	c.httpSet(w, name, domain, encoded, c.maxAge)
	return nil
}

func (c *CookieHandler) DeleteCookie(w http.ResponseWriter, r *http.Request, name string) {
	c.httpSet(w, name, r.Host, "", -1)
}

func (c *CookieHandler) httpSet(w http.ResponseWriter, name, domain, value string, maxage int) {
	http.SetCookie(w, &http.Cookie{
		Name:     SetCookiePrefix(name, domain, c.path, c.secureOnly),
		Value:    value,
		Domain:   strings.Split(domain, ":")[0],
		Path:     c.path,
		MaxAge:   maxage,
		HttpOnly: c.httpOnly,
		Secure:   c.secureOnly,
		SameSite: c.sameSite,
	})
	varyValues := w.Header().Values("vary")
	for _, vary := range varyValues {
		if vary == "Cookie" {
			return
		}
	}
	w.Header().Add("vary", "Cookie")
}
