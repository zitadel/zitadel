//go:build integration

package system_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

var (
	CTX       context.Context
	SystemCTX context.Context
	ErrCTX    context.Context
	Tester    *integration.Tester
	Client    system.SystemServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()
		CTX = ctx

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		SystemCTX, ErrCTX = Tester.WithAuthorization(ctx, integration.SystemUser), errCtx
		Client = Tester.Client.System
		return m.Run()
	}())
}
