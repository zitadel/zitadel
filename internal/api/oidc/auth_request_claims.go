package oidc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zitadel/logging"
)

// maxAuthRequestClaimsLength limits the size of the `claims` parameter stored on the auth request.
const maxAuthRequestClaimsLength = 4096

type authRequestClaimsKey struct{}

// authRequestClaimsHandler stores the `claims` parameter of an authorization request in the context,
// since it is not part of the parsed [oidc.AuthRequest].
//
// POC: the claims are only stored to be forwarded to an IdP, see [command.Commands.AuthFromProvider].
func authRequestClaimsHandler(authorizationPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != authorizationPath {
				next.ServeHTTP(w, r)
				return
			}
			if err := r.ParseForm(); err != nil {
				next.ServeHTTP(w, r)
				return
			}
			claims := r.Form.Get("claims")
			if claims == "" {
				next.ServeHTTP(w, r)
				return
			}
			if len(claims) > maxAuthRequestClaimsLength || !json.Valid([]byte(claims)) {
				logging.WithFields("length", len(claims)).Info("ignoring invalid claims parameter of auth request")
				next.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authRequestClaimsKey{}, claims)))
		})
	}
}

func authRequestClaimsFromContext(ctx context.Context) string {
	claims, _ := ctx.Value(authRequestClaimsKey{}).(string)
	return claims
}
