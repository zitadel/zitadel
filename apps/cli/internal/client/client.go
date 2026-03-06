package client

import (
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// New builds an http.Client that injects the Bearer token from the given TokenSource on every request.
func New(tokenSource oauth2.TokenSource) *http.Client {
	return &http.Client{
		Transport: &authTransport{
			base:   http.DefaultTransport,
			source: tokenSource,
		},
	}
}

type authTransport struct {
	base   http.RoundTripper
	source oauth2.TokenSource
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.source.Token()
	if err != nil {
		return nil, err
	}
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", token.Type()+" "+token.AccessToken)
	return t.base.RoundTrip(r)
}

// InstanceURL normalises a raw instance string to a base URL.
//
//	"mycompany.zitadel.cloud" → "https://mycompany.zitadel.cloud"
//	"https://auth.internal"   → "https://auth.internal"
func InstanceURL(raw string) string {
	raw = strings.TrimRight(raw, "/")
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "https://" + raw
	}
	return raw
}
