//go:build integration

package idp_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
)

var (
	CTX       context.Context
	IamCTX    context.Context
	UserCTX   context.Context
	SystemCTX context.Context
	ErrCTX    context.Context
	Tester    *integration.Tester
	Client    idp_pb.IDPServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		UserCTX = Tester.WithAuthorization(ctx, integration.Login)
		IamCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		CTX, ErrCTX = Tester.WithAuthorization(ctx, integration.OrgOwner), errCtx
		Client = Tester.Client.IDPv2
		return m.Run()
	}())
}
