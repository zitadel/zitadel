package command

import (
	"testing"
)

func Test_validateHostedDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		// valid domains
		{name: "empty (no restriction)", domain: "", wantErr: false},
		{name: "simple domain", domain: "example.com", wantErr: false},
		{name: "subdomain", domain: "corp.example.com", wantErr: false},
		{name: "deep subdomain", domain: "a.b.c.example.com", wantErr: false},
		{name: "hyphen in label", domain: "my-company.com", wantErr: false},
		{name: "numeric label", domain: "example123.com", wantErr: false},
		{name: "long TLD", domain: "example.cloud", wantErr: false},
		// invalid domains
		{name: "no dot", domain: "example", wantErr: true},
		{name: "leading dot", domain: ".example.com", wantErr: true},
		{name: "trailing dot", domain: "example.com.", wantErr: true},
		{name: "leading hyphen in label", domain: "-example.com", wantErr: true},
		{name: "trailing hyphen in label", domain: "example-.com", wantErr: true},
		{name: "single-char TLD", domain: "example.c", wantErr: true},
		{name: "spaces", domain: "my company.com", wantErr: true},
		{name: "with scheme", domain: "https://example.com", wantErr: true},
		{name: "with path", domain: "example.com/path", wantErr: true},
		{name: "at sign", domain: "user@example.com", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHostedDomain(tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHostedDomain(%q) error = %v, wantErr %v", tt.domain, err, tt.wantErr)
			}
		})
	}
}
