package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/caos/zitadel/internal/api"
)

type key int

const (
	nonceKey key = 0

	DefaultNonceLength = uint(32)
)

func SecurityHeaders(csp *CSP, nonceLength ...uint) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		if csp == nil {
			csp = &DefaultSCP
		}
		length := DefaultNonceLength
		if len(nonceLength) > 0 {
			length = nonceLength[0]
		}
		return &headers{
			csp:         csp,
			handler:     handler,
			nonceLength: length,
		}
	}
}

type headers struct {
	csp         *CSP
	handler     http.Handler
	nonceLength uint
}

func (h *headers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nonce := GetNonce(r)
	if nonce == "" {
		var err error
		nonce, err = generateNonce(h.nonceLength)
		if err != nil {

		}
		r = saveContext(r, nonceKey, nonce)
	}
	headers := w.Header()
	headers.Set(api.ContentSecurityPolicy, h.csp.Value(nonce))
	headers.Set(api.XXSSProtection, "1; mode=block")
	headers.Set(api.StrictTransportSecurity, "max-age=31536000; includeSubDomains")
	headers.Set(api.XFrameOptions, "DENY")
	headers.Set(api.XContentTypeOptions, "nosniff")
	headers.Set(api.ReferrerPolicy, "same-origin")
	headers.Set(api.FeaturePolicy, "payment 'none'")
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
