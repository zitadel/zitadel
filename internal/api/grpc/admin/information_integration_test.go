//go:build integration

package admin_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

var (
	Tester *integration.Tester
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(time.Minute)
		defer cancel()
		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		return m.Run()
	}())
}

func TestServer_Healthz(t *testing.T) {
	_, err := Tester.Client.Admin.Healthz(context.TODO(), &admin.HealthzRequest{})
	require.NoError(t, err)
}
