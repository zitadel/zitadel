//go:build integration

package system_test

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_Limits_Bulk(t *testing.T) {
	const len = 5
	type instance struct{ domain, id string }
	instances := make([]*instance, len)
	for i := 0; i < len; i++ {
		domain := integration.RandString(5) + ".integration.localhost"
		resp, err := integration.SystemClient().CreateInstance(CTX, &system.CreateInstanceRequest{
			InstanceName: "testinstance",
			CustomDomain: domain,
			Owner: &system.CreateInstanceRequest_Machine_{
				Machine: &system.CreateInstanceRequest_Machine{
					UserName: "owner",
					Name:     "owner",
				},
			},
		})
		require.NoError(t, err)
		instances[i] = &instance{domain, resp.GetInstanceId()}
	}
	resp, err := integration.SystemClient().BulkSetLimits(CTX, &system.BulkSetLimitsRequest{
		Limits: []*system.SetLimitsRequest{{
			InstanceId: instances[0].id,
			Block:      gu.Ptr(true),
		}, {
			InstanceId: instances[1].id,
			Block:      gu.Ptr(false),
		}, {
			InstanceId: instances[2].id,
			Block:      gu.Ptr(true),
		}, {
			InstanceId: instances[3].id,
			Block:      gu.Ptr(false),
		}, {
			InstanceId: instances[4].id,
			Block:      gu.Ptr(true),
		}},
	})
	require.NoError(t, err)
	details := resp.GetTargetDetails()
	require.Len(t, details, len)
	t.Run("the first instance is blocked", func(t *testing.T) {
		require.Equal(t, instances[0].id, details[0].GetResourceOwner(), "resource owner must be instance id")
		testBlockingAPI(t, publicAPIBlockingTest(instances[0].domain), true, true)
	})
	t.Run("the second instance isn't blocked", func(t *testing.T) {
		require.Equal(t, instances[1].id, details[1].GetResourceOwner(), "resource owner must be instance id")
		testBlockingAPI(t, publicAPIBlockingTest(instances[1].domain), false, true)
	})
	t.Run("the third instance is blocked", func(t *testing.T) {
		require.Equal(t, instances[2].id, details[2].GetResourceOwner(), "resource owner must be instance id")
		testBlockingAPI(t, publicAPIBlockingTest(instances[2].domain), true, true)
	})
	t.Run("the fourth instance isn't blocked", func(t *testing.T) {
		require.Equal(t, instances[3].id, details[3].GetResourceOwner(), "resource owner must be instance id")
		testBlockingAPI(t, publicAPIBlockingTest(instances[3].domain), false, true)
	})
	t.Run("the fifth instance is blocked", func(t *testing.T) {
		require.Equal(t, instances[4].id, details[4].GetResourceOwner(), "resource owner must be instance id")
		testBlockingAPI(t, publicAPIBlockingTest(instances[4].domain), true, true)
	})
}
