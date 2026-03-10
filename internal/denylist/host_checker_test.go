package denylist

import (
	"io"
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestIsHostBlocked(t *testing.T) {
	t.Parallel()
	var denyList = []AddressChecker{
		mustNewHostChecker(t, "192.168.5.0/24"),
		mustNewHostChecker(t, "127.0.0.1"),
		mustNewHostChecker(t, "test.com"),
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
			name: "in range",
			args: args{
				address: mustNewURL(t, "https://192.168.5.4/hodor"),
			},
			want: NewAddressDeniedError("192.168.5.0/24"),
		},
		{
			name: "exact ip",
			args: args{
				address: mustNewURL(t, "http://127.0.0.1:8080/hodor"),
			},
			want: NewAddressDeniedError("127.0.0.1"),
		},
		{
			name: "address match",
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
			name: "address not match",
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
			name: "looked up ip matches",
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
			name: "lookup failure",
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
			got := IsHostBlocked(denyList, tt.args.address, tt.fields.lookup)
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

func mustNewHostChecker(t *testing.T, ip string) AddressChecker {
	t.Helper()
	checker, err := NewHostChecker(ip)
	if err != nil {
		t.Errorf("unable to parse cidr of %q because: %v", ip, err)
		t.FailNow()
	}
	return checker
}

func mustNewURL(t *testing.T, raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		t.Errorf("unable to parse address of %q because: %v", raw, err)
		t.FailNow()
	}
	return u
}
