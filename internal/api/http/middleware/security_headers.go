package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

type securityKey int

const (
	nonceKey securityKey = 0

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
					})
				}
			}
			errorHandler(err).ServeHTTP(w, r)
			return
		}
		r = saveContext(r, nonceKey, nonce)
	}
	allowedHosts := authz.GetInstance(r.Context()).SecurityPolicyAllowedOrigins()
	headers := w.Header()
	headers.Set(http_utils.ContentSecurityPolicy, h.csp.Value(nonce, r.Host, allowedHosts))
	headers.Set(http_utils.XXSSProtection, "1; mode=block")
	headers.Set(http_utils.StrictTransportSecurity, "max-age=31536000; includeSubDomains")
	if len(allowedHosts) == 0 {
		headers.Set(http_utils.XFrameOptions, "DENY")
	}
	headers.Set(http_utils.XContentTypeOptions, "nosniff")
	headers.Set(http_utils.ReferrerPolicy, "same-origin")
	headers.Set(http_utils.FeaturePolicy, "payment 'none'")
	headers.Set(http_utils.PermissionsPolicy, "payment=()")
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

func saveContext(r *http.Request, key, value any) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

func getContext(r *http.Request, key any) any {
	return r.Context().Value(key)
}
