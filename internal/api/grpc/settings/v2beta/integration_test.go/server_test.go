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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Instance, err = integration.FirstInstance(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		AdminCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Instance.Client.SettingsV2beta
		return m.Run()
	}())
}
