package risk

import (
	"net/http"
	"strings"
)

// ExtractHTTPContext enriches a Signal with data extracted from HTTP headers.
// It reads from the provided header map (typically from domain.UserAgent.Header
// or copied from the gRPC gateway context) and the risk config for the
// geo-country header name.
//
// The caller is responsible for merging the returned fields into the Signal.
func ExtractHTTPContext(headers http.Header, geoCountryHeader string) HTTPContext {
	if headers == nil {
		return HTTPContext{}
	}
	ctx := HTTPContext{
		AcceptLanguage: headers.Get("Accept-Language"),
		Referer:        headers.Get("Referer"),
		SecFetchSite:   headers.Get("Sec-Fetch-Site"),
		IsHTTPS:        strings.EqualFold(headers.Get("X-Forwarded-Proto"), "https"),
		ForwardedChain: parseForwardedChain(headers.Get("X-Forwarded-For")),
	}
	if geoCountryHeader != "" {
		ctx.Country = strings.ToUpper(strings.TrimSpace(headers.Get(geoCountryHeader)))
	}
	return ctx
}

// HTTPContext holds HTTP-derived data extracted from request headers.
type HTTPContext struct {
	AcceptLanguage string
	Country        string
	ForwardedChain []string
	Referer        string
	SecFetchSite   string
	IsHTTPS        bool
}

// ApplyTo merges the extracted HTTP context into a Signal, only setting fields
// that are not already populated.
func (h HTTPContext) ApplyTo(s *Signal) {
	if s.AcceptLanguage == "" {
		s.AcceptLanguage = h.AcceptLanguage
	}
	if s.Country == "" {
		s.Country = h.Country
	}
	if len(s.ForwardedChain) == 0 {
		s.ForwardedChain = h.ForwardedChain
	}
	if s.Referer == "" {
		s.Referer = h.Referer
	}
	if s.SecFetchSite == "" {
		s.SecFetchSite = h.SecFetchSite
	}
	if !s.IsHTTPS {
		s.IsHTTPS = h.IsHTTPS
	}
}

// parseForwardedChain splits the X-Forwarded-For header value into individual
// IP addresses, trimming whitespace. Returns nil for empty input.
func parseForwardedChain(xff string) []string {
	xff = strings.TrimSpace(xff)
	if xff == "" {
		return nil
	}
	parts := strings.Split(xff, ",")
	chain := make([]string, 0, len(parts))
	for _, part := range parts {
		ip := strings.TrimSpace(part)
		if ip != "" {
			chain = append(chain, ip)
		}
	}
	if len(chain) == 0 {
		return nil
	}
	return chain
}
