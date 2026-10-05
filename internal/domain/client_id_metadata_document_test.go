package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsClientIDMetadataDocumentURL(t *testing.T) {
	tests := []struct {
		clientID string
		want     bool
	}{
		{"https://app.example.com/oauth/client", true},
		{"https://app.example.com/", true},
		{"https://app.example.com:8443/client", true},
		{"HTTPS://app.example.com/client", true},
		{"https://app.example.com/client?version=1", true},
		{"https://app.example.com/.well-known/client", true},
		{"https://app.example.com/client...json", true},
		{"https://app.example.com", false},
		{"https://:443/client", false},
		{"https://:8443/", false},
		{"https://user@app.example.com/client", false},
		{"https://user:password@app.example.com/client", false},
		{"https://app.example.com/client#fragment", false},
		{"https://app.example.com/client#", false},
		{"https://app.example.com/./client", false},
		{"https://app.example.com/oauth/../client", false},
		{"https://app.example.com/oauth/..", false},
		{"https://app.example.com/oauth/%2e%2e/client", false},
		{"https://app.example.com/oauth/%2E./client", false},
		{"https://app.example.com/oauth%2f..%2fclient", false},
		{"https://app.example.com/oauth%5Cclient", false},
		{"http://app.example.com/client", false},
		{"320948502934092851", false},
		{"320948502934092851@project", false},
		{"app.example.com/client", false},
		{"/oauth/client", false},
		{"", false},
		{"ftp://app.example.com/client", false},
	}
	for _, tt := range tests {
		t.Run(tt.clientID, func(t *testing.T) {
			assert.Equal(t, tt.want, IsClientIDMetadataDocumentURL(tt.clientID))
		})
	}
}

func TestClientIDMetadataDocumentURLAllowed(t *testing.T) {
	allowed := []string{
		"https://app.example.com/oauth/client.json",
		"https://clients.example.com/connectors/",
	}
	tests := []struct {
		clientID string
		want     bool
	}{
		{"https://app.example.com/oauth/client.json", true},
		{"https://clients.example.com/connectors/1234/client.json", true},
		{"https://clients.example.com/connectors/", true},
		{"https://app.example.com/oauth/client.json?x=1", false},
		{"https://app.example.com/oauth/client", false},
		{"https://app.example.com:443/oauth/client.json", false},
		{"https://APP.example.com/oauth/client.json", false},
		{"https://clients.example.com/connectors", false},
		{"https://clients.example.com/other/client.json", false},
		{"https://clients.example.com.evil.example/connectors/client.json", false},
		{"https://evil.example/https://clients.example.com/connectors/", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.clientID, func(t *testing.T) {
			assert.Equal(t, tt.want, ClientIDMetadataDocumentURLAllowed(allowed, tt.clientID))
		})
	}
	t.Run("nothing allowed", func(t *testing.T) {
		assert.False(t, ClientIDMetadataDocumentURLAllowed(nil, "https://app.example.com/oauth/client.json"))
	})
}
