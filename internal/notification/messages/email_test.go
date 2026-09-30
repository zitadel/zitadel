package messages

import (
	"io"
	"net/mail"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmail_GetContent(t *testing.T) {
	tests := []struct {
		name        string
		headers     map[string]string
		wantHeaders map[string]string
		wantMissing []string
	}{
		{
			name:    "no additional headers",
			headers: nil,
		},
		{
			name: "additional headers are added",
			headers: map[string]string{
				"X-Instance-ID": "instance1",
				"X-Org-ID":      "org1",
			},
			wantHeaders: map[string]string{
				"X-Instance-ID": "instance1",
				"X-Org-ID":      "org1",
			},
		},
		{
			name: "headers set by ZITADEL are not overwritten",
			headers: map[string]string{
				"from":         "attacker@other.example",
				"Reply-To":     "attacker@other.example",
				"SUBJECT":      "other subject",
				"Content-Type": "text/plain",
				"X-Valid":      "valid",
			},
			wantHeaders: map[string]string{
				"X-Valid": "valid",
			},
		},
		{
			name: "invalid headers are ignored",
			headers: map[string]string{
				"X-Line-Break":                 "value\r\nBcc: attacker@other.example",
				"X-Name\r\nBcc":                "attacker@other.example",
				"X Space":                      "value",
				"X-Colon: injected\r\nX-Other": "value",
				"":                             "value",
				"X-Valid":                      "valid",
			},
			wantHeaders: map[string]string{
				"X-Valid": "valid",
			},
			wantMissing: []string{"X-Line-Break", "X-Name", "Bcc", "X Space", "X-Colon", "X-Other"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Email{
				Recipients:     []string{"user@example.com"},
				SenderEmail:    "noreply@example.com",
				SenderName:     "Sender",
				ReplyToAddress: "reply@example.com",
				Subject:        "Subject",
				Content:        "<html><body>Content</body></html>",
				Headers:        tt.headers,
			}
			content, err := msg.GetContent()
			require.NoError(t, err)

			parsed, err := mail.ReadMessage(strings.NewReader(content))
			require.NoError(t, err)
			header := textproto.MIMEHeader(parsed.Header)

			// headers set by ZITADEL are never changed by additional headers
			assert.Equal(t, "Sender <noreply@example.com>", header.Get("From"))
			assert.Equal(t, "reply@example.com", header.Get("Reply-To"))
			assert.Equal(t, "noreply@example.com", header.Get("Return-Path"))
			assert.Equal(t, "user@example.com", header.Get("To"))
			assert.Equal(t, "Subject", header.Get("Subject"))
			assert.Equal(t, []string{`text/html; charset="UTF-8"`}, header.Values("Content-Type"))
			assert.Len(t, header.Values("From"), 1)
			assert.Len(t, header.Values("Reply-To"), 1)
			assert.Len(t, header.Values("Subject"), 1)

			for name, value := range tt.wantHeaders {
				assert.Equal(t, []string{value}, header.Values(name), name)
			}
			for _, name := range tt.wantMissing {
				assert.Empty(t, header.Values(name), name)
			}

			body, err := io.ReadAll(parsed.Body)
			require.NoError(t, err)
			// the content is preceded by an additional line break
			assert.Equal(t, "\r\n<html><body>Content</body></html>", string(body))
		})
	}
}
