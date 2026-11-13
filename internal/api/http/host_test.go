package http

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		want    string
		wantErr bool
	}{
		{name: "hostname", host: "example.com", want: "example.com"},
		{name: "hostname with port", host: "example.com:8080", want: "example.com:8080"},
		{name: "ipv4", host: "192.168.1.1", want: "192.168.1.1"},
		{name: "ipv4 with port", host: "192.168.1.1:8080", want: "192.168.1.1:8080"},
		{name: "ipv6", host: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]", want: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]"},
		{name: "ipv6 with port", host: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080", want: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:8080"},
		{name: "hostname with underscore", host: "my_host.internal:8080", want: "my_host.internal:8080"},
		{name: "fully qualified hostname with trailing dot", host: "example.com.", want: "example.com."},
		{name: "punycode hostname", host: "xn--bcher-kva.example", want: "xn--bcher-kva.example"},
		{name: "empty port is removed", host: "example.com:", want: "example.com"},
		{name: "empty port of ipv6 is removed", host: "[::1]:", want: "[::1]"},
		{name: "leading zeros of port are removed", host: "example.com:0443", want: "example.com:443"},
		// characters allowed in a Host header, which no instance or trusted domain contains,
		// so these hosts are rejected by the instance lookup or the trusted domain check
		{name: "sub-delimiter semicolon", host: "trusted.example.com;evil", want: "trusted.example.com;evil"},
		{name: "sub-delimiter comma", host: "a,b.example.com", want: "a,b.example.com"},
		{name: "sub-delimiter parentheses", host: "(trusted)", want: "(trusted)"},
		{name: "wildcard", host: "*.example.com", want: "*.example.com"},
		{name: "tilde", host: "~evil", want: "~evil"},

		{name: "empty", host: "", wantErr: true},
		{name: "unicode hostname", host: "bücher.example", wantErr: true},
		{name: "non-numeric port", host: "example.com:port", wantErr: true},
		{name: "port zero", host: "example.com:0", wantErr: true},
		{name: "port out of range", host: "example.com:65536", wantErr: true},
		{name: "port only", host: ":8080", wantErr: true},
		{name: "signed port", host: "example.com:+80", wantErr: true},
		{name: "unclosed ipv6 bracket", host: "[::1", wantErr: true},
		{name: "characters after ipv6 bracket", host: "[::1]x:8080", wantErr: true},
		{name: "userinfo without user", host: "localhost:@attacker.com", wantErr: true},
		{name: "userinfo without user and with port", host: "localhost:@attacker.com:8080", wantErr: true},
		{name: "userinfo", host: "a:b@evil.com", wantErr: true},
		{name: "userinfo with trusted host and port", host: "trusted.example.com:443@evil.example", wantErr: true},
		{name: "userinfo with trusted host", host: "trusted.example.com@evil.example", wantErr: true},
		{name: "backslash", host: `trusted.example.com\@evil.example`, wantErr: true},
		{name: "fragment", host: "trusted.example.com#@evil.example", wantErr: true},
		{name: "query", host: "trusted.example.com?@evil.example", wantErr: true},
		{name: "path", host: "trusted.example.com/@evil.example", wantErr: true},
		{name: "path after port", host: "evil.com:443/", wantErr: true},
		{name: "whitespace", host: "trusted.example.com @evil.example", wantErr: true},
		{name: "percent encoding", host: "trusted.example.com%40evil.example", wantErr: true},
		{name: "ipv6 without brackets", host: "::1", wantErr: true},
		{name: "ipv6 with zone", host: "[fe80::1%25en0]:8080", wantErr: true},
		{name: "bracket without ipv6", host: "]", wantErr: true},
		{name: "bracket inside hostname", host: "a]b", wantErr: true},
		{name: "invalid ipv6", host: "[::1::2]", wantErr: true},
		{name: "brackets around hostname", host: "[example.com]:8080", wantErr: true},
		{name: "brackets around ipv4", host: "[1.2.3.4]", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateHost(tt.host)
			if tt.wantErr {
				assert.True(t, zerrors.IsErrorInvalidArgument(err), "expected invalid argument error, got: %v", err)
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// FuzzValidateHost ensures that every accepted host is read the same way by the trusted domain check,
// by the URLs (e.g. password reset links) built from it and by the URL parsers of browsers and mail clients.
func FuzzValidateHost(f *testing.F) {
	for _, seed := range []string{
		"example.com",
		"example.com:8080",
		"example.com:",
		"example.com:0443",
		"[::1]:8080",
		"192.168.1.1:443",
		"my_host.internal",
		"trusted.example.com:443@evil.example",
		"trusted.example.com@evil.example",
		"a:b@evil.com",
		"localhost:@attacker.com",
		"evil.com:443/",
		`trusted.example.com\@evil.example`,
		"trusted.example.com%40evil.example",
		"trusted.example.com;evil",
		"::1",
		"]",
		"[fe80::1%25en0]:8080",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, rawHost string) {
		host, err := ValidateHost(rawHost)
		if err != nil {
			return
		}
		// URL parsers of browsers and mail clients (WHATWG) end the host at these characters, take the part before "@" as userinfo,
		// or strip whitespace and control characters, so an accepted host must not contain any of them.
		assert.False(t, strings.ContainsFunc(host, func(r rune) bool {
			return r <= ' ' || r >= 0x7f || strings.ContainsRune(`@/\?#%`, r)
		}), "accepted host %q contains a delimiter or whitespace", host)

		// the canonical form is stable
		again, err := ValidateHost(host)
		require.NoError(t, err, "canonical host %q is not valid", host)
		assert.Equal(t, host, again)

		domainCtx := NewDomainCtx(host, host, "https")
		// the host checked against the trusted domains ...
		checked := domainCtx.RequestedDomain()
		// ... must be the host a link built from the origin points to
		link, err := url.Parse(domainCtx.Origin() + "/ui/v2/login/password/set")
		require.NoError(t, err, "accepted host %q does not build a valid link", host)
		assert.Nil(t, link.User, "accepted host %q contains userinfo", host)
		// an IPv6 address is checked with brackets if no port is given ("[::1]"), url.URL.Hostname() returns it without
		checkedHostname := checked
		if strings.HasPrefix(checked, "[") && strings.HasSuffix(checked, "]") {
			checkedHostname = checked[1 : len(checked)-1]
		}
		assert.Equal(t, checkedHostname, link.Hostname(), "accepted host %q is checked as %q but links point to %q", host, checked, link.Hostname())
	})
}
