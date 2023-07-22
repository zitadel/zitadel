package middleware

import (
	"net/http"
)

func ModeHandler() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if mode := r.URL.Query().Get("mode"); r.URL.Path == "/oauth/v2/authorize" && mode != "" {
				http.SetCookie(w, &http.Cookie{
					Name:     "mode",
					Value:    mode,
					SameSite: http.SameSiteLaxMode,
					Secure:   true,
					Path:     "/",
				})
			}
			handler.ServeHTTP(w, r)
		})
	}
}
