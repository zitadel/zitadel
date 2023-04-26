package middleware

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func Test_hostNameFromContext(t *testing.T) {
	type args struct {
		ctx        context.Context
		headerName string
	}
	type res struct {
		want string
		err  bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"empty context, error",
			args{
				ctx:        context.Background(),
				headerName: "header",
			},
			res{
				want: "",
				err:  true,
			},
		},
		{
			"header not found",
			args{
				ctx:        metadata.NewIncomingContext(context.Background(), nil),
				headerName: "header",
			},
			res{
				want: "",
				err:  true,
			},
		},
		{
			"header not found",
			args{
				ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("header", "value")),
				headerName: "header",
			},
			res{
				want: "value",
				err:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hostFromContext(tt.args.ctx, tt.args.headerName)
			if (err != nil) != tt.res.err {
				t.Errorf("hostFromContext() error = %v, wantErr %v", err, tt.res.err)
				return
			}
			if got != tt.res.want {
				t.Errorf("hostFromContext() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}

func Test_setInstance(t *testing.T) {
	type args struct {
		ctx        context.Context
		req        interface{}
		info       *grpc.UnaryServerInfo
		handler    grpc.UnaryHandler
		verifier   authz.InstanceVerifier
		headerName string
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
				ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("header", "host2")),
				req:        &mockRequest{},
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
				ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("header", "host")),
				req:        &mockRequest{},
				verifier:   &mockInstanceVerifier{"host"},
				headerName: "header",
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
			got, err := setInstance(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler, tt.args.verifier, tt.args.headerName, nil)
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
	host string
}

func (m *mockInstanceVerifier) InstanceByHost(_ context.Context, host string) (authz.Instance, error) {
	if host != m.host {
		return nil, fmt.Errorf("invalid host")
	}
	return &mockInstance{}, nil
}

func (m *mockInstanceVerifier) InstanceByID(context.Context) (authz.Instance, error) { return nil, nil }

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
	return "localhost"
}

func (m *mockInstance) RequestedHost() string {
	return "localhost:8080"
}

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}
