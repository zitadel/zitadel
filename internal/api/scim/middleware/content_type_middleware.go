package middleware

import (
	"mime"
	"net/http"
	"strings"

	"github.com/zitadel/logging"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ContentTypeScim                = "application/scim+json"
	ContentTypeJson                = "application/json"
	ContentTypeApplicationWildcard = "application/*"
	ContentTypeWildcard            = "*/*"
)

func ContentTypeMiddleware(next middleware.HandlerFuncWithError) middleware.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set(zhttp.ContentType, ContentTypeScim)

		if !validateContentType(r.Header.Get(zhttp.ContentType)) {
			return zerrors.ThrowInvalidArgumentf(nil, "SMCM-12x4", "Invalid content type header")
		}

		if !validateContentType(r.Header.Get(zhttp.Accept)) {
			return zerrors.ThrowInvalidArgumentf(nil, "SMCM-12x5", "Invalid accept header")
		}

		return next(w, r)
	}
}

func validateContentType(contentType string) bool {
	if contentType == "" {
		return true
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		logging.OnError(err).Warn("failed to parse content type header")
		return false
	}

	if mediaType != "" &&
		!strings.EqualFold(mediaType, ContentTypeWildcard) &&
		!strings.EqualFold(mediaType, ContentTypeApplicationWildcard) &&
		!strings.EqualFold(mediaType, ContentTypeJson) &&
		!strings.EqualFold(mediaType, ContentTypeScim) {
		return false
	}

	charset, ok := params["charset"]
	return !ok || strings.EqualFold(charset, "utf-8")
}
