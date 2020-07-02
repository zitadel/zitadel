package http

import (
	"net/http"

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
}

func NewUserAgentHandler(config *UserAgentCookieConfig, idGenerator id.Generator) (*UserAgentHandler, error) {
	key, err := crypto.LoadKey(config.Key, config.Key.EncryptionKeyID)
	if err != nil {
		return nil, err
	}
	cookieKey := []byte(key)
	handler := NewCookieHandler(
		WithEncryption(cookieKey, cookieKey),
		WithDomain(config.Domain),
		WithUnsecure(),
	)
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
