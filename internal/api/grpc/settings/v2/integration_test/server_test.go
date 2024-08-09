//go:build integration

package settings_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

var (
	CTX, AdminCTX context.Context
	Instance      *integration.Instance
	Client        settings.SettingsServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		first, err := integration.FirstInstance(ctx)
		if err != nil {
			panic(err)
		}
		Instance, err = first.UseIsolatedInstance(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		AdminCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Instance.Client.SettingsV2
		return m.Run()
	}())
}
