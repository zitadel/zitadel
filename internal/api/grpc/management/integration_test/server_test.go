//go:build integration

package management_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

var (
	CTX, IAMOwnerCTX, OrgCTX context.Context
	Instance                 *integration.Instance
	Client                   mgmt_pb.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		CTX = ctx
		IAMOwnerCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		OrgCTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		Client = Instance.Client.Mgmt
		return m.Run()
	}())
}

// ensureOrgEmailProviders turns on the org_email_providers instance feature,
// which is off by default and gates every org email provider endpoint.
func ensureOrgEmailProviders(t *testing.T) {
	t.Helper()
	enabled := true
	integration.EnsureInstanceFeature(t, CTX, Instance,
		&feature.SetInstanceFeaturesRequest{OrgEmailProviders: &enabled},
		func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.True(tt, got.GetOrgEmailProviders().GetEnabled())
		},
	)
}
