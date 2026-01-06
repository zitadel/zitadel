package smtp

import (
	"net"
	"net/smtp"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func OAuthBearerAuth(config OAuthBearerAuthConfig, host string) smtp.Auth {
	return &oauthBearerAuth{
		host:   host,
		config: config,
	}
}

type oauthBearerAuth struct {
	host   string
	config OAuthBearerAuthConfig
}

type OAuthBearerAuthConfig struct {
	User        string
	BearerToken string
}

func (a *oauthBearerAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-dKgdBJ", "wrong host name")
	}

	host, port, err := net.SplitHostPort(a.host)
	if err != nil {
		host = a.host
		port = "587"
	}

	resp := []byte("n,a=" + a.config.User + ",\x01host=" + host + "\x01port=" + port + "\x01" + "auth=Bearer " + a.config.BearerToken + "\x01\x01")
	return "OAUTHBEARER", resp, nil
}

func (a *oauthBearerAuth) Next(_ []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	return nil, zerrors.ThrowInternal(nil, "SMTP-TkAf1W", "unexpected server challenge for XOAUTH2 auth method")
}
