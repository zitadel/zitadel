//go:build integration

package user_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v3alpha"
)

var (
	CTX       context.Context
	IamCTX    context.Context
	UserCTX   context.Context
	SystemCTX context.Context
	ErrCTX    context.Context
	Tester    *integration.Tester
	Client    user.ZITADELUsersClient
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
		Client = Tester.Client.UserV3Alpha
		return m.Run()
	}())
}
