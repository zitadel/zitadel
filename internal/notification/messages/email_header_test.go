package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsReservedEmailHeader(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "From", want: true},
		{name: "from", want: true},
		{name: "Reply-to", want: true},
		{name: "REPLY-TO", want: true},
		{name: "MIME-Version", want: true},
		{name: "content-type", want: true},
		{name: "Bcc", want: true},
		{name: "X-Instance-ID", want: false},
		{name: "X-From", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsReservedEmailHeader(tt.name))
		})
	}
}

func TestIsValidEmailHeaderName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "X-Instance-ID", want: true},
		{name: "X_Header.1", want: true},
		{name: "", want: false},
		{name: "X Header", want: false},
		{name: "X-Header:", want: false},
		{name: "X-Header\r\nBcc", want: false},
		{name: "X-Hëader", want: false},
		{name: "X-(Meta)", want: false},
		{name: "X-Meta/1", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidEmailHeaderName(tt.name))
		})
	}
}

func TestIsValidEmailHeaderValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "simple", value: "instance1", want: true},
		{name: "empty", value: "", want: true},
		{name: "with spaces", value: "some value", want: true},
		{name: "carriage return", value: "value\rBcc: a@example.com", want: false},
		{name: "line feed", value: "value\nBcc: a@example.com", want: false},
		{name: "crlf", value: "value\r\nBcc: a@example.com", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidEmailHeaderValue(tt.value))
		})
	}
}
