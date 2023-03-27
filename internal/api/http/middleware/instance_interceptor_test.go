package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func Test_instanceInterceptor_Handler(t *testing.T) {
	type fields struct {
		verifier   authz.InstanceVerifier
		headerName string
	}
	type args struct {
		request *http.Request
	}
	type res struct {
		statusCode int
		context    context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"setInstance error",
			fields{
				verifier:   &mockInstanceVerifier{},
				headerName: "header",
			},
			args{
				request: httptest.NewRequest("", "/url", nil),
			},
			res{
				statusCode: 404,
				context:    nil,
			},
		},
		{
			"setInstance ok",
			fields{
				verifier:   &mockInstanceVerifier{"host"},
				headerName: "header",
			},
			args{
				request: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r.Header.Set("header", "host")
					return r
				}(),
			},
			res{
				statusCode: 200,
				context:    authz.WithInstance(context.Background(), &mockInstance{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &instanceInterceptor{
				verifier:   tt.fields.verifier,
				headerName: tt.fields.headerName,
			}
			next := &testHandler{}
			got := a.HandlerFunc(next.ServeHTTP)
			rr := httptest.NewRecorder()
			got.ServeHTTP(rr, tt.args.request)
			assert.Equal(t, tt.res.statusCode, rr.Code)
			assert.Equal(t, tt.res.context, next.context)
		})
	}
}

func Test_instanceInterceptor_HandlerFunc(t *testing.T) {
	type fields struct {
		verifier   authz.InstanceVerifier
		headerName string
	}
	type args struct {
		request *http.Request
	}
	type res struct {
		statusCode int
		context    context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"setInstance error",
			fields{
				verifier:   &mockInstanceVerifier{},
				headerName: "header",
			},
			args{
				request: httptest.NewRequest("", "/url", nil),
			},
			res{
				statusCode: 404,
				context:    nil,
			},
		},
		{
			"setInstance ok",
			fields{
				verifier:   &mockInstanceVerifier{"host"},
				headerName: "header",
			},
			args{
				request: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r.Header.Set("header", "host")
					return r
				}(),
			},
			res{
				statusCode: 200,
				context:    authz.WithInstance(context.Background(), &mockInstance{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &instanceInterceptor{
				verifier:   tt.fields.verifier,
				headerName: tt.fields.headerName,
			}
			next := &testHandler{}
			got := a.HandlerFunc(next.ServeHTTP)
			rr := httptest.NewRecorder()
			got.ServeHTTP(rr, tt.args.request)
			assert.Equal(t, tt.res.statusCode, rr.Code)
			assert.Equal(t, tt.res.context, next.context)
		})
	}
}

func Test_setInstance(t *testing.T) {
	type args struct {
		r          *http.Request
		verifier   authz.InstanceVerifier
		headerName string
	}
	type res struct {
		want context.Context
		err  bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"hostname not found, error",
			args{
				r: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					return r
				}(),
				verifier:   &mockInstanceVerifier{},
				headerName: "",
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"invalid host, error",
			args{
				r: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r.Header.Set("header", "host2")
					return r
				}(),
				verifier:   &mockInstanceVerifier{"host"},
				headerName: "header",
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"valid host",
			args{
				r: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r.Header.Set("header", "host")
					return r
				}(),
				verifier:   &mockInstanceVerifier{"host"},
				headerName: "header",
			},
			res{
				want: authz.WithInstance(context.Background(), &mockInstance{}),
				err:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setInstance(tt.args.r, tt.args.verifier, tt.args.headerName)
			if (err != nil) != tt.res.err {
				t.Errorf("setInstance() error = %v, wantErr %v", err, tt.res.err)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("setInstance() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}

type testHandler struct {
	context context.Context
}

func (t *testHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	t.context = r.Context()
}

type mockInstanceVerifier struct {
	host string
}

func (m *mockInstanceVerifier) InstanceByHost(_ context.Context, host string) (authz.Instance, error) {
	if host != m.host {
		return nil, fmt.Errorf("invalid host")
	}
	return &mockInstance{}, nil
}

func (m *mockInstanceVerifier) InstanceByID(context.Context) (authz.Instance, error) {
	return nil, nil
}

type mockInstance struct{}

func (m *mockInstance) InstanceID() string {
	return "instanceID"
}

func (m *mockInstance) ProjectID() string {
	return "projectID"
}

func (m *mockInstance) ConsoleClientID() string {
	return "consoleClientID"
}

func (m *mockInstance) ConsoleApplicationID() string {
	return "consoleApplicationID"
}

func (m *mockInstance) DefaultLanguage() language.Tag {
	return language.English
}

func (m *mockInstance) DefaultOrganisationID() string {
	return "orgID"
}

func (m *mockInstance) RequestedDomain() string {
	return "zitadel.cloud"
}

func (m *mockInstance) RequestedHost() string {
	return "zitadel.cloud:443"
}

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}
