package middleware

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_v3 "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
)

func Test_setInstance(t *testing.T) {
	type args struct {
		ctx      context.Context
		req      interface{}
		info     *grpc.UnaryServerInfo
		handler  grpc.UnaryHandler
		verifier authz.InstanceVerifier
	}
	type res struct {
		want interface{}
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
				ctx: context.Background(),
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"invalid host, error",
			args{
				ctx:      http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host2"}),
				req:      &mockRequest{},
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"valid host",
			args{
				ctx:      http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host"}),
				req:      &mockRequest{},
				verifier: &mockInstanceVerifier{instanceHost: "host"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return req, nil
				},
			},
			res{
				want: &mockRequest{},
				err:  false,
			},
		},
		{
			"explicit instance unset, hostname not found, error",
			args{
				ctx:      http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host2"}),
				req:      &mockRequestWithExplicitInstance{},
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"explicit instance unset, invalid host, error",
			args{
				ctx:      http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host2"}),
				req:      &mockRequestWithExplicitInstance{},
				verifier: &mockInstanceVerifier{instanceHost: "host"},
			},
			res{
				want: nil,
				err:  true,
			},
		},
		{
			"explicit instance unset, valid host",
			args{
				ctx:      http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host"}),
				req:      &mockRequestWithExplicitInstance{},
				verifier: &mockInstanceVerifier{instanceHost: "host"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return req, nil
				},
			},
			res{
				want: &mockRequestWithExplicitInstance{},
				err:  false,
			},
		},
		{
			name: "explicit instance set, id not found, error",
			args: args{
				ctx: context.Background(),
				req: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Id{
							Id: "not existing instance id",
						},
					},
				},
				verifier: &mockInstanceVerifier{id: "existing instance id"},
			},
			res: res{
				want: nil,
				err:  true,
			},
		},
		{
			name: "explicit instance set, id found, ok",
			args: args{
				ctx: context.Background(),
				req: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Id{
							Id: "existing instance id",
						},
					},
				},
				verifier: &mockInstanceVerifier{id: "existing instance id"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return req, nil
				},
			},
			res: res{
				want: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Id{
							Id: "existing instance id",
						},
					},
				},
				err: false,
			},
		},
		{
			name: "explicit instance set, domain not found, error",
			args: args{
				ctx: context.Background(),
				req: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Domain{
							Domain: "not existing instance domain",
						},
					},
				},
				verifier: &mockInstanceVerifier{instanceHost: "existing instance domain"},
			},
			res: res{
				want: nil,
				err:  true,
			},
		},
		{
			name: "explicit instance set, domain found, ok",
			args: args{
				ctx: context.Background(),
				req: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Domain{
							Domain: "existing instance domain",
						},
					},
				},
				verifier: &mockInstanceVerifier{instanceHost: "existing instance domain"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return req, nil
				},
			},
			res: res{
				want: &mockRequestWithExplicitInstance{
					instance: object_v3.Instance{
						Property: &object_v3.Instance_Domain{
							Domain: "existing instance domain",
						},
					},
				},
				err: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setInstance(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler, tt.args.verifier, "", nil)
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

type mockRequest struct{}

type mockRequestWithExplicitInstance struct {
	instance object_v3.Instance
}

func (m *mockRequestWithExplicitInstance) GetInstance() *object_v3.Instance {
	return &m.instance
}

type mockInstanceVerifier struct {
	id           string
	instanceHost string
	publicHost   string
	err          error
}

func (m *mockInstanceVerifier) InstanceByHost(_ context.Context, instanceHost, publicHost string) (authz.Instance, error) {
	if m.err != nil {
		return nil, m.err
	}
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

func (m *mockInstanceVerifier) InstanceByID(_ context.Context, id string) (authz.Instance, error) {
	if m.err != nil {
		return nil, m.err
	}
	if id != m.id {
		return nil, fmt.Errorf("not found")
	}
	return &mockInstance{}, nil
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

func (m *mockInstance) ManagementConsoleClientID() string {
	return "consoleClientID"
}

func (m *mockInstance) ManagementConsoleApplicationID() string {
	return "consoleApplicationID"
}

func (m *mockInstance) DefaultLanguage() language.Tag {
	return language.English
}

func (m *mockInstance) AllowedLanguages() []language.Tag {
	return []language.Tag{language.English}
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

func Test_setInstance_errorCodes(t *testing.T) {
	i18n.SupportLanguages(language.English)
	translator := i18n.NewZitadelTranslator(language.English)

	cases := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{
			name:     "not found from verifier propagates as NotFound",
			err:      zerrors.ThrowNotFound(nil, "TEST-001", "Errors.Instance.NotFound"),
			wantCode: codes.NotFound,
		},
		{
			name:     "internal error from verifier propagates as Internal",
			err:      zerrors.ThrowInternal(errors.New("FATAL: the database system is shutting down (SQLSTATE 57P03)"), "TEST-002", "Errors.Internal"),
			wantCode: codes.Internal,
		},
		{
			name:     "unavailable error from verifier propagates as Unavailable",
			err:      zerrors.ThrowUnavailable(nil, "TEST-003", "Errors.Unavailable"),
			wantCode: codes.Unavailable,
		},
	}

	for _, tc := range cases {
		verifier := &mockInstanceVerifier{err: tc.err}

		t.Run("byRequestedHost/"+tc.name, func(t *testing.T) {
			ctx := http_util.WithDomainContext(context.Background(), &http_util.DomainCtx{InstanceHost: "host"})
			_, err := setInstance(ctx, &mockRequest{}, &grpc.UnaryServerInfo{}, nil, verifier, "", translator)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := status.Code(err); got != tc.wantCode {
				t.Errorf("got code %v, want %v", got, tc.wantCode)
			}
		})

		t.Run("byID/"+tc.name, func(t *testing.T) {
			req := &mockRequestWithExplicitInstance{
				instance: object_v3.Instance{
					Property: &object_v3.Instance_Id{Id: "any"},
				},
			}
			_, err := setInstance(context.Background(), req, &grpc.UnaryServerInfo{}, nil, verifier, "", translator)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := status.Code(err); got != tc.wantCode {
				t.Errorf("got code %v, want %v", got, tc.wantCode)
			}
		})

		t.Run("byDomain/"+tc.name, func(t *testing.T) {
			req := &mockRequestWithExplicitInstance{
				instance: object_v3.Instance{
					Property: &object_v3.Instance_Domain{Domain: "any"},
				},
			}
			_, err := setInstance(context.Background(), req, &grpc.UnaryServerInfo{}, nil, verifier, "", translator)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if got := status.Code(err); got != tc.wantCode {
				t.Errorf("got code %v, want %v", got, tc.wantCode)
			}
		})
	}
}
