package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

func LogHandler(service string, ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return logging.NewHandler(h, service, ignoredPrefix...)
	}
}
