package signals

import (
	"net/http"
	"strings"
	"testing"
)

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"US", true},
		{"us", true},
		{"De", true},
		{"ABC", true},
		{"", true}, // vacuously true
		{"U1", false},
		{"U ", false},
		{"U-S", false},
		{"US!", false},
		{"12", false},
		{"<script>", false},
		{strings.Repeat("A", 100), true},
		{"' OR 1=1 --", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isAlpha(tt.input); got != tt.want {
				t.Errorf("isAlpha(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hel"},
		{"hello", 0, ""},
		{"", 10, ""},
		{strings.Repeat("x", 1000), 512, strings.Repeat("x", 512)},
	}
	for _, tt := range tests {
		got := truncateString(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestStripPort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"192.168.1.1:8080", "192.168.1.1"},
		{"192.168.1.1", "192.168.1.1"},
		{"[::1]:443", "::1"},
		{"::1", "::1"},
		{"", ""},
		{"example.com:443", "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := stripPort(tt.input); got != tt.want {
				t.Errorf("stripPort(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseForwardedChain(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int // expected length, -1 for nil
	}{
		{"empty", "", -1},
		{"whitespace", "   ", -1},
		{"single", "192.168.1.1", 1},
		{"two hops", "10.0.0.1, 192.168.1.1", 2},
		{"three with spaces", " 10.0.0.1 , 172.16.0.1 , 192.168.1.1 ", 3},
		{"max hops", strings.TrimRight(strings.Repeat("10.0.0.1,", maxForwardedHops+5), ","), maxForwardedHops},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseForwardedChain(tt.input)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("parseForwardedChain(%q) len = %d, want %d", tt.input, len(got), tt.want)
			}
		})
	}
}

func TestExtractHTTPContext_CountryValidation(t *testing.T) {
	tests := []struct {
		name     string
		country  string
		wantKeep bool
	}{
		{"valid alpha-2", "US", true},
		{"valid alpha-3", "USA", true},
		{"lowercase (uppercased)", "de", true},
		{"empty", "", false},
		{"single char", "U", false},
		{"four chars", "ABCD", false},
		{"numeric", "12", false},
		{"injection attempt", "' OR 1=1 --", false},
		{"huge payload", strings.Repeat("A", 1000), false},
		{"html injection", "<b>US</b>", false},
		{"with space", "U S", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			h.Set("X-Country", tt.country)
			ctx := ExtractHTTPContext(h, "X-Country")
			if tt.wantKeep && ctx.Country == "" {
				t.Errorf("expected country to be preserved for %q, got empty", tt.country)
			}
			if !tt.wantKeep && ctx.Country != "" {
				t.Errorf("expected country to be rejected for %q, got %q", tt.country, ctx.Country)
			}
		})
	}
}

func TestExtractHTTPContext_NilHeaders(t *testing.T) {
	ctx := ExtractHTTPContext(nil, "X-Country")
	if ctx.Country != "" || ctx.AcceptLanguage != "" || ctx.Referer != "" {
		t.Error("expected empty context for nil headers")
	}
}

func TestExtractHTTPContext_Truncation(t *testing.T) {
	h := http.Header{}
	h.Set("Accept-Language", strings.Repeat("en-US,", 100))
	h.Set("Referer", strings.Repeat("https://example.com/", 200))
	h.Set("User-Agent", strings.Repeat("Mozilla/5.0 ", 100))

	ctx := ExtractHTTPContext(h, "")
	if len(ctx.AcceptLanguage) > maxAcceptLanguageLen {
		t.Errorf("AcceptLanguage not truncated: len=%d, max=%d", len(ctx.AcceptLanguage), maxAcceptLanguageLen)
	}
	if len(ctx.Referer) > maxRefererLen {
		t.Errorf("Referer not truncated: len=%d, max=%d", len(ctx.Referer), maxRefererLen)
	}
}

func TestExtractHTTPContext_HTTPS(t *testing.T) {
	h := http.Header{}
	h.Set("X-Forwarded-Proto", "https")
	ctx := ExtractHTTPContext(h, "")
	if !ctx.IsHTTPS {
		t.Error("expected IsHTTPS=true for X-Forwarded-Proto: https")
	}

	h.Set("X-Forwarded-Proto", "http")
	ctx = ExtractHTTPContext(h, "")
	if ctx.IsHTTPS {
		t.Error("expected IsHTTPS=false for X-Forwarded-Proto: http")
	}
}
