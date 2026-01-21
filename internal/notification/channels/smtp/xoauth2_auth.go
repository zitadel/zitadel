package smtp

import (
	"context"
	"crypto/sha256"
	"net/smtp"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var tokenSourceCache = TokenSourceCache{}

type TokenSourceCache map[string]tokenSourceCacheEntry

type tokenSourceCacheEntry struct {
	lastUsedAt  time.Time
	tokenSource oauth2.TokenSource
}

func (cache TokenSourceCache) Cleanup() {
	maxAge := 7 * 24 * time.Hour
	now := time.Now().UTC()
	oldestAllowed := now.Add(maxAge)
	for k, v := range tokenSourceCache {
		if v.lastUsedAt.Before(oldestAllowed) {
			delete(tokenSourceCache, k)
		}
	}
}

type XOAuth2AuthConfig struct {
	User                  string
	TokenEndpoint         string
	Scopes                []string
	ClientCredentialsAuth *OAuth2ClientCredentials
}

type OAuth2ClientCredentials struct {
	ClientId     string
	ClientSecret string
}

func (cfg XOAuth2AuthConfig) Hash() string {
	sha := sha256.New()
	sha.Write([]byte(cfg.User))
	sha.Write([]byte(cfg.TokenEndpoint))
	for _, s := range cfg.Scopes {
		sha.Write([]byte(s))
	}
	if cfg.ClientCredentialsAuth != nil {
		sha.Write([]byte(cfg.ClientCredentialsAuth.ClientId))
		sha.Write([]byte(cfg.ClientCredentialsAuth.ClientSecret))
	}
	return string(sha.Sum(nil))
}

type xoauth2Auth struct {
	host   string
	config XOAuth2AuthConfig
}

func XOAuth2Auth(config XOAuth2AuthConfig, host string) smtp.Auth {
	return &xoauth2Auth{
		host:   host,
		config: config,
	}
}
func (a *xoauth2Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-eRJLyi", "wrong host name")
	}

	hash := a.config.Hash()
	entry, ok := tokenSourceCache[hash]
	if !ok {
		config := &clientcredentials.Config{
			ClientID:     a.config.ClientCredentialsAuth.ClientId,
			ClientSecret: a.config.ClientCredentialsAuth.ClientSecret,
			Scopes:       a.config.Scopes,
			TokenURL:     a.config.TokenEndpoint,
		}
		entry.tokenSource = config.TokenSource(context.Background())
	}

	entry.lastUsedAt = time.Now().UTC()
	tokenSourceCache[hash] = entry

	token, err := entry.tokenSource.Token()
	if err != nil {
		return "", nil, zerrors.ThrowInternal(err, "SMTP-fjHcJu", "Failed to get token to connect with smtp server")
	}

	resp := []byte("user=" + a.config.User + "\x01" + "auth=Bearer " + token.AccessToken + "\x01\x01")
	return "XOAUTH2", resp, nil
}

func (a *xoauth2Auth) Next(_ []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	return nil, zerrors.ThrowInternal(nil, "SMTP-Pqsrj9", "unexpected server challenge for XOAUTH2 auth method")
}
