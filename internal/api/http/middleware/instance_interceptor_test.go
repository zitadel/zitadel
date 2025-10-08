package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/i18n"
)

func Test_instanceInterceptor_Handler(t *testing.T) {
	type fields struct {
		verifier authz.InstanceVerifier
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
				verifier: &mockInstanceVerifier{},
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
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			args{
				request: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r = r.WithContext(zitadel_http.WithDomainContext(r.Context(), &zitadel_http.DomainCtx{InstanceHost: "host"}))
					return r
				}(),
			},
			res{
				statusCode: 200,
				context:    authz.WithInstance(zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "host"}), &mockInstance{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &instanceInterceptor{
				verifier:   tt.fields.verifier,
				translator: i18n.NewZitadelTranslator(language.English),
			}
			next := &testHandler{}
			got := a.HandlerFunc(next)
			rr := httptest.NewRecorder()
			got.ServeHTTP(rr, tt.args.request)
			assert.Equal(t, tt.res.statusCode, rr.Code)
			assert.Equal(t, tt.res.context, next.context)
		})
	}
}

func Test_instanceInterceptor_HandlerFunc(t *testing.T) {
	type fields struct {
		verifier authz.InstanceVerifier
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
				verifier: &mockInstanceVerifier{},
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
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			args{
				request: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r = r.WithContext(zitadel_http.WithDomainContext(r.Context(), &zitadel_http.DomainCtx{InstanceHost: "host"}))
					return r
				}(),
			},
			res{
				statusCode: 200,
				context:    authz.WithInstance(zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "host"}), &mockInstance{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &instanceInterceptor{
				verifier:   tt.fields.verifier,
				translator: i18n.NewZitadelTranslator(language.English),
			}
			next := &testHandler{}
			got := a.HandlerFunc(next)
			rr := httptest.NewRecorder()
			got.ServeHTTP(rr, tt.args.request)
			assert.Equal(t, tt.res.statusCode, rr.Code)
			assert.Equal(t, tt.res.context, next.context)
		})
	}
}

func Test_instanceInterceptor_HandlerFuncWithError(t *testing.T) {
	type fields struct {
		verifier authz.InstanceVerifier
	}
	type args struct {
		request *http.Request
	}
	type res struct {
		wantErr bool
		context context.Context
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
				verifier: &mockInstanceVerifier{},
			},
			args{
				request: httptest.NewRequest("", "/url", nil),
			},
			res{
				wantErr: true,
				context: nil,
			},
		},
		{
			"setInstance ok",
			fields{
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			args{
				request: func() *http.Request {
					r := httptest.NewRequest("", "/url", nil)
					r = r.WithContext(zitadel_http.WithDomainContext(r.Context(), &zitadel_http.DomainCtx{InstanceHost: "host"}))
					return r
				}(),
			},
			res{
				context: authz.WithInstance(zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "host"}), &mockInstance{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &instanceInterceptor{
				verifier:   tt.fields.verifier,
				translator: i18n.NewZitadelTranslator(language.English),
			}
			var ctx context.Context
			got := a.HandlerFuncWithError(func(w http.ResponseWriter, r *http.Request) error {
				ctx = r.Context()
				return nil
			})
			rr := httptest.NewRecorder()
			err := got(rr, tt.args.request)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("got error %v, want %v", err, tt.res.wantErr)
			}

			assert.Equal(t, tt.res.context, ctx)
		})
	}
}

func Test_setInstance(t *testing.T) {
	type args struct {
		ctx      context.Context
		verifier authz.InstanceVerifier
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
			"no domain context, not found error",
			args{
				ctx:      context.Background(),
				verifier: &mockInstanceVerifier{},
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"instanceHost found, ok",
			args{
				ctx:      zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "host", Protocol: "https"}),
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			res{
				want: authz.WithInstance(zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "host", Protocol: "https"}), &mockInstance{}),
				err:  false,
			},
		},
		{
			"instanceHost not found, error",
			args{
				ctx:      zitadel_http.WithDomainContext(context.Background(), &zitadel_http.DomainCtx{InstanceHost: "fromorigin:9999", Protocol: "https"}),
				verifier: &mockInstanceVerifier{instanceHost: "unknowndomain"},
			},
			res{
				want: nil,
				err:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setInstance(tt.args.ctx, tt.args.verifier)
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
	instanceHost string
	publicHost   string
}

func (m *mockInstanceVerifier) InstanceByHost(_ context.Context, instanceHost, publicHost string) (authz.Instance, error) {
	if instanceHost != m.instanceHost {
		return nil, fmt.Errorf("invalid host")
	}
	if publicHost == "" {
		return &mockInstance{}, nil
	}
	if publicHost != instanceHost && publicHost != m.publicHost {
		return nil, fmt.Errorf("invalid host")
	}
	return &mockInstance{}, nil
}

func (m *mockInstanceVerifier) InstanceByID(context.Context, string) (authz.Instance, error) {
	return nil, nil
}

type mockInstance struct{}

func (m *mockInstance) Block() *bool {
	panic("shouldn't be called here")
}

func (m *mockInstance) AuditLogRetention() *time.Duration {
	panic("shouldn't be called here")
}

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

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func (m *mockInstance) EnableImpersonation() bool {
	return false
}

func (m *mockInstance) Features() feature.Features {
	return feature.Features{}
}

func (m *mockInstance) ExecutionRouter() target.Router {
	return target.NewRouter(nil)
}
