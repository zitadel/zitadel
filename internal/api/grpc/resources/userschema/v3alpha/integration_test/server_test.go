//go:build integration

package userschema_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	schema "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
)

var (
	IAMOwnerCTX, SystemCTX context.Context
	Instance               *integration.Instance
	Client                 schema.ZITADELUserSchemasClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		IAMOwnerCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		SystemCTX = integration.WithSystemAuthorization(ctx)
		Client = Instance.Client.UserSchemaV3

		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T, iamOwnerCTX context.Context) {
	f, err := Instance.Client.FeatureV2.GetInstanceFeatures(iamOwnerCTX, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.UserSchema.GetEnabled() {
		return
	}
	_, err = Instance.Client.FeatureV2.SetInstanceFeatures(iamOwnerCTX, &feature.SetInstanceFeaturesRequest{
		UserSchema: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := iamOwnerCTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := Instance.Client.FeatureV2.GetInstanceFeatures(iamOwnerCTX, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(ttt, err)
			if f.UserSchema.GetEnabled() {
				return
			}
		},
		retryDuration,
		time.Second,
		"timed out waiting for ensuring instance feature")
}
