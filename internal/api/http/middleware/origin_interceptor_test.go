package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_composeOrigin(t *testing.T) {
	type args struct {
		h            http.Header
		enforceHttps bool
		// defaults to x-zitadel-public-host
		publicHeaders []string
	}
	tests := []struct {
		name    string
		args    args
		want    *http_util.DomainCtx
		wantErr bool
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
	}, {
		name: "public host",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"public.host"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			PublicHost:   "public.host",
			Protocol:     "http",
		},
	}, {
		name: "public host with port",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"public.host:8080"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			PublicHost:   "public.host:8080",
			Protocol:     "http",
		},
	}, {
		name: "invalid public host (non-numeric port) is rejected",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"evil.com:443/"},
			},
		},
		wantErr: true,
	}, {
		name: "invalid public host (userinfo) is rejected",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"a:b@evil.com"},
			},
		},
		wantErr: true,
	}, {
		// the trusted domain check must not see "trusted.example.com", while links would point to "evil.example"
		name: "invalid public host (userinfo with trusted host and port) is rejected",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"trusted.example.com:443@evil.example"},
			},
		},
		wantErr: true,
	}, {
		name: "invalid public host (userinfo without port) is rejected",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"trusted.example.com@evil.example"},
			},
		},
		wantErr: true,
	}, {
		name: "invalid instance host (userinfo with trusted host and port) is ignored",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"trusted.example.com:443@evil.example"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	}, {
		name: "invalid instance host is ignored",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"localhost:@attacker.com"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	}, {
		name: "x-forwarded-for does not contain a host",
		args: args{
			h: http.Header{
				"X-Forwarded-For": []string{"host=evil.example"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	}, {
		name: "multiple instance host values use the first value",
		args: args{
			h: http.Header{
				"X-Forwarded-Host": []string{"first.host", "second.host"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "first.host",
			Protocol:     "http",
		},
	}, {
		name: "public host with empty port is normalized",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"public.host:"},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			PublicHost:   "public.host",
			Protocol:     "http",
		},
	}, {
		name: "multiple public host values are rejected",
		args: args{
			h: http.Header{
				"X-Zitadel-Public-Host": []string{"trusted.example.com", "evil.example"},
			},
		},
		wantErr: true,
	}, {
		// first match wins: the invalid higher-ranked header is skipped and the valid lower-ranked
		// header is used. The chosen host is still validated and checked against the trusted domains,
		// so no poisoned host can reach the links.
		name: "invalid public host header is skipped when a lower-ranked one is valid",
		args: args{
			h: http.Header{
				"X-Public-A": []string{"trusted.example.com:443@evil.example"},
				"X-Public-B": []string{"trusted.example.com"},
			},
			publicHeaders: []string{"x-public-a", "x-public-b"},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			PublicHost:   "trusted.example.com",
			Protocol:     "http",
		},
	}, {
		name: "unparsable forwarded public host is rejected",
		args: args{
			h: http.Header{
				"Forwarded": []string{`host="a:b@evil.com`},
			},
			publicHeaders: []string{http_util.Forwarded},
		},
		wantErr: true,
	}, {
		name: "forwarded public host with invalid host is rejected",
		args: args{
			h: http.Header{
				"Forwarded": []string{`host="trusted.example.com:443@evil.example"`},
			},
			publicHeaders: []string{http_util.Forwarded},
		},
		wantErr: true,
	}, {
		name: "forwarded public host without host is ignored",
		args: args{
			h: http.Header{
				"Forwarded": []string{"for=192.0.2.60;proto=https"},
			},
			publicHeaders: []string{http_util.Forwarded},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "https",
		},
	}, {
		name: "unparsable forwarded instance host is ignored",
		args: args{
			h: http.Header{
				"Forwarded": []string{`host="a:b@evil.com`},
			},
		},
		want: &http_util.DomainCtx{
			InstanceHost: "host.header",
			Protocol:     "http",
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicHeaders := tt.args.publicHeaders
			if publicHeaders == nil {
				publicHeaders = []string{"x-zitadel-public-host"}
			}
			got, err := composeDomainContext(
				&http.Request{
					Host:   "host.header",
					Header: tt.args.h,
				},
				tt.args.enforceHttps,
				[]string{http_util.Forwarded, http_util.ForwardedFor, http_util.ForwardedHost, http_util.ForwardedProto},
				publicHeaders,
			)
			if tt.wantErr {
				assert.True(t, zerrors.IsErrorInvalidArgument(err), "expected invalid argument error, got: %v", err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "headers: %+v, enforceHttps: %t", tt.args.h, tt.args.enforceHttps)
		})
	}
}

func TestWithOrigin_keepsHeaderConfig(t *testing.T) {
	// a config slice with spare capacity must not be written to (shared between concurrent requests)
	instanceHostHeaders := make([]string, 1, 10)
	instanceHostHeaders[0] = "x-zitadel-instance-host"

	WithOrigin(false, "host", ":authority", instanceHostHeaders, []string{"x-zitadel-public-host"})

	assert.Equal(t, []string{"x-zitadel-instance-host"}, instanceHostHeaders)
	assert.Empty(t, instanceHostHeaders[:2][1], "backing array of the config slice was modified")
}

func TestWithOrigin(t *testing.T) {
	tests := []struct {
		name           string
		publicHost     string
		wantStatus     int
		wantNextCalled bool
		wantPublicHost string
	}{
		{
			name:           "no public host",
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "valid public host",
			publicHost:     "trusted.example.com",
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
			wantPublicHost: "trusted.example.com",
		},
		{
			name:       "invalid public host (userinfo with trusted host and port)",
			publicHost: "trusted.example.com:443@evil.example",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid public host (userinfo)",
			publicHost: "a:b@evil.com",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				nextCalled bool
				domainCtx  *http_util.DomainCtx
			)
			handler := WithOrigin(false, "", "", nil, []string{"x-zitadel-public-host"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				domainCtx = http_util.DomainContext(r.Context())
			}))
			req := httptest.NewRequest(http.MethodGet, "http://host.header/", nil)
			if tt.publicHost != "" {
				req.Header.Set("x-zitadel-public-host", tt.publicHost)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantStatus == http.StatusBadRequest {
				assert.Contains(t, rec.Body.String(), "HTTP-Rk7pX")
			}
			assert.Equal(t, tt.wantNextCalled, nextCalled)
			if tt.wantNextCalled {
				assert.Equal(t, tt.wantPublicHost, domainCtx.PublicHost)
			}
		})
	}
}
