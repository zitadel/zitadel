package risk

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHTTPContext(t *testing.T) {
	tests := []struct {
		name             string
		headers          http.Header
		geoCountryHeader string
		want             HTTPContext
	}{
		{
			name:    "nil headers",
			headers: nil,
			want:    HTTPContext{},
		},
		{
			name:    "empty headers",
			headers: http.Header{},
			want:    HTTPContext{},
		},
		{
			name: "all fields populated",
			headers: http.Header{
				"Accept-Language":   {"en-US,en;q=0.9,de;q=0.8"},
				"Referer":          {"https://login.example.com/login"},
				"Sec-Fetch-Site":   {"same-origin"},
				"X-Forwarded-Proto": {"https"},
				"X-Forwarded-For":  {"1.2.3.4, 10.0.0.1, 192.168.1.1"},
				"Cf-Ipcountry":    {"CH"},
			},
			geoCountryHeader: "CF-IPCountry",
			want: HTTPContext{
				AcceptLanguage: "en-US,en;q=0.9,de;q=0.8",
				Referer:        "https://login.example.com/login",
				SecFetchSite:   "same-origin",
				IsHTTPS:        true,
				ForwardedChain: []string{"1.2.3.4", "10.0.0.1", "192.168.1.1"},
				Country:        "CH",
			},
		},
		{
			name: "no geo header configured",
			headers: http.Header{
				"Cf-Ipcountry": {"US"},
			},
			geoCountryHeader: "",
			want:             HTTPContext{},
		},
		{
			name: "HTTP not HTTPS",
			headers: http.Header{
				"X-Forwarded-Proto": {"http"},
			},
			want: HTTPContext{IsHTTPS: false},
		},
		{
			name: "country lowercased and trimmed",
			headers: http.Header{
				"X-Country": {"  de  "},
			},
			geoCountryHeader: "X-Country",
			want:             HTTPContext{Country: "DE"},
		},
		{
			name: "single forwarded IP",
			headers: http.Header{
				"X-Forwarded-For": {"1.2.3.4"},
			},
			want: HTTPContext{ForwardedChain: []string{"1.2.3.4"}},
		},
		{
			name: "empty forwarded for",
			headers: http.Header{
				"X-Forwarded-For": {""},
			},
			want: HTTPContext{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractHTTPContext(tt.headers, tt.geoCountryHeader)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHTTPContext_ApplyTo(t *testing.T) {
	t.Run("fills empty signal fields", func(t *testing.T) {
		signal := Signal{IP: "1.2.3.4"}
		ctx := HTTPContext{
			AcceptLanguage: "en-US",
			Country:        "CH",
			ForwardedChain: []string{"1.2.3.4", "10.0.0.1"},
			Referer:        "https://example.com",
			SecFetchSite:   "same-origin",
			IsHTTPS:        true,
		}
		ctx.ApplyTo(&signal)
		assert.Equal(t, "en-US", signal.AcceptLanguage)
		assert.Equal(t, "CH", signal.Country)
		assert.Equal(t, []string{"1.2.3.4", "10.0.0.1"}, signal.ForwardedChain)
		assert.Equal(t, "https://example.com", signal.Referer)
		assert.Equal(t, "same-origin", signal.SecFetchSite)
		assert.True(t, signal.IsHTTPS)
	})

	t.Run("does not overwrite existing signal fields", func(t *testing.T) {
		signal := Signal{
			AcceptLanguage: "de-DE",
			Country:        "DE",
			ForwardedChain: []string{"5.6.7.8"},
			Referer:        "https://original.com",
			SecFetchSite:   "cross-site",
			IsHTTPS:        true,
		}
		ctx := HTTPContext{
			AcceptLanguage: "en-US",
			Country:        "US",
			ForwardedChain: []string{"1.2.3.4"},
			Referer:        "https://new.com",
			SecFetchSite:   "none",
			IsHTTPS:        true,
		}
		ctx.ApplyTo(&signal)
		assert.Equal(t, "de-DE", signal.AcceptLanguage)
		assert.Equal(t, "DE", signal.Country)
		assert.Equal(t, []string{"5.6.7.8"}, signal.ForwardedChain)
		assert.Equal(t, "https://original.com", signal.Referer)
		assert.Equal(t, "cross-site", signal.SecFetchSite)
	})
}

func TestParseForwardedChain(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"  ", nil},
		{"1.2.3.4", []string{"1.2.3.4"}},
		{"1.2.3.4, 10.0.0.1", []string{"1.2.3.4", "10.0.0.1"}},
		{"1.2.3.4,10.0.0.1,192.168.0.1", []string{"1.2.3.4", "10.0.0.1", "192.168.0.1"}},
		{" 1.2.3.4 , 10.0.0.1 ", []string{"1.2.3.4", "10.0.0.1"}},
		{",,,", nil},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseForwardedChain(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
