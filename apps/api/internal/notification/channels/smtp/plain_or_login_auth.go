package smtp

import (
	"bytes"
	"net/smtp"
	"slices"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// golang net/smtp SMTP AUTH LOGIN or PLAIN Auth Handler
// Reference: https://gist.github.com/andelf/5118732?permalink_comment_id=4825669#gistcomment-4825669

func PlainOrLoginAuth(username, password, host string) smtp.Auth {
	return &plainOrLoginAuth{username: username, password: password, host: host}
}

type plainOrLoginAuth struct {
	username   string
	password   string
	host       string
	authMethod string
}

func (a *plainOrLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-RRi75", "wrong host name")
	}
	if !slices.Contains(server.Auth, "PLAIN") {
		a.authMethod = "LOGIN"
		return a.authMethod, nil, nil
	} else {
		a.authMethod = "PLAIN"
		resp := []byte("\x00" + a.username + "\x00" + a.password)
		return a.authMethod, resp, nil
	}
}

func (a *plainOrLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	if a.authMethod == "PLAIN" {
		// We've already sent everything.
		return nil, zerrors.ThrowInternal(nil, "SMTP-AAf43", "unexpected server challenge for PLAIN auth method")
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, zerrors.ThrowInternal(nil, "SMTP-HjW21", "unexpected server challenge")
	}
}
