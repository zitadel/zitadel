package integration

import (
	"context"
	_ "embed"
	"sync"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

var (
	//go:embed config/system-user-key.pem
	systemUserKey []byte
	//go:embed config/system-user-with-no-permissions.pem
	systemUserWithNoPermissions []byte
)

var (
	// SystemClient creates a system connection once and reuses it on every use.
	// Each client call automatically gets the authorization context for the system user.
	SystemClient                     = sync.OnceValue[system.SystemServiceClient](systemClient)
	SystemToken                      string
	SystemUserWithNoPermissionsToken string
)

func systemClient() system.SystemServiceClient {
	cc, err := grpc.NewClient(loadedConfig.Host(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = WithSystemAuthorization(ctx)
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		panic(err)
	}
	return system.NewSystemServiceClient(cc)
}

func createSystemUserToken() string {
	const ISSUER = "tester"
	audience := http_util.BuildOrigin(loadedConfig.Host(), loadedConfig.Secure)
	signer, err := client.NewSignerFromPrivateKeyByte(systemUserKey, "")
	if err != nil {
		panic(err)
	}
	token, err := client.SignedJWTProfileAssertion(ISSUER, []string{audience}, time.Hour, signer)
	if err != nil {
		panic(err)
	}
	return token
}

func createSystemUserWithNoPermissionsToken() string {
	const ISSUER = "system-user-with-no-permissions"
	audience := http_util.BuildOrigin(loadedConfig.Host(), loadedConfig.Secure)
	signer, err := client.NewSignerFromPrivateKeyByte(systemUserWithNoPermissions, "")
	if err != nil {
		panic(err)
	}
	token, err := client.SignedJWTProfileAssertion(ISSUER, []string{audience}, time.Hour, signer)
	if err != nil {
		panic(err)
	}
	return token
}

func WithSystemAuthorization(ctx context.Context) context.Context {
	return WithAuthorizationToken(ctx, SystemToken)
}

func WithSystemUserWithNoPermissionsAuthorization(ctx context.Context) context.Context {
	return WithAuthorizationToken(ctx, SystemUserWithNoPermissionsToken)
}
