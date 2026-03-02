package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zitadel/sloggcp"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// RecoverHandler recovers from panics in the HTTP handler chain
// and calls the provided writeResponse function to write an appropriate response to the client.
//
// The request context is canceled with the panic error as the cause.
func RecoverHandler(writeResponse func(w http.ResponseWriter, r *http.Request, err error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancelCause(r.Context())
			r = r.WithContext(ctx)

			defer func() {
				var err error
				if rec := recover(); rec != nil {
					recErr, ok := rec.(error)
					if !ok {
						recErr = fmt.Errorf("%v", rec)
					}
					err = zerrors.ThrowInternal(recErr, zerrors.IDRecover, "Errors.Internal")
					logRecovered(ctx, err)
					writeResponse(w, r, err)
				}
				cancel(err)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// FallbackRecoverHandler recovers from panics in the HTTP handler chain
// and returns a 500 Internal Server Error response.
// The request context is canceled with the panic error as the cause,
// so that any ongoing operations can be stopped and cleaned up.
//
// The response is sent as a text/plain response.
// It is used as a last line of defense to prevent the server from crashing
// due to panics in the handlers.
// Protocols (OIDC, SAML, HTML etc.) should use [RecoverHandler] to write
// properly formatted error responses to the clients.
func FallbackRecoverHandler() func(http.Handler) http.Handler {
	return RecoverHandler(func(w http.ResponseWriter, r *http.Request, _ error) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	})
}

// RecoverHandlerWithError is similar to [RecoverHandler] but returns an error instead of writing a response directly.
func RecoverHandlerWithError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		ctx, cancel := context.WithCancelCause(r.Context())
		r = r.WithContext(ctx)
		defer func() {
			if rec := recover(); rec != nil {
				recErr, ok := rec.(error)
				if !ok {
					recErr = fmt.Errorf("%v", rec)
				}
				err = zerrors.ThrowInternal(recErr, zerrors.IDRecover, "Errors.Internal")
				logRecovered(ctx, err)
			}
			cancel(err)
		}()
		return next(w, r)
	}
}

func logRecovered(ctx context.Context, err error) {
	logger := logging.FromCtx(ctx)
	logger.Log(ctx, sloggcp.LevelAlert, "recovered from panic", "err", err)
}
