package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"

	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	PrefixSecure cookiePrefix = "__Secure-"
	PrefixHost   cookiePrefix = "__Host-"
)

type cookiePrefix string

type CookieHandler struct {
	securecookie *securecookie.SecureCookie
	secureOnly   bool
	httpOnly     bool
	sameSite     http.SameSite
	path         string
	maxAge       int
	prefix       cookiePrefix
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

func WithPrefix(prefix cookiePrefix) CookieHandlerOpt {
	return func(c *CookieHandler) {
		c.prefix = prefix
	}
}

func SetCookiePrefix(name string, secureOnly bool, prefix cookiePrefix) string {
	if !secureOnly {
		return name
	}
	return string(prefix) + name
}

func (c *CookieHandler) GetCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (c *CookieHandler) GetEncryptedCookieValue(r *http.Request, name string, value interface{}) error {
	cookie, err := r.Cookie(SetCookiePrefix(name, c.secureOnly, c.prefix))
	if err != nil {
		return err
	}
	if c.securecookie == nil {
		return zerrors.ThrowInternal(nil, "HTTP-X6XpnL", "securecookie not configured")
	}
	return c.securecookie.Decode(name, cookie.Value, value)
}

func (c *CookieHandler) SetCookie(w http.ResponseWriter, name, domain, value string) {
	c.httpSet(w, name, domain, value, c.maxAge)
}

func (c *CookieHandler) SetEncryptedCookie(w http.ResponseWriter, name, domain string, value interface{}, sameSiteNone bool) error {
	if c.securecookie == nil {
		return zerrors.ThrowInternal(nil, "HTTP-s2HUtx", "securecookie not configured")
	}
	encoded, err := c.securecookie.Encode(name, value)
	if err != nil {
		return err
	}
	sameSite := c.sameSite
	if sameSiteNone {
		sameSite = http.SameSiteNoneMode
	}
	c.httpSetWithSameSite(w, name, domain, encoded, c.maxAge, sameSite)
	return nil
}

func (c *CookieHandler) DeleteCookie(w http.ResponseWriter, name string) {
	c.httpSet(w, name, "", "", -1)
}

func (c *CookieHandler) httpSet(w http.ResponseWriter, name, domain, value string, maxAge int) {
	c.httpSetWithSameSite(w, name, domain, value, maxAge, c.sameSite)
}

func (c *CookieHandler) httpSetWithSameSite(w http.ResponseWriter, name, host, value string, maxAge int, sameSite http.SameSite) {
	domain := strings.Split(host, ":")[0]
	// same site none requires the secure flag, so we'll set it even if the cookie is set on non-TLS for localhost
	secure := c.secureOnly || (sameSite == http.SameSiteNoneMode && domain == "localhost")
	// prefix the cookie for secure cookies (TLS only, therefore not for samesite none on http://localhost)
	prefixedName := SetCookiePrefix(name, c.secureOnly, c.prefix)
	// in case the host prefix is set, we need to make sure the domain is not set (otherwise the browser will reject the cookie)
	if secure && c.prefix == PrefixHost {
		domain = ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:     prefixedName,
		Value:    value,
		Domain:   domain,
		Path:     c.path,
		MaxAge:   maxAge,
		HttpOnly: c.httpOnly,
		Secure:   secure,
		SameSite: sameSite,
	})
	varyValues := w.Header().Values("vary")
	for _, vary := range varyValues {
		if vary == "Cookie" {
			return
		}
	}
	w.Header().Add("vary", "Cookie")
}
