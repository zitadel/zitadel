//go:build integration

package settings_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

var (
	CTX, AdminCTX context.Context
	Tester        *integration.Tester
	Client        settings.SettingsServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(3 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX = ctx
		AdminCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		Client = Tester.Client.SettingsV2beta
		return m.Run()
	}())
}
