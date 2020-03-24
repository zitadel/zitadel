package middleware

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/caos/zitadel/internal/api"
)

var (
	DefaultCORSOptions = cors.Options{
		AllowCredentials: true,
		AllowedHeaders: []string{
			api.Origin,
			api.ContentType,
			api.Accept,
			api.AcceptLanguage,
			api.Authorization,
			api.ZitadelOrgID,
			"x-grpc-web", //TODO: needed
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
			api.Location,
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
