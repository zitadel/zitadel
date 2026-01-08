package smtp

import (
	"context"
	"net/smtp"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func XOAuth2Auth(config XOAuth2AuthConfig, host string) smtp.Auth {
	return &xoauth2Auth{
		host:   host,
		config: config,
	}
}

type xoauth2Auth struct {
	host        string
	config      XOAuth2AuthConfig
	tokenSource oauth2.TokenSource
}

type XOAuth2AuthConfig struct {
	User          string
	ClientId      string
	ClientSecret  string
	TokenEndpoint string
	Scopes        []string
}

func (a *xoauth2Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-eRJLyi", "wrong host name")
	}

	if a.tokenSource == nil {
		config := &clientcredentials.Config{
			ClientID:     a.config.ClientId,
			ClientSecret: a.config.ClientSecret,
			Scopes:       a.config.Scopes,
			TokenURL:     a.config.TokenEndpoint,
		}
		a.tokenSource = config.TokenSource(context.Background())
	}

	token, err := a.tokenSource.Token()
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
