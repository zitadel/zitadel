package http

import (
	"net/http"

	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/id"
)

type UserAgent struct {
	ID string
}

type UserAgentHandler struct {
	handler     *CookieHandler
	cookieName  string
	idGenerator id.Generator
}

type UserAgentCookieConfig struct {
	Name   string
	Domain string
	Key    *crypto.KeyConfig
	MaxAge types.Duration
}

func NewUserAgentHandler(config *UserAgentCookieConfig, idGenerator id.Generator, localDevMode bool) (*UserAgentHandler, error) {
	key, err := crypto.LoadKey(config.Key, config.Key.EncryptionKeyID)
	if err != nil {
		return nil, err
	}
	cookieKey := []byte(key)
	opts := []CookieHandlerOpt{
		WithEncryption(cookieKey, cookieKey),
		WithDomain(config.Domain),
		WithMaxAge(int(config.MaxAge.Seconds())),
	}
	if localDevMode {
		opts = append(opts, WithUnsecure())
	}
	handler := NewCookieHandler(opts...)
	return &UserAgentHandler{
		cookieName:  config.Name,
		handler:     handler,
		idGenerator: idGenerator,
	}, nil
}

func (ua *UserAgentHandler) NewUserAgent() (*UserAgent, error) {
	agentID, err := ua.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	return &UserAgent{ID: agentID}, nil
}

func (ua *UserAgentHandler) GetUserAgent(r *http.Request) (*UserAgent, error) {
	userAgent := new(UserAgent)
	err := ua.handler.GetEncryptedCookieValue(r, ua.cookieName, userAgent)
	if err != nil {
		return nil, errors.ThrowPermissionDenied(err, "HTTP-YULqH4", "cannot read user agent cookie")
	}
	return userAgent, nil
}

func (ua *UserAgentHandler) SetUserAgent(w http.ResponseWriter, agent *UserAgent) error {
	err := ua.handler.SetEncryptedCookie(w, ua.cookieName, agent)
	if err != nil {
		return errors.ThrowPermissionDenied(err, "HTTP-AqgqdA", "cannot set user agent cookie")
	}
	return nil
}
