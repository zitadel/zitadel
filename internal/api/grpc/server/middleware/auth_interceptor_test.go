package middleware

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/caos/zitadel/internal/api/authz"
)

var (
	mockMethods = authz.MethodMapping{
		"need.authentication": authz.Option{
			Permission: "authenticated",
		},
	}
)

type verifierMock struct{}

func (v *verifierMock) VerifyAccessToken(ctx context.Context, token, clientID string) (string, string, string, string, error) {
	return "", "", "", "", nil
}
func (v *verifierMock) SearchMyMemberships(ctx context.Context) ([]*authz.Membership, error) {
	return nil, nil
}

func (v *verifierMock) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (string, []string, error) {
	return "", nil, nil
}
func (v *verifierMock) ExistsOrg(ctx context.Context, orgID string) error {
	return nil
}
func (v *verifierMock) VerifierClientID(ctx context.Context, appName string) (string, error) {
	return "", nil
}

func Test_authorize(t *testing.T) {
	type args struct {
		ctx         context.Context
		req         interface{}
		info        *grpc.UnaryServerInfo
		handler     grpc.UnaryHandler
		verifier    *authz.TokenVerifier
		authConfig  authz.Config
		authMethods authz.MethodMapping
	}
	type res struct {
		want    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no token needed ok",
			args{
				ctx:     context.Background(),
				req:     &mockReq{},
				info:    mockInfo("/no/token/needed"),
				handler: emptyMockHandler,
				verifier: func() *authz.TokenVerifier {
					verifier := authz.Start(&verifierMock{})
					verifier.RegisterServer("need", "need", authz.MethodMapping{})
					return verifier
				}(),
				authMethods: mockMethods,
			},
			res{
				&mockReq{},
				false,
			},
		},
		{
			"auth header missing error",
			args{
				ctx:     context.Background(),
				req:     &mockReq{},
				info:    mockInfo("/need/authentication"),
				handler: emptyMockHandler,
				verifier: func() *authz.TokenVerifier {
					verifier := authz.Start(&verifierMock{})
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				}(),
				authConfig:  authz.Config{},
				authMethods: mockMethods,
			},
			res{
				nil,
				true,
			},
		},
		{
			"unauthorized error",
			args{
				ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "wrong")),
				req:     &mockReq{},
				info:    mockInfo("/need/authentication"),
				handler: emptyMockHandler,
				verifier: func() *authz.TokenVerifier {
					verifier := authz.Start(&verifierMock{})
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				}(),
				authConfig:  authz.Config{},
				authMethods: mockMethods,
			},
			res{
				nil,
				true,
			},
		},
		{
			"authorized ok",
			args{
				ctx:     metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer token")),
				req:     &mockReq{},
				info:    mockInfo("/need/authentication"),
				handler: emptyMockHandler,
				verifier: func() *authz.TokenVerifier {
					verifier := authz.Start(&verifierMock{})
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				}(),
				authConfig:  authz.Config{},
				authMethods: mockMethods,
			},
			res{
				&mockReq{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorize(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler, tt.args.verifier, tt.args.authConfig)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("authorize() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}
