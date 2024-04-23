//go:build integration

package execution_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client execution.ExecutionServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.ExecutionV3

		CTX, _ = Tester.WithAuthorization(ctx, integration.IAMOwner), errCtx
		return m.Run()
	}())
}
