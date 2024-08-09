package middleware

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"golang.org/x/text/language"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/feature"
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

func (m *mockInstanceVerifier) InstanceByID(context.Context) (authz.Instance, error) { return nil, nil }

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
