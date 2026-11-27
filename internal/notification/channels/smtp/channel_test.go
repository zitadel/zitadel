package smtp

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEhloName(t *testing.T) {
	tests := []struct {
		senderAddress string
		want          string
	}{
		{"noreply@accounts.example.com", "accounts.example.com"},
		{"No Reply <noreply@example.com>", "example.com"},
		{"", ""},
		{"no-at-sign", ""},
	}
	for _, tt := range tests {
		t.Run(tt.senderAddress, func(t *testing.T) {
			assert.Equal(t, tt.want, ehloName(tt.senderAddress))
		})
	}
}

// serveRejectingLocalhost behaves like Google's SMTP relay: it closes the connection when greeted as "localhost".
func serveRejectingLocalhost(t *testing.T) (addr string, greetings chan string) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { l.Close() })
	greetings = make(chan string, 1)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		write := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
		write("220 test ESMTP")
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			line = strings.TrimSpace(line)
			cmd := strings.ToUpper(strings.SplitN(line, " ", 2)[0])
			switch cmd {
			case "EHLO", "HELO":
				greetings <- line
				if strings.EqualFold(strings.TrimSpace(line[4:]), "localhost") {
					write("421 4.7.0 Try again later, closing connection. (EHLO)")
					return
				}
				write("250 test")
			case "MAIL":
				write("250 OK")
			case "QUIT":
				write("221 bye")
				return
			default:
				write("502 not implemented")
			}
		}
	}()
	return l.Addr().String(), greetings
}

func TestConnectToSMTP_EhloName(t *testing.T) {
	tests := []struct {
		name     string
		ehloName string
		wantEhlo string
		wantErr  bool
	}{
		{"sender domain", "accounts.example.com", "EHLO accounts.example.com", false},
		{"net/smtp default", "", "EHLO localhost", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, greetings := serveRejectingLocalhost(t)
			client, err := SMTP{Host: addr}.connectToSMTP(false, tt.ehloName)
			require.NoError(t, err)
			defer client.Close()

			err = client.Mail("noreply@accounts.example.com")
			assert.Equal(t, tt.wantEhlo, <-greetings)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NoError(t, client.Quit())
		})
	}
}
