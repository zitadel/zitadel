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

// HTTPContext holds HTTP-derived data extracted from request headers.
type HTTPContext struct {
	AcceptLanguage string
	Country        string
	ForwardedChain []string
	Referer        string
	SecFetchSite   string
	IsHTTPS        bool
}

// ExtractHTTPContext extracts signal-relevant data from HTTP headers.
// Header values are truncated to prevent oversized payloads.
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
		country := strings.ToUpper(strings.TrimSpace(headers.Get(geoCountryHeader)))
		// Validate ISO 3166-1 alpha-2 format (2 uppercase letters) and
		// truncate to prevent storage inflation from user-controlled headers.
		if len(country) >= 2 && len(country) <= 3 && isAlpha(country) {
			ctx.Country = country
		}
	}
	return ctx
}

// isAlpha returns true if s contains only ASCII letters.
func isAlpha(s string) bool {
	for i := 0; i < len(s); i++ {
		if (s[i] < 'A' || s[i] > 'Z') && (s[i] < 'a' || s[i] > 'z') {
			return false
		}
	}
	return true
}

// parseForwardedChain splits X-Forwarded-For into individual IPs,
// capped at maxForwardedHops entries.
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
