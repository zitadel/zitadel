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

func (v *verifierMock) VerifyAccessToken(ctx context.Context, token string) (string, string, string, error) {
	return "", "", "", nil
}
func (v *verifierMock) ResolveGrant(ctx context.Context) (*authz.Grant, error) {
	return nil, nil
}
func (v *verifierMock) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return "", nil
}

func Test_authorize(t *testing.T) {
	type args struct {
		ctx         context.Context
		req         interface{}
		info        *grpc.UnaryServerInfo
		handler     grpc.UnaryHandler
		verifier    authz.TokenVerifierOld
		authConfig  *authz.Config
		authMethods authz.MethodMapping
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"no token needed ok",
			args{
				ctx:         context.Background(),
				req:         &mockReq{},
				info:        mockInfo("no.token.needed"),
				handler:     emptyMockHandler,
				verifier:    nil,
				authConfig:  nil,
				authMethods: mockMethods,
			},
			&mockReq{},
			false,
		},
		{
			"auth header missing error",
			args{
				ctx:         context.Background(),
				req:         &mockReq{},
				info:        mockInfo("need.authentication"),
				handler:     emptyMockHandler,
				verifier:    nil,
				authConfig:  nil,
				authMethods: mockMethods,
			},
			nil,
			true,
		},
		{
			"unauthorized error",
			args{
				ctx:         metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "wrong")),
				req:         &mockReq{},
				info:        mockInfo("need.authentication"),
				handler:     emptyMockHandler,
				verifier:    nil,
				authConfig:  nil,
				authMethods: mockMethods,
			},
			nil,
			true,
		},
		{
			"authorized ok",
			args{
				ctx:         metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer token")),
				req:         &mockReq{},
				info:        mockInfo("need.authentication"),
				handler:     emptyMockHandler,
				verifier:    &verifierMock{},
				authConfig:  nil,
				authMethods: mockMethods,
			},
			&mockReq{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorize(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler, tt.args.verifier, tt.args.authConfig, tt.args.authMethods)
			if (err != nil) != tt.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authorize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
