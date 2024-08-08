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
	Tester        *integration.Tester
	Client        settings.SettingsServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Tester, err = integration.NewTester(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		AdminCTX = Tester.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Tester.Client.SettingsV2
		return m.Run()
	}())
}
