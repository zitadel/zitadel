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

func UserAgentIDFromCtx(ctx context.Context) (string, bool) {
	userAgentID, ok := ctx.Value(userAgentKey).(string)
	return userAgentID, ok
}

func UserAgentCookieHandler(cookieHandler *http_utils.UserAgentHandler, nextHandlerFunc func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ua, err := cookieHandler.GetUserAgent(r)
			if err != nil {
				ua, err = cookieHandler.NewUserAgent()
			}
			if err == nil {
				ctx := context.WithValue(r.Context(), userAgentKey, ua.ID)
				r = r.WithContext(ctx)
				cookieHandler.SetUserAgent(w, ua)
			}
			handlerFunc(w, r)
			nextHandlerFunc(handlerFunc)
		}
	}
}
