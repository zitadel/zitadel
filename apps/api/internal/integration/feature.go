package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

func EnsureInstanceFeature(t *testing.T, ctx context.Context, instance *Instance, features *feature.SetInstanceFeaturesRequest, assertFeatures func(t *assert.CollectT, got *feature.GetInstanceFeaturesResponse)) {
	ctx = instance.WithAuthorizationToken(ctx, UserTypeIAMOwner)
	_, err := instance.Client.FeatureV2.SetInstanceFeatures(ctx, features)
	require.NoError(t, err)
	retryDuration, tick := WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t,
		func(tt *assert.CollectT) {
			got, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(tt, err)
			assertFeatures(tt, got)
		},
		retryDuration,
		tick,
		"timed out waiting for ensuring instance feature")
}
