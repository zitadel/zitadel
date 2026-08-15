package denylist

import (
	"io"
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIsURLBlocked(t *testing.T) {
	t.Parallel()

	denyList := []AddressChecker{
		NewHostChecker("192.168.5.0/24"),
		NewHostChecker("127.0.0.1"),
		NewHostChecker("test.com"),
	}

	type fields struct {
		lookup func(host string) ([]net.IP, error)
	}
	type args struct {
		address *url.URL
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "in range (CIDR match)",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("192.168.5.4")}, nil
				},
			},
			args: args{
				address: mustNewURL(t, "https://192.168.5.4/hodor"),
			},
			want: NewAddressDeniedError("192.168.5.0/24"),
		},
		{
			name: "exact ip literal match",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("127.0.0.1")}, nil
				},
			},
			args: args{
				address: mustNewURL(t, "http://127.0.0.1:8080/hodor"),
			},
			want: NewAddressDeniedError("127.0.0.1"),
		},
		{
			name: "domain string match",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("194.264.52.4")}, nil
				},
			},
			args: args{
				address: mustNewURL(t, "https://test.com:42/hodor"),
			},
			want: NewAddressDeniedError("test.com"),
		},
		{
			name: "allowed domain string pass",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("194.264.52.4")}, nil
				},
			},
			args: args{
				address: mustNewURL(t, "https://test2.com/hodor"),
			},
			want: nil,
		},
		{
			name: "unmatched domain resolves to blocked IP (SSRF pivot attempt)",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return []net.IP{net.ParseIP("127.0.0.1")}, nil
				},
			},
			args: args{
				address: mustNewURL(t, "https://test2.com/hodor"),
			},
			want: NewAddressDeniedError("127.0.0.1"),
		},
		{
			name: "dns resolver failures are isolated",
			fields: fields{
				lookup: func(host string) ([]net.IP, error) {
					return nil, io.EOF
				},
			},
			args: args{
				address: mustNewURL(t, "https://test2.com/hodor"),
			},
			want: zerrors.ThrowInternal(io.EOF, "NET-4m9s2", "lookup failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsURLBlocked(denyList, tt.args.address, tt.fields.lookup)
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

func TestIsHostnameBlocked_EdgeCases(t *testing.T) {
	t.Parallel()

	denyList := []AddressChecker{
		NewHostChecker("127.0.0.1"),
		NewHostChecker("blocked.example.com"),
	}

	t.Run("blocked by host lookup via hostname method", func(t *testing.T) {
		t.Parallel()
		err := IsHostnameBlocked(denyList, "safe.example.com", func(_ string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		})
		assert.ErrorIs(t, err, NewAddressDeniedError("127.0.0.1"))
	})

	t.Run("blocked by exact domain string mapping", func(t *testing.T) {
		t.Parallel()
		err := IsHostnameBlocked(denyList, "blocked.example.com", func(_ string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("192.0.2.1")}, nil
		})
		assert.ErrorIs(t, err, NewAddressDeniedError("blocked.example.com"))
	})

	t.Run("ipv6 target check compatibility", func(t *testing.T) {
		t.Parallel()
		// URL parsing or raw string handling yields plain IPv6 targets, make sure checker behaves cleanly
		err := IsHostnameBlocked(denyList, "::1", func(_ string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("::1")}, nil
		})
		assert.NoError(t, err)
	})
}

func mustNewURL(t *testing.T, raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		t.Errorf("unable to parse address of %q because: %v", raw, err)
		t.FailNow()
	}
	return u
}
