package client

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// UserAgent is the global User-Agent string sent with all API requests.
// It is mutated by main.go at startup to include version and OS/Arch.
var UserAgent = "zitadel-cli/dev"

// EnableDebug toggles verbose request/response logging to stderr.
var EnableDebug = false

// New builds an http.Client that injects the Bearer token from the given TokenSource on every request.
func New(tokenSource oauth2.TokenSource) *http.Client {
	var transport http.RoundTripper = &authTransport{
		base:   http.DefaultTransport,
		source: tokenSource,
	}

	if EnableDebug {
		transport = &loggingTransport{base: transport}
	}

	return &http.Client{
		Transport: transport,
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
	r.Header.Set("User-Agent", UserAgent)
	return t.base.RoundTrip(r)
}

type loggingTransport struct {
	base http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Clone request body for logging
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	slog.Debug("Sending request",
		"method", req.Method,
		"url", req.URL.String(),
		"headers", req.Header,
		"body", string(reqBody),
	)

	resp, err := t.base.RoundTrip(req)

	if err != nil {
		slog.Error("Request failed",
			"error", err,
			"duration", time.Since(start),
		)
		return nil, err
	}

	// Clone response body for logging
	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	slog.Debug("Received response",
		"status", resp.StatusCode,
		"headers", resp.Header,
		"body", string(respBody),
		"duration", time.Since(start),
	)

	return resp, nil
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
