package middleware

import "net/http"

// HandlerFuncWithError is a http handler func which can return an error
// the error should then get handled later on in the pipeline by an error handler
// the error handler can be dependent on the interface standard (e.g. SCIM, Problem Details, ...)
type HandlerFuncWithError = func(w http.ResponseWriter, r *http.Request) error

// MiddlewareWithErrorFunc is a http middleware which can return an error
// the error should then get handled later on in the pipeline by an error handler
// the error handler can be dependent on the interface standard (e.g. SCIM, Problem Details, ...)
type MiddlewareWithErrorFunc = func(HandlerFuncWithError) HandlerFuncWithError

// ErrorHandlerFunc handles errors and returns a regular http handler
type ErrorHandlerFunc = func(HandlerFuncWithError) http.Handler

func ChainedWithErrorHandler(errorHandler ErrorHandlerFunc, middlewares ...MiddlewareWithErrorFunc) func(HandlerFuncWithError) http.Handler {
	return func(next HandlerFuncWithError) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return errorHandler(next)
	}
}
