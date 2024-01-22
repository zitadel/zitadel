package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_composeOrigin(t *testing.T) {
	type args struct {
		h               http.Header
		fallBackToHttps bool
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
			fallBackToHttps: false,
		},
		want: "https://host.header",
	}, {
		name: "forwarded host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"host=forwarded.host"},
			},
			fallBackToHttps: false,
		},
		want: "http://forwarded.host",
	}, {
		name: "forwarded proto and host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host"},
			},
			fallBackToHttps: false,
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with multiple complete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http;host=forwarded.host2"},
			},
			fallBackToHttps: false,
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with multiple incomplete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http"},
			},
			fallBackToHttps: false,
		},
		want: "https://forwarded.host",
	}, {
		name: "forwarded proto and host with incomplete entries in different values",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=http", "proto=https;host=forwarded.host", "proto=http"},
			},
			fallBackToHttps: true,
		},
		want: "http://forwarded.host",
	}, {
		name: "x-forwarded-proto https",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
			},
			fallBackToHttps: false,
		},
		want: "https://host.header",
	}, {
		name: "x-forwarded-proto http",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"http"},
			},
			fallBackToHttps: true,
		},
		want: "http://host.header",
	}, {
		name: "fallback to http",
		args: args{
			fallBackToHttps: false,
		},
		want: "http://host.header",
	}, {
		name: "fallback to https",
		args: args{
			fallBackToHttps: true,
		},
		want: "https://host.header",
	}, {
		name: "x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
			fallBackToHttps: false,
		},
		want: "http://x-forwarded.host",
	}, {
		name: "x-forwarded-proto and x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
				"X-Forwarded-Host":  []string{"x-forwarded.host"},
			},
			fallBackToHttps: false,
		},
		want: "https://x-forwarded.host",
	}, {
		name: "forwarded host and x-forwarded-host",
		args: args{
			h: http.Header{
				"Forwarded":        []string{"host=forwarded.host"},
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
			fallBackToHttps: false,
		},
		want: "http://forwarded.host",
	}, {
		name: "forwarded host and x-forwarded-proto",
		args: args{
			h: http.Header{
				"Forwarded":         []string{"host=forwarded.host"},
				"X-Forwarded-Proto": []string{"https"},
			},
			fallBackToHttps: false,
		},
		want: "https://forwarded.host",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, composeOrigin(
				&http.Request{
					Host:   "host.header",
					Header: tt.args.h,
				},
				tt.args.fallBackToHttps,
			), "headers: %+v, fallBackToHttps: %t", tt.args.h, tt.args.fallBackToHttps)
		})
	}
}
