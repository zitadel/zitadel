package http

import (
	"net/http"

	"github.com/caos/zitadel/internal/errors"
)

type UserAgent struct {
	ID string
}

func (u *UserAgent) GetID() string {
	return u.ID
}

type Agent interface {
	GetID() string
}

type UserAgentHandler struct {
	handler    *CookieHandler
	cookieName string
}

type UserAgentCookieConfig struct {
	Name   string
	Domain string
	Key    string
}

func NewUserAgentHandler(config *UserAgentCookieConfig) *UserAgentHandler {
	handler := NewCookieHandler(
		WithEncryption([]byte(config.Key), []byte(config.Key)),
		WithDomain(config.Domain),
		WithUnsecure(),
	)
	return &UserAgentHandler{
		cookieName: config.Name,
		handler:    handler,
	}
}

func (ua *UserAgentHandler) GetUserAgent(r *http.Request) (Agent, error) {
	userAgent := new(UserAgent)
	err := ua.handler.GetEncryptedCookieValue(r, ua.cookieName, userAgent)
	if err != nil {
		return nil, errors.ThrowPermissionDenied(err, "HTTP-YULqH4", "cannot read user agent cookie")
	}
	return userAgent, nil
}

func (ua *UserAgentHandler) SetUserAgent(w http.ResponseWriter, agent Agent) error {
	err := ua.handler.SetEncryptedCookie(w, ua.cookieName, &UserAgent{ID: agent.GetID()})
	if err != nil {
		return errors.ThrowPermissionDenied(err, "HTTP-AqgqdA", "cannot set user agent cookie")
	}
	return nil
}
