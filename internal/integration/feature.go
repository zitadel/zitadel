package integration

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
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

type RelationalTableFeatureMatrix struct {
	Name      string
	Inst      *Instance
	InstOwner context.Context
}

// TODO(IAM-Marco): Once we have gotten rid of eventstore, this can be removed
func RelationalTablesEnableMatrix(t *testing.T, ctx context.Context, sysAuthZ context.Context) []RelationalTableFeatureMatrix {
	return []RelationalTableFeatureMatrix{
		func() RelationalTableFeatureMatrix {
			inst := NewInstance(sysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, UserTypeIAMOwner)
			return RelationalTableFeatureMatrix{
				Name:      "when relational tables are disabled",
				Inst:      inst,
				InstOwner: instOwner,
			}
		}(),
		func() RelationalTableFeatureMatrix {
			inst := NewInstance(sysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, UserTypeIAMOwner)
			_, err := inst.Client.FeatureV2.SetInstanceFeatures(instOwner, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(false)})
			require.NoError(t, err)
			return RelationalTableFeatureMatrix{
				Name:      "when relational tables are enabled",
				Inst:      inst,
				InstOwner: instOwner,
			}
		}(),
	}
}
