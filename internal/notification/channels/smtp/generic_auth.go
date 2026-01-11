package smtp

import (
	"net/smtp"
	"slices"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type GenericAuth struct {
	Host           string
	PlainOrLogin   *plainOrLoginAuth
	XOAuth2        *xoauth2Auth
	selectedMethod smtp.Auth
}

func (a *GenericAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.Host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-BtXgkX", "wrong host name")
	}

	switch {
	case slices.Contains(server.Auth, "XOAUTH2") && a.XOAuth2 != nil:
		a.selectedMethod = a.XOAuth2
	case slices.Contains(server.Auth, "PLAIN") && a.PlainOrLogin != nil:
		a.selectedMethod = a.PlainOrLogin
	case slices.Contains(server.Auth, "LOGIN") && a.PlainOrLogin != nil:
		a.selectedMethod = a.PlainOrLogin
	}

	if a.selectedMethod == nil {
		return "", nil, zerrors.ThrowInternalf(nil, "SMTP-3mBHTz", "no supported auth method found (server supports: %v)", server.Auth)
	}
	return a.selectedMethod.Start(server)
}

func (a *GenericAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if a.selectedMethod != nil {
		return nil, zerrors.ThrowInternal(nil, "SMTP-G99DUr", "no auth method selected")
	}

	return a.selectedMethod.Next(fromServer, more)
}
