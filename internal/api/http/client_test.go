package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/denylist"
)

func TestClientConfig_NewClient_Allowed(t *testing.T) {
	t.Parallel()

	// Spin up a simple local test server to serve a dummy response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	}))
	defer server.Close()

	// An empty denylist permits all safe destinations
	cfg := &ClientConfig{
		MaxBodySize:  1024,
		Timeout:      2 * time.Second,
		MaxRedirects: 3,
		DenyList:     []denylist.AddressChecker{},
	}

	client := cfg.NewClient()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	if resp != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "hello world", string(body))
	}
}

func TestClientConfig_Redirect_HTTPSDowngradeBlocked(t *testing.T) {
	t.Parallel()

	cfg := &ClientConfig{
		MaxBodySize:         1024,
		Timeout:             2 * time.Second,
		MaxRedirects:        3,
		AllowHTTPSDowngrade: false, // Explicitly disallow downgrades
		DenyList:            []denylist.AddressChecker{},
	}

	client := cfg.NewClient()

	// Mock a historical request chain moving from a secure schema to an unencrypted target
	viaReq, _ := http.NewRequest(http.MethodGet, "https://secure-identity.com/oauth", nil)
	currentReq, _ := http.NewRequest(http.MethodGet, "http://insecure-identity.com/callback", nil)

	err := client.CheckRedirect(currentReq, []*http.Request{viaReq})
	assert.ErrorIs(t, err, ErrHTTPsDowngrade)
}

func TestMaxBytesRoundTripper_EnforcesCap(t *testing.T) {
	t.Parallel()

	// Spin up a local server delivering an exact 10-byte payload response string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("0123456789"))
	}))
	defer server.Close()

	tests := []struct {
		name        string
		maxBytes    int64
		expectError error
	}{
		{
			name:        "body cleanly within threshold limit",
			maxBytes:    15,
			expectError: nil,
		},
		{
			name:        "body matches boundary size exactly",
			maxBytes:    10,
			expectError: nil,
		},
		{
			name:        "body exceeds allowable limit",
			maxBytes:    5,
			expectError: ErrResponseTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ClientConfig{
				MaxBodySize:  tt.maxBytes,
				Timeout:      2 * time.Second,
				MaxRedirects: 3,
				DenyList:     []denylist.AddressChecker{},
			}

			client := cfg.NewClient()
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
			resp, err := client.Do(req)

			if tt.expectError != nil {
				// Depending on optimization, errors can trip directly on call or during chunk reading
				if err != nil {
					assert.ErrorIs(t, err, tt.expectError)
					return
				}
				defer resp.Body.Close()
				_, readErr := io.ReadAll(resp.Body)
				assert.ErrorIs(t, readErr, tt.expectError)
			} else {
				assert.NoError(t, err)
				defer resp.Body.Close()
				body, readErr := io.ReadAll(resp.Body)
				assert.NoError(t, readErr)
				assert.Equal(t, 10, len(body))
			}
		})
	}
}

func TestClientConfig_MergeDeprecatedDenylists(t *testing.T) {
	t.Parallel()

	cfg := &ClientConfig{DenyList: []denylist.AddressChecker{denylist.NewHostChecker("1.1.1.1")}}

	v1 := []denylist.AddressChecker{denylist.NewHostChecker("2.2.2.2")}
	v2 := []denylist.AddressChecker{denylist.NewHostChecker("3.3.3.3")}

	cfg.MergeDeprecatedDenylists(v1, v2)

	assert.Len(t, cfg.DenyList, 3)
}
