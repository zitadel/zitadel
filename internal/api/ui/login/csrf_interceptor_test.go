package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/csrf"

	"github.com/zitadel/zitadel/internal/api/authz"
)

// Test_createCSRFInterceptor_origin verifies the origin validation gorilla/csrf
// performs since v1.7.3: requests are assumed to be HTTPS unless the request is
// marked as plaintext, which the interceptor must do when ZITADEL is not
// exposed with TLS (ExternalSecure=false). Without the marker every plain-HTTP
// POST would be rejected with "origin invalid".
func Test_createCSRFInterceptor_origin(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	tests := []struct {
		name           string
		externalSecure bool
		origin         string
		wantStatus     int
	}{
		{name: "https, same origin", externalSecure: true, origin: "https://example.com", wantStatus: http.StatusOK},
		{name: "https, http origin rejected", externalSecure: true, origin: "http://example.com", wantStatus: http.StatusForbidden},
		{name: "https, foreign origin rejected", externalSecure: true, origin: "https://evil.example", wantStatus: http.StatusForbidden},
		{name: "plaintext, same origin", externalSecure: false, origin: "http://example.com", wantStatus: http.StatusOK},
		{name: "plaintext, foreign origin rejected", externalSecure: false, origin: "http://evil.example", wantStatus: http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			handler := createCSRFInterceptor("csrf", key, tt.externalSecure, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token = csrf.Token(r)
				w.WriteHeader(http.StatusOK)
			}))
			ctx := authz.NewMockContext("instance1", "org1", "user1")

			// first request obtains the cookie and a token
			get := httptest.NewRequest(http.MethodGet, "http://example.com/login", nil).WithContext(ctx)
			getRec := httptest.NewRecorder()
			handler.ServeHTTP(getRec, get)
			if getRec.Code != http.StatusOK {
				t.Fatalf("GET status = %d, want 200", getRec.Code)
			}
			cookies := getRec.Result().Cookies()
			if len(cookies) == 0 || token == "" {
				t.Fatal("expected csrf cookie and token from GET")
			}

			post := httptest.NewRequest(http.MethodPost, "http://example.com/login", nil).WithContext(ctx)
			for _, c := range cookies {
				post.AddCookie(c)
			}
			post.Header.Set("X-CSRF-Token", token)
			post.Header.Set("Origin", tt.origin)
			postRec := httptest.NewRecorder()
			handler.ServeHTTP(postRec, post)
			if postRec.Code != tt.wantStatus {
				t.Errorf("POST status = %d, want %d", postRec.Code, tt.wantStatus)
			}
		})
	}
}
