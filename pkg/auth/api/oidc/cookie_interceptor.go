package oidc

import (
	"context"
	"net/http"

	http_utils "github.com/caos/zitadel/internal/api/http"
)

type key int

var (
	userAgentKey key
)

const (
	UserAgentCookie  = "ua.cookie"
	UserAgentContext = "ua"
)

type UserAgentStorage interface {
	//GetUserAgent(context.Context, string) (*model.Agent, error)
	//CreateUserAgent(context.Context, *model.CreateAgent) (*model.Agent, error)
}

func UserAgentIDFromCtx(ctx context.Context) (string, bool) {
	userAgentID, ok := ctx.Value(userAgentKey).(string)
	return userAgentID, ok
}

func UserAgentCookieHandler(cookieHandler *http_utils.UserAgentHandler, storage UserAgentStorage, nextHandlerFunc func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//ua, err := cookieHandler.GetUserAgent(r)
			//var agent *model.Agent
			//if err == nil {
			//	agent, err = storage.GetUserAgent(r.Context(), ua.GetID())
			//}
			//if err != nil {
			//	agent, err = storage.CreateUserAgent(r.Context(), &model.CreateAgent{})
			//}
			//if err == nil {
			//	ctx := context.WithValue(r.Context(), userAgentKey, agent.ID)
			//	r = r.WithContext(ctx)
			//	cookieHandler.SetUserAgent(w, agent)
			//}
			handlerFunc(w, r)
			nextHandlerFunc(handlerFunc)
		}
	}
}
