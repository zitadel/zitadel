package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/denylist"
)

type ClientConfig struct {
	MaxBodySize         int64
	Timeout             time.Duration
	MaxRedirects        int
	AllowHTTPSDowngrade bool
	DenyList            []denylist.AddressChecker
}

// NewClient returns a new http.Client with the configured settings.
// The client is protected against DNS rebinding attacks, redirects, HTTPs downgrades, and response body size limits.
func (c *ClientConfig) NewClient() *http.Client {
	baseTransport := denylist.NewHTTPTransport(c.DenyList)

	protectedTransport := &MaxBytesRoundTripper{
		Underlying: baseTransport,
		MaxBytes:   c.MaxBodySize,
	}

	return &http.Client{
		Transport: protectedTransport,
		Timeout:   c.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= c.MaxRedirects {
				return ErrTooManyRedirects
			}
			if !c.AllowHTTPSDowngrade && len(via) > 0 {
				prev := via[len(via)-1]
				if prev != nil && prev.URL != nil && req.URL != nil &&
					strings.EqualFold(prev.URL.Scheme, "https") &&
					!strings.EqualFold(req.URL.Scheme, "https") {
					return ErrHTTPsDowngrade
				}
			}
			return denylist.IsURLBlocked(c.DenyList, req.URL, nil)
		},
	}
}

// MergeDeprecatedDenylists merges the two deprecated (actions) denylists into the main denylist.
func (c *ClientConfig) MergeDeprecatedDenylists(actionsV1, actionsV2 []denylist.AddressChecker) {
	c.DenyList = append(c.DenyList, actionsV1...)
	c.DenyList = append(c.DenyList, actionsV2...)
}

var (
	// ErrResponseTooLarge is returned when the response body exceeds the configured limit.
	ErrResponseTooLarge = errors.New("response body exceeded maximum allowed size")
	//ErrTooManyRedirects is returned when the number of redirects exceeds the configured limit.
	ErrTooManyRedirects = errors.New("stopped after too many redirects")
	// ErrHTTPsDowngrade is returned when the client attempts to downgrade to HTTP.
	ErrHTTPsDowngrade = errors.New("redirect downgrade from https to http is not allowed")
)

// MaxBytesRoundTripper wraps an existing RoundTripper to protect against OOM.
type MaxBytesRoundTripper struct {
	Underlying http.RoundTripper
	MaxBytes   int64
}

func (m *MaxBytesRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := m.Underlying
	if transport == nil {
		transport = http.DefaultTransport
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.ContentLength > m.MaxBytes {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: Content-Length is %d (limit %d)", ErrResponseTooLarge, resp.ContentLength, m.MaxBytes)
	}

	resp.Body = &strictMaxBytesReader{
		limitReader: io.LimitReader(resp.Body, m.MaxBytes+1), // We initialize with +1 to detect overflows during Read
		closer:      resp.Body,
		limit:       m.MaxBytes,
	}

	return resp, nil
}

var _ http.RoundTripper = (*MaxBytesRoundTripper)(nil)

// strictMaxBytesReader enforces a hard limit and returns an explicit error if exceeded.
type strictMaxBytesReader struct {
	limitReader io.Reader
	closer      io.Closer
	limit       int64
	bytesRead   int64
}

func (s *strictMaxBytesReader) Read(p []byte) (int, error) {
	n, err := s.limitReader.Read(p)
	s.bytesRead += int64(n)

	if s.bytesRead > s.limit {
		// Because LimitReader stops at limit+1, excess is mathematically guaranteed to be 1
		safeN := n - 1
		if safeN < 0 {
			safeN = 0
		}

		// This instantly slices off that single offending byte from the view of the caller
		p = p[:safeN]
		return safeN, ErrResponseTooLarge
	}

	return n, err
}

func (s *strictMaxBytesReader) Close() error {
	return s.closer.Close()
}

var _ io.ReadCloser = (*strictMaxBytesReader)(nil)
