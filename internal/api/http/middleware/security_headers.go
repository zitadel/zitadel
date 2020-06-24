package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	http2 "github.com/caos/zitadel/internal/api/http"
)

type key int

const (
	nonceKey key = 0

	DefaultNonceLength = uint(32)
)

func SecurityHeaders(csp *CSP, errorHandler func(error) http.Handler, nonceLength ...uint) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		if csp == nil {
			csp = &DefaultSCP
		}
		length := DefaultNonceLength
		if len(nonceLength) > 0 {
			length = nonceLength[0]
		}
		return &headers{
			csp:          csp,
			handler:      handler,
			errorHandler: errorHandler,
			nonceLength:  length,
		}
	}
}

type headers struct {
	csp          *CSP
	handler      http.Handler
	errorHandler func(err error) http.Handler
	nonceLength  uint
}

func (h *headers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nonce := GetNonce(r)
	if nonce == "" {
		var err error
		nonce, err = generateNonce(h.nonceLength)
		if err != nil {
			errorHandler := h.errorHandler
			if errorHandler == nil {
				errorHandler = func(err error) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					})
				}
			}
			errorHandler(err).ServeHTTP(w, r)
			return
		}
		r = saveContext(r, nonceKey, nonce)
	}
	headers := w.Header()
	headers.Set(http2.ContentSecurityPolicy, h.csp.Value(nonce))
	headers.Set(http2.XXSSProtection, "1; mode=block")
	headers.Set(http2.StrictTransportSecurity, "max-age=31536000; includeSubDomains")
	headers.Set(http2.XFrameOptions, "DENY")
	headers.Set(http2.XContentTypeOptions, "nosniff")
	headers.Set(http2.ReferrerPolicy, "same-origin")
	headers.Set(http2.FeaturePolicy, "payment 'none'")
	//PLANNED: add expect-ct

	h.handler.ServeHTTP(w, r)
}

func GetNonce(r *http.Request) string {
	nonce, _ := getContext(r, nonceKey).(string)
	return nonce
}

func generateNonce(length uint) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func saveContext(r *http.Request, key, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

func getContext(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}
