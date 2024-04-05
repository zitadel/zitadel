package smtp

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"slices"
)

// golang net/smtp SMTP AUTH LOGIN Auth Handler
// Reference: https://gist.github.com/andelf/5118732?permalink_comment_id=4825669#gistcomment-4825669
func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

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
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises PLAIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
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
		return nil, errors.New("unexpected server challenge")
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("unexpected server challenge: %s", fromServer)
	}
}
