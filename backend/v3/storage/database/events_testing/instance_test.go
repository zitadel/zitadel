//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TestInstanceReduces(t *testing.T) {
	instanceRepo := repository.InstanceRepository(pool)

	t.Run("test instance add reduces", func(t *testing.T) {
		instanceName := gofakeit.Name()
		beforeCreate := time.Now()
		instance, err := SystemClient.CreateInstance(CTX, &system.CreateInstanceRequest{
			InstanceName: instanceName,
			Owner: &system.CreateInstanceRequest_Machine_{
				Machine: &system.CreateInstanceRequest_Machine{
					UserName:            "owner",
					Name:                "owner",
					PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
				},
			},
		})
		afterCreate := time.Now()

		require.NoError(t, err)
		t.Cleanup(func() {
			_, err = SystemClient.RemoveInstance(CTX, &system.RemoveInstanceRequest{
				InstanceId: instance.GetInstanceId(),
			})
			if err != nil {
				t.Logf("Failed to delete instance on cleanup: %v", err)
			}
		})

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			instance, err := instanceRepo.Get(CTX,
				database.WithCondition(instanceRepo.IDCondition(instance.GetInstanceId())),
			)
			require.NoError(t, err)
			// event instance.added
			assert.Equal(t, instanceName, instance.Name)
			// event instance.default.org.set
			assert.NotNil(t, instance.DefaultOrgID)
			// event instance.iam.project.set
			assert.NotNil(t, instance.IAMProjectID)
			// event instance.iam.console.set
			assert.NotNil(t, instance.ConsoleAppID)
			// event instance.default.language.set
			assert.NotNil(t, instance.DefaultLanguage)
			// event instance.added
			assert.WithinRange(t, instance.CreatedAt, beforeCreate, afterCreate)
			// event instance.added
			assert.WithinRange(t, instance.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})

	t.Run("test instance update reduces", func(t *testing.T) {
		instanceName := gofakeit.Name()
		res, err := SystemClient.CreateInstance(CTX, &system.CreateInstanceRequest{
			InstanceName: instanceName,
			Owner: &system.CreateInstanceRequest_Machine_{
				Machine: &system.CreateInstanceRequest_Machine{
					UserName:            "owner",
					Name:                "owner",
					PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err = SystemClient.RemoveInstance(CTX, &system.RemoveInstanceRequest{
				InstanceId: res.GetInstanceId(),
			})
			if err != nil {
				t.Logf("Failed to delete instance on cleanup: %v", err)
			}
		})

		// check instance exists
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			instance, err := instanceRepo.Get(CTX,
				database.WithCondition(instanceRepo.IDCondition(res.GetInstanceId())),
			)
			require.NoError(t, err)
			assert.Equal(t, instanceName, instance.Name)
		}, retryDuration, tick)

		instanceName += "new"
		beforeUpdate := time.Now()
		_, err = SystemClient.UpdateInstance(CTX, &system.UpdateInstanceRequest{
			InstanceId:   res.InstanceId,
			InstanceName: instanceName,
		})
		afterUpdate := time.Now()
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			instance, err := instanceRepo.Get(CTX,
				database.WithCondition(instanceRepo.IDCondition(res.GetInstanceId())),
			)
			require.NoError(t, err)
			// event instance.changed
			assert.Equal(t, instanceName, instance.Name)
			assert.WithinRange(t, instance.UpdatedAt, beforeUpdate, afterUpdate)
		}, retryDuration, tick)
	})

	t.Run("test instance delete reduces", func(t *testing.T) {
		instanceName := gofakeit.Name()
		res, err := SystemClient.CreateInstance(CTX, &system.CreateInstanceRequest{
			InstanceName: instanceName,
			Owner: &system.CreateInstanceRequest_Machine_{
				Machine: &system.CreateInstanceRequest_Machine{
					UserName:            "owner",
					Name:                "owner",
					PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
				},
			},
		})
		require.NoError(t, err)

		// check instance exists
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			instance, err := instanceRepo.Get(CTX,
				database.WithCondition(instanceRepo.IDCondition(res.GetInstanceId())),
			)
			require.NoError(t, err)
			assert.Equal(t, instanceName, instance.Name)
		}, retryDuration, tick)

		_, err = SystemClient.RemoveInstance(CTX, &system.RemoveInstanceRequest{
			InstanceId: res.InstanceId,
		})
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			instance, err := instanceRepo.Get(CTX,
				database.WithCondition(instanceRepo.IDCondition(res.GetInstanceId())),
			)
			// event instance.removed
			assert.Nil(t, instance)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
