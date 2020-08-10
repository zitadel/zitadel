package middleware

import (
	"net/http"

	"github.com/rs/cors"

	http_utils "github.com/caos/zitadel/internal/api/http"
)

var (
	DefaultCORSOptions = cors.Options{
		AllowCredentials: true,
		AllowedHeaders: []string{
			http_utils.Origin,
			http_utils.ContentType,
			http_utils.Accept,
			http_utils.AcceptLanguage,
			http_utils.Authorization,
			http_utils.ZitadelOrgID,
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
			http_utils.Location,
		},
		AllowOriginFunc: func(_ string) bool {
			return true
		},
	}
)

func CORSInterceptorOpts(opts cors.Options, h http.Handler) http.Handler {
	return cors.New(opts).Handler(h)
}

func CORSInterceptor(h http.Handler) http.Handler {
	return CORSInterceptorOpts(DefaultCORSOptions, h)
}
