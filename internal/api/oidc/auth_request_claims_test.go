package oidc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_authRequestClaimsHandler(t *testing.T) {
	const claims = `{"userinfo":{"urn:swissid:qor":{"value":"qor2"}}}`
	tests := []struct {
		name string
		req  *http.Request
		want string
	}{
		{
			name: "other path",
			req:  httptest.NewRequest(http.MethodGet, "/oauth/v2/token?claims="+url.QueryEscape(claims), nil),
			want: "",
		},
		{
			name: "no claims",
			req:  httptest.NewRequest(http.MethodGet, "/oauth/v2/authorize?client_id=client", nil),
			want: "",
		},
		{
			name: "invalid json",
			req:  httptest.NewRequest(http.MethodGet, "/oauth/v2/authorize?claims=%7Binvalid", nil),
			want: "",
		},
		{
			name: "too long",
			req:  httptest.NewRequest(http.MethodGet, "/oauth/v2/authorize?claims="+url.QueryEscape(`"`+strings.Repeat("a", maxAuthRequestClaimsLength)+`"`), nil),
			want: "",
		},
		{
			name: "query",
			req:  httptest.NewRequest(http.MethodGet, "/oauth/v2/authorize?claims="+url.QueryEscape(claims), nil),
			want: claims,
		},
		{
			name: "form post",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/oauth/v2/authorize", strings.NewReader(url.Values{"claims": {claims}}.Encode()))
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r
			}(),
			want: claims,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			handler := authRequestClaimsHandler("/oauth/v2/authorize")(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				got = authRequestClaimsFromContext(r.Context())
			}))
			handler.ServeHTTP(httptest.NewRecorder(), tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}
