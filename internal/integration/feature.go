package integration

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

func EnsureInstanceFeature(t *testing.T, ctx context.Context, instance *Instance, features *feature.SetInstanceFeaturesRequest) {
	ctx = instance.WithAuthorizationToken(ctx, UserTypeIAMOwner)
	_, err := instance.Client.FeatureV2.SetInstanceFeatures(ctx, features)
	require.NoError(t, err)
}

type RelationalTableFeatureMatrix struct {
	Name      string
	Inst      *Instance
	InstOwner context.Context
	Enabled   bool
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
			EnsureInstanceFeature(t, ctx, inst, &feature.SetInstanceFeaturesRequest{
				EnableRelationalTables: gu.Ptr(true),
			})
			return RelationalTableFeatureMatrix{
				Name:      "when relational tables are enabled",
				Inst:      inst,
				InstOwner: instOwner,
				Enabled:   true,
			}
		}(),
	}
}
