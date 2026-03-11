package signals

import (
	"net/http"
	"strings"
)

const (
	maxAcceptLanguageLen = 256
	maxRefererLen        = 2048
	maxUserAgentLen      = 512
	maxForwardedHops     = 32
)

// truncateString returns s truncated to maxLen bytes.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// ExtractHTTPContext enriches a Signal with data extracted from HTTP headers.
// It reads from the provided header map (typically from domain.UserAgent.Header
// or copied from the gRPC gateway context) and the risk config for the
// geo-country header name.
//
// Header values are truncated to prevent oversized payloads from reaching
// storage. The X-Forwarded-For chain is capped at maxForwardedHops entries.
func ExtractHTTPContext(headers http.Header, geoCountryHeader string) HTTPContext {
	if headers == nil {
		return HTTPContext{}
	}
	ctx := HTTPContext{
		AcceptLanguage: truncateString(headers.Get("Accept-Language"), maxAcceptLanguageLen),
		Referer:        truncateString(headers.Get("Referer"), maxRefererLen),
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
// IP addresses, trimming whitespace. The result is capped at maxForwardedHops
// entries to prevent abuse via oversized headers. Returns nil for empty input.
func parseForwardedChain(xff string) []string {
	xff = strings.TrimSpace(xff)
	if xff == "" {
		return nil
	}
	parts := strings.Split(xff, ",")
	cap := len(parts)
	if cap > maxForwardedHops {
		cap = maxForwardedHops
	}
	chain := make([]string, 0, cap)
	for _, part := range parts {
		ip := strings.TrimSpace(part)
		if ip != "" {
			chain = append(chain, ip)
			if len(chain) >= maxForwardedHops {
				break
			}
		}
	}
	if len(chain) == 0 {
		return nil
	}
	return chain
}
