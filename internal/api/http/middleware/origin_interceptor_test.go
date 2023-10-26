package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_composeOrigin(t *testing.T) {
	type args struct {
		h http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "no proxy headers",
		want: "http://host.header",
	}, {
		name: "forwarded proto",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https"},
			},
		},
		want: "https://host.header",
	}, {
		name: "forwarded host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"host=forwarded.host"},
			},
		},
		want: "http://forwarded.host",
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
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, composeOrigin(&http.Request{
				Host:   "host.header",
				Header: tt.args.h,
			}), "headers: %+v", tt.args.h)
		})
	}
}
