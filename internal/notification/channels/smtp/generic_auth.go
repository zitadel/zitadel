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
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-TODO-GET_ERROR_CODE", "wrong host name")
	}

	if slices.Contains(server.Auth, "XOAUTH2") && a.XOAuth2 != nil {
		a.selectedMethod = a.XOAuth2
	} else if slices.Contains(server.Auth, "PLAIN") && a.PlainOrLogin != nil {
		a.selectedMethod = a.PlainOrLogin
	} else if slices.Contains(server.Auth, "LOGIN") && a.PlainOrLogin != nil {
		a.selectedMethod = a.PlainOrLogin
	}

	if a.selectedMethod == nil {
		return "", nil, zerrors.ThrowInternalf(nil, "SMTP-TODO-GET_ERROR_CODE", "no supported auth method found (server supports: %v)", server.Auth)
	}
	return a.selectedMethod.Start(server)
}

func (a *GenericAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if a.selectedMethod != nil {
		return nil, zerrors.ThrowInternal(nil, "SMTP-TODO-GET_ERROR_CODE", "no auth method selected")
	}

	return a.selectedMethod.Next(fromServer, more)
}
