package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func Test_composeOrigin(t *testing.T) {
	type args struct {
		h            http.Header
		enforceHttps bool
	}
	tests := []struct {
		name string
		args args
		want *http_util.DomainCtx
	}{{
		name: "no proxy headers",
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	}, {
		name: "forwarded proto",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "https",
		},
	}, {
		name: "forwarded host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"host=forwarded.host"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "http",
		},
	}, {
		name: "forwarded proto and host",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "https",
		},
	}, {
		name: "forwarded proto and host with multiple complete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http;host=forwarded.host2"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "https",
		},
	}, {
		name: "forwarded proto and host with multiple incomplete entries",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=https;host=forwarded.host, proto=http"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "https",
		},
	}, {
		name: "forwarded proto and host with incomplete entries in different values",
		args: args{
			h: http.Header{
				"Forwarded": []string{"proto=http", "proto=https;host=forwarded.host", "proto=http"},
			},
			enforceHttps: true,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "https",
		},
	}, {
		name: "x-forwarded-proto https",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "https",
		},
	}, {
		name: "x-forwarded-proto http",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"http"},
			},
			enforceHttps: true,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "https",
		},
	}, {
		name: "fallback to http",
		args: args{
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	}, {
		name: "enforce https",
		args: args{
			enforceHttps: true,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "https",
		},
	}, {
		name: "x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "x-forwarded.host",
			Protocol:     "http",
		},
	}, {
		name: "x-forwarded-proto and x-forwarded-host",
		args: args{
			h: http.Header{
				"X-Forwarded-Proto": []string{"https"},
				"X-Forwarded-Host":  []string{"x-forwarded.host"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "x-forwarded.host",
			Protocol:     "https",
		},
	}, {
		name: "forwarded host and x-forwarded-host",
		args: args{
			h: http.Header{
				"Forwarded":        []string{"host=forwarded.host"},
				"X-Forwarded-Host": []string{"x-forwarded.host"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "http",
		},
	}, {
		name: "forwarded host and x-forwarded-proto",
		args: args{
			h: http.Header{
				"Forwarded":         []string{"host=forwarded.host"},
				"X-Forwarded-Proto": []string{"https"},
			},
			enforceHttps: false,
		},
		want: &http_util.DomainCtx{
			InstanceHost: "forwarded.host",
			Protocol:     "https",
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, composeDomainContext(
				&http.Request{
					Host:   "host.header",
					Header: tt.args.h,
				},
				tt.args.enforceHttps,
				[]string{http_util.Forwarded, http_util.ForwardedFor, http_util.ForwardedHost, http_util.ForwardedProto},
				[]string{"x-zitadel-public-host"},
			), "headers: %+v, enforceHttps: %t", tt.args.h, tt.args.enforceHttps)
		})
	}
}
