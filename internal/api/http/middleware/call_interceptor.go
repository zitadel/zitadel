package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/api/call"
)

func DefaultCallHandler(handler http.Handler) http.Handler {
	return CallHandler(handler)
}

func CallHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(call.WithTimestamp(r.Context())))
	})
}
