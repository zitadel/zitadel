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
	Instance      *integration.Instance
	Client        settings.SettingsServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = ctx
		AdminCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Instance.Client.SettingsV2beta
		return m.Run()
	}())
}
