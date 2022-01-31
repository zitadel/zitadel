package middleware

import (
	"context"
	"net/http"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/v2/internal/config/types"
)

type cookieKey int

var (
	userAgentKey cookieKey = 0
)

func UserAgentIDFromCtx(ctx context.Context) (string, bool) {
	userAgentID, ok := ctx.Value(userAgentKey).(string)
	return userAgentID, ok
}

type UserAgent struct {
	ID string
}

type userAgentHandler struct {
	cookieHandler *http_utils.CookieHandler
	cookieName    string
	idGenerator   id.Generator
	nextHandler   http.Handler
}

type UserAgentCookieConfig struct {
	Name   string
	Key    *crypto.KeyConfig
	MaxAge types.Duration
}

func NewUserAgentHandler(config *UserAgentCookieConfig, domain string, idGenerator id.Generator, localDevMode bool) (func(http.Handler) http.Handler, error) {
	key, err := crypto.LoadKey(config.Key, config.Key.EncryptionKeyID)
	if err != nil {
		return nil, err
	}
	cookieKey := []byte(key)
	opts := []http_utils.CookieHandlerOpt{
		http_utils.WithEncryption(cookieKey, cookieKey),
		http_utils.WithDomain(domain),
		http_utils.WithMaxAge(int(config.MaxAge.Seconds())),
	}
	if localDevMode {
		opts = append(opts, http_utils.WithUnsecure())
	}
	return func(handler http.Handler) http.Handler {
		return &userAgentHandler{
			nextHandler:   handler,
			cookieName:    config.Name,
			cookieHandler: http_utils.NewCookieHandler(opts...),
			idGenerator:   idGenerator,
		}
	}, nil
}

func (ua *userAgentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agent, err := ua.getUserAgent(r)
	if err != nil {
		agent, err = ua.newUserAgent()
	}
	if err == nil {
		ctx := context.WithValue(r.Context(), userAgentKey, agent.ID)
		r = r.WithContext(ctx)
		ua.setUserAgent(w, agent)
	}
	ua.nextHandler.ServeHTTP(w, r)
}

func (ua *userAgentHandler) newUserAgent() (*UserAgent, error) {
	agentID, err := ua.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	return &UserAgent{ID: agentID}, nil
}

func (ua *userAgentHandler) getUserAgent(r *http.Request) (*UserAgent, error) {
	userAgent := new(UserAgent)
	err := ua.cookieHandler.GetEncryptedCookieValue(r, ua.cookieName, userAgent)
	if err != nil {
		return nil, errors.ThrowPermissionDenied(err, "HTTP-YULqH4", "cannot read user agent cookie")
	}
	return userAgent, nil
}

func (ua *userAgentHandler) setUserAgent(w http.ResponseWriter, agent *UserAgent) error {
	err := ua.cookieHandler.SetEncryptedCookie(w, ua.cookieName, agent)
	if err != nil {
		return errors.ThrowPermissionDenied(err, "HTTP-AqgqdA", "cannot set user agent cookie")
	}
	return nil
}
