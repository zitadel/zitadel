package net

import (
	builtin_net "net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestHostnameToIPList(t *testing.T) {
	t.Parallel()
	addrErr := &builtin_net.AddrError{Addr: "invalid.local", Err: "Invalid address"}
	tests := []struct {
		name        string
		hostname    string
		lookupFunc  IPLookupFunc
		want        []builtin_net.IP
		expectedErr error
	}{
		{
			name:       "valid IP address",
			hostname:   "192.168.1.1",
			lookupFunc: nil,
			want:       []builtin_net.IP{builtin_net.ParseIP("192.168.1.1")},
		},
		{
			name:     "domain with lookup function",
			hostname: "example.com",
			lookupFunc: func(s string) ([]builtin_net.IP, error) {
				return []builtin_net.IP{builtin_net.ParseIP("127.0.0.1")}, nil
			},
			want: []builtin_net.IP{builtin_net.ParseIP("127.0.0.1")},
		},
		{
			name:        "domain without lookup function",
			hostname:    "example.com",
			lookupFunc:  nil,
			want:        nil,
			expectedErr: zerrors.ThrowInvalidArgument(nil, "NET-naSn77", "lookup function must not be nil"),
		},
		{
			name:     "domain with lookup error",
			hostname: "invalid.local",
			lookupFunc: func(s string) ([]builtin_net.IP, error) {
				return nil, addrErr
			},
			want:        nil,
			expectedErr: zerrors.ThrowInternal(addrErr, "NET-4m9s2", "lookup failed"),
		},
		{
			name:     "domain with multiple IPs",
			hostname: "multi.example.com",
			lookupFunc: func(s string) ([]builtin_net.IP, error) {
				return []builtin_net.IP{
					builtin_net.ParseIP("127.0.0.1"),
					builtin_net.ParseIP("127.0.0.2"),
				}, nil
			},
			want: []builtin_net.IP{
				builtin_net.ParseIP("127.0.0.1"),
				builtin_net.ParseIP("127.0.0.2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := HostnameToIPList(tt.hostname, tt.lookupFunc)

			require.ErrorIs(t, err, tt.expectedErr)
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}
