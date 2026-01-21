package smtp

import (
	"context"
	"crypto/sha256"
	"net/smtp"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var tokenSourceCache *lru.TwoQueueCache[string, oauth2.TokenSource]

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

	var err error
	if tokenSourceCache == nil {
		tokenSourceCache, err = lru.New2Q[string, oauth2.TokenSource](100)
		if err != nil {
			return "", nil, zerrors.ThrowInternal(nil, "SMTP-xcQ4WA", "failed to create token source cache")
		}
	}

	hash := a.config.Hash()
	tokenSource, ok := tokenSourceCache.Get(hash)
	if !ok {
		if a.config.ClientCredentialsAuth != nil {
			config := &clientcredentials.Config{
				ClientID:     a.config.ClientCredentialsAuth.ClientId,
				ClientSecret: a.config.ClientCredentialsAuth.ClientSecret,
				Scopes:       a.config.Scopes,
				TokenURL:     a.config.TokenEndpoint,
			}
			tokenSource = config.TokenSource(context.Background())
			tokenSourceCache.Add(hash, tokenSource)
		}
	}

	token, err := tokenSource.Token()
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
