package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileRules(t *testing.T) {
	tests := []struct {
		name    string
		configs []RuleConfig
		wantErr string
	}{
		{
			name:    "no rules",
			configs: nil,
		},
		{
			name: "valid rule",
			configs: []RuleConfig{{
				Match:   RuleMatch{Hosts: []string{"smtp.example.com"}, Users: []string{"token"}},
				Headers: map[string]string{"X-Instance-ID": "{{.InstanceID}}"},
			}},
		},
		{
			name: "text and placeholders",
			configs: []RuleConfig{{
				Headers: map[string]string{
					"X-Header": "prefix-{{.InstanceID}}-{{.OrgID}}-suffix",
					"X-Empty":  "",
					"X-Text":   "text",
				},
			}},
		},
		{
			name: "invalid header name",
			configs: []RuleConfig{{
				Headers: map[string]string{"X Invalid": "value"},
			}},
			wantErr: `smtp rule 0: header "X Invalid": invalid name`,
		},
		{
			name: "header name with line break",
			configs: []RuleConfig{{
				Headers: map[string]string{"X-Header\r\nBcc": "value"},
			}},
			wantErr: "invalid name",
		},
		{
			name: "reserved header name, case insensitive",
			configs: []RuleConfig{{}, {
				Headers: map[string]string{"reply-to": "value"},
			}},
			wantErr: `smtp rule 1: header "reply-to": reserved name`,
		},
		{
			name: "header value with line break",
			configs: []RuleConfig{{
				Headers: map[string]string{"X-Header": "value\r\nBcc: attacker@example.com"},
			}},
			wantErr: "value must not contain line breaks",
		},
		{
			name: "empty host",
			configs: []RuleConfig{{
				Match: RuleMatch{Hosts: []string{"smtp.example.com", " "}},
			}},
			wantErr: "smtp rule 0: match host 1: must not be empty",
		},
		{
			name: "empty user",
			configs: []RuleConfig{{
				Match: RuleMatch{Users: []string{""}},
			}},
			wantErr: "smtp rule 0: match user 0: must not be empty",
		},
		{
			name: "empty sender domain",
			configs: []RuleConfig{{
				Match: RuleMatch{SenderDomains: []string{"example.com", ""}},
			}},
			wantErr: "smtp rule 0: match sender domain 1: must not be empty",
		},
		{
			name: "header name with non token characters",
			configs: []RuleConfig{{
				Headers: map[string]string{"X-(Meta)": "value"},
			}},
			wantErr: `smtp rule 0: header "X-(Meta)": invalid name`,
		},
		{
			name: "header names differing in case only",
			configs: []RuleConfig{{
				Headers: map[string]string{"X-Header": "first", "x-header": "second"},
			}},
			wantErr: `smtp rule 0: header "x-header": duplicate name`,
		},
	}
	unsupported := []struct {
		name  string
		value string
	}{
		{"unknown placeholder", "{{.UserID}}"},
		{"placeholder with spaces", "{{ .InstanceID }}"},
		{"trim markers", "{{- .InstanceID -}}"},
		{"condition", "{{if .InstanceID}}{{.OrgID}}{{end}}"},
		{"function", `{{printf "%s" .InstanceID}}`},
		{"pipeline", "{{.InstanceID | len}}"},
		{"chained field", "{{.InstanceID.Foo}}"},
		{"variable", "{{$id := .InstanceID}}{{$id}}"},
		{"unterminated action", "{{.InstanceID"},
		{"stray closing braces", "text}}"},
		{"nested template", `{{define "x"}}{{.OrgID}}{{end}}{{template "x" .}}`},
	}
	for _, u := range unsupported {
		tests = append(tests, struct {
			name    string
			configs []RuleConfig
			wantErr string
		}{
			name: "unsupported value, " + u.name,
			configs: []RuleConfig{{
				Headers: map[string]string{"X-Header": u.value},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, err := CompileRules(tt.configs)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, rules, len(tt.configs))
		})
	}
}

func TestRules_Match(t *testing.T) {
	data := RuleData{InstanceID: "instance1", OrgID: "org1"}
	provider := func(host, user, from string) *Config {
		return &Config{
			SMTP: SMTP{Host: host, PlainAuth: &PlainAuthConfig{User: user, Password: "secret"}},
			From: from,
		}
	}
	defaultProvider := provider("smtp.example.com:587", "token", "noreply@example.com")

	tests := []struct {
		name    string
		configs []RuleConfig
		config  *Config
		want    Rule
	}{
		{
			name:   "no rules",
			config: defaultProvider,
			want:   Rule{},
		},
		{
			name:    "no provider",
			configs: []RuleConfig{{RestrictCustomHTML: true}},
			config:  nil,
			want:    Rule{},
		},
		{
			name:    "empty match applies to all providers",
			configs: []RuleConfig{{RestrictCustomHTML: true}},
			config:  defaultProvider,
			want:    Rule{RestrictCustomHTML: true},
		},
		// hosts
		{
			name: "host without port matches any port",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"smtp.example.com"}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "host with port matches",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"smtp.example.com:587"}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "host with other port does not match",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"smtp.example.com:25"}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{},
		},
		{
			name: "host is case insensitive",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{" SMTP.Example.com "}},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.EXAMPLE.com:587", "token", "noreply@example.com"),
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "provider host without port",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"smtp.example.com"}},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.example.com", "token", "noreply@example.com"),
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "IPv6 host with and without port",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"[2001:db8::1]"}},
				RestrictCustomHTML: true,
			}, {
				Match:              RuleMatch{Hosts: []string{"[2001:db8::2]:2525"}},
				RestrictCustomHTML: true,
			}},
			config: provider("[2001:DB8::1]:2525", "token", "noreply@example.com"),
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "IPv6 host with other port does not match",
			configs: []RuleConfig{{
				Match:              RuleMatch{Hosts: []string{"[2001:db8::1]:25"}},
				RestrictCustomHTML: true,
			}},
			config: provider("[2001:db8::1]:2525", "token", "noreply@example.com"),
			want:   Rule{},
		},
		// users
		{
			name: "user matches",
			configs: []RuleConfig{{
				Match:              RuleMatch{Users: []string{"other", " token "}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "user is case sensitive",
			configs: []RuleConfig{{
				Match:              RuleMatch{Users: []string{"TOKEN"}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{},
		},
		{
			name: "other user does not match, even with the same host and sender domain",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					Users:         []string{"token"},
					SenderDomains: []string{"example.com"},
				},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.example.com:587", "other", "noreply@example.com"),
			want:   Rule{},
		},
		{
			name: "user matches the XOAuth2 user",
			configs: []RuleConfig{{
				Match:              RuleMatch{Users: []string{"oauth-user"}},
				RestrictCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com:587", XOAuth2Auth: &XOAuth2AuthConfig{User: "oauth-user"}}},
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "provider without authentication does not match a user",
			configs: []RuleConfig{{
				Match:              RuleMatch{Users: []string{"token"}},
				RestrictCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com:587"}},
			want:   Rule{},
		},
		// sender domains
		{
			name: "sender domain matches",
			configs: []RuleConfig{{
				Match:              RuleMatch{SenderDomains: []string{"EXAMPLE.com"}},
				RestrictCustomHTML: true,
			}},
			config: defaultProvider,
			want:   Rule{RestrictCustomHTML: true},
		},
		{
			name: "subdomain of sender domain does not match",
			configs: []RuleConfig{{
				Match:              RuleMatch{SenderDomains: []string{"example.com"}},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.example.com:587", "token", "noreply@mail.example.com"),
			want:   Rule{},
		},
		{
			name: "sender without domain does not match",
			configs: []RuleConfig{{
				Match:              RuleMatch{SenderDomains: []string{"example.com"}},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.example.com:587", "token", "example.com"),
			want:   Rule{},
		},
		{
			name: "same host, own sender domain does not match",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					SenderDomains: []string{"example.com"},
				},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.example.com:587", "token", "noreply@other.example"),
			want:   Rule{},
		},
		{
			name: "same sender domain, other host does not match",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					SenderDomains: []string{"example.com"},
				},
				RestrictCustomHTML: true,
			}},
			config: provider("smtp.other.example:587", "token", "noreply@example.com"),
			want:   Rule{},
		},
		// order and headers
		{
			name: "first match wins",
			configs: []RuleConfig{
				{
					Match:   RuleMatch{Hosts: []string{"smtp.other.example"}},
					Headers: map[string]string{"X-Rule": "first"},
				},
				{
					Match:   RuleMatch{Hosts: []string{"smtp.example.com"}},
					Headers: map[string]string{"X-Rule": "second"},
				},
				{
					Headers: map[string]string{"X-Rule": "third"},
				},
			},
			config: defaultProvider,
			want:   Rule{Headers: map[string]string{"X-Rule": "second"}},
		},
		{
			name: "headers are rendered",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					Users:         []string{"token"},
					SenderDomains: []string{"example.com"},
				},
				RestrictCustomHTML: true,
				Headers: map[string]string{
					"X-Instance-ID": "{{.InstanceID}}",
					"X-Org-ID":      "{{.OrgID}}",
					"X-Tag":         "{{.InstanceID}}-{{.OrgID}}-{{.InstanceID}}",
					"X-Static":      "static",
				},
			}},
			config: defaultProvider,
			want: Rule{
				RestrictCustomHTML: true,
				Headers: map[string]string{
					"X-Instance-Id": "instance1",
					"X-Org-Id":      "org1",
					"X-Tag":         "instance1-org1-instance1",
					"X-Static":      "static",
				},
			},
		},
		{
			name: "header names are canonicalized",
			configs: []RuleConfig{{
				Headers: map[string]string{
					"x-pm-metadata-instance-id": "{{.InstanceID}}",
					"X-PM-METADATA-ORG_ID":      "{{.OrgID}}",
				},
			}},
			config: defaultProvider,
			want: Rule{
				Headers: map[string]string{
					"X-Pm-Metadata-Instance-Id": "instance1",
					"X-Pm-Metadata-Org_id":      "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, err := CompileRules(tt.configs)
			require.NoError(t, err)
			assert.Equal(t, tt.want, rules.Match(tt.config, data))
		})
	}
}
