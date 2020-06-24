package middleware

import (
	"net/http"

	"github.com/rs/cors"

	http2 "github.com/caos/zitadel/internal/api/http"
)

var (
	DefaultCORSOptions = cors.Options{
		AllowCredentials: true,
		AllowedHeaders: []string{
			http2.Origin,
			http2.ContentType,
			http2.Accept,
			http2.AcceptLanguage,
			http2.Authorization,
			http2.ZitadelOrgID,
		},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		ExposedHeaders: []string{
			http2.Location,
		},
		AllowedOrigins: []string{
			"http://localhost:*",
		},
	}
)

func CORSInterceptorOpts(opts cors.Options, h http.Handler) http.Handler {
	return cors.New(opts).Handler(h)
}

func CORSInterceptor(h http.Handler) http.Handler {
	return CORSInterceptorOpts(DefaultCORSOptions, h)
}
