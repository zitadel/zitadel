package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_composeOrigin(t *testing.T) {
	const hostHeaderValue = "host.header"
	type args struct {
		h http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "no proxy headers",
		want: "http://" + hostHeaderValue,
	}, {
		name: "forwarded host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"host=forwarded.host"},
			},
		},
		want: "http://forwarded.host",
	}, /*{ // TODO: Incomment once we support the proto directive in the Forwarded and the X-Forwarded-* headers
		name: "forwarded proto",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https"},
			},
		},
		want: "https://host.header",
	}, {
		name: "forwarded proto and host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host"},
			},
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with multiple complete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http;host=forwarded.host2"},
			},
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with multiple incomplete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http"},
			},
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with incomplete entries in different values",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=http", "proto=https;host=forwarded.host", "proto=http"},
			},
		},
		want: "http://forwarded.host",
	}, {
		name: "x-forwarded-proto",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
			},
		},
		want: "https://host.header",
	}, {
		name: "x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
		},
		want: "http://x-forwarded.host",
	}, {
		name: "x-forwarded-proto and x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
				"X-Forwarded-Host":  []string{"x-forwarded.host"},
			},
		},
		want: "https://x-forwarded.host",
	}, {
		name: "forwarded host and x-forwarded-host",
		args: args{
			h: http.Header{
				"Forwarded":        []string{"host=forwarded.host"},
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
		},
		want: "http://forwarded.host",
	}, {
		name: "forwarded host and x-forwarded-proto",
		args: args{
			h: http.Header{
				"Forwarded":         []string{"host=forwarded.host"},
				"X-Forwarded-Proto": []string{"https"},
			},
		},
		want: "https://forwarded.host",
	}*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin, err := composeOrigin(&http.Request{
				Host:   hostHeaderValue,
				Header: tt.args.h,
			}, false, "", "")
			require.NoError(t, err)
			assert.Equalf(t, tt.want, origin.Full, "headers: %+v", tt.args.h)
		})
	}
}
