package actions

import (
	"net/url"
	"testing"
)

func Test_isHostBlocked(t *testing.T) {
	var denyList = []AddressChecker{
		mustNewIPChecker(t, "192.168.5.0/24"),
		mustNewIPChecker(t, "127.0.0.1"),
		&DomainChecker{Domain: "test.com"},
	}
	type args struct {
		address *url.URL
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "in range",
			args: args{
				address: mustNewURL(t, "https://192.168.5.4/hodor"),
			},
			want: true,
		},
		{
			name: "exact ip",
			args: args{
				address: mustNewURL(t, "http://127.0.0.1:8080/hodor"),
			},
			want: true,
		},
		{
			name: "address match",
			args: args{
				address: mustNewURL(t, "https://test.com:42/hodor"),
			},
			want: true,
		},
		{
			name: "address not match",
			args: args{
				address: mustNewURL(t, "https://test2.com/hodor"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHostBlocked(denyList, tt.args.address); got != tt.want {
				t.Errorf("isHostBlocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustNewIPChecker(t *testing.T, ip string) AddressChecker {
	t.Helper()
	checker, err := NewIPChecker(ip)
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
