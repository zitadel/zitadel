//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TestInstanceDomainReduces(t *testing.T) {
	instance := integration.NewInstance(CTX)

	instanceRepo := repository.InstanceRepository(pool)
	instanceDomainRepo := instanceRepo.Domains(true)

	t.Cleanup(func() {
		_, err := instance.Client.InstanceV2Beta.DeleteInstance(CTX, &v2beta.DeleteInstanceRequest{
			InstanceId: instance.Instance.Id,
		})
		if err != nil {
			t.Logf("Failed to delete instance on cleanup: %v", err)
		}
	})

	// Wait for instance to be created
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := instanceRepo.Get(CTX,
			database.WithCondition(instanceRepo.IDCondition(instance.Instance.Id)),
		)
		assert.NoError(ttt, err)
	}, retryDuration, tick)

	t.Run("test instance custom domain add reduces", func(t *testing.T) {
		// Add a domain to the instance
		domainName := gofakeit.DomainName()
		beforeAdd := time.Now()
		_, err := instance.Client.InstanceV2Beta.AddCustomDomain(CTX, &v2beta.AddCustomDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		t.Cleanup(func() {
			_, err := instance.Client.InstanceV2Beta.RemoveCustomDomain(CTX, &v2beta.RemoveCustomDomainRequest{
				InstanceId: instance.Instance.Id,
				Domain:     domainName,
			})
			if err != nil {
				t.Logf("Failed to delete instance domain on cleanup: %v", err)
			}
		})

		// Test that domain add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					),
				),
			)
			require.NoError(ttt, err)
			// event instance.domain.added
			assert.Equal(ttt, domainName, domain.Domain)
			assert.Equal(ttt, instance.Instance.Id, domain.InstanceID)
			assert.False(ttt, *domain.IsPrimary)
			assert.WithinRange(ttt, domain.CreatedAt, beforeAdd, afterAdd)
			assert.WithinRange(ttt, domain.UpdatedAt, beforeAdd, afterAdd)
		}, retryDuration, tick)
	})

	t.Run("test instance custom domain set primary reduces", func(t *testing.T) {
		// Add a domain to the instance
		domainName := gofakeit.DomainName()
		_, err := instance.Client.InstanceV2Beta.AddCustomDomain(CTX, &v2beta.AddCustomDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)

		t.Cleanup(func() {
			// first we change the primary domain to something else
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
						instanceDomainRepo.IsPrimaryCondition(false),
					),
				),
				database.WithLimit(1),
			)
			require.NoError(t, err)
			_, err = SystemClient.SetPrimaryDomain(CTX, &system.SetPrimaryDomainRequest{
				InstanceId: instance.Instance.Id,
				Domain:     domain.Domain,
			})
			require.NoError(t, err)

			_, err = instance.Client.InstanceV2Beta.RemoveCustomDomain(CTX, &v2beta.RemoveCustomDomainRequest{
				InstanceId: instance.Instance.Id,
				Domain:     domainName,
			})
			if err != nil {
				t.Logf("Failed to delete instance domain on cleanup: %v", err)
			}
		})

		// Wait for domain to be created
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					),
				),
			)
			require.NoError(ttt, err)
			require.False(ttt, *domain.IsPrimary)
			assert.Equal(ttt, domainName, domain.Domain)
		}, retryDuration, tick)

		// Set domain as primary
		beforeSetPrimary := time.Now()
		_, err = SystemClient.SetPrimaryDomain(CTX, &system.SetPrimaryDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)
		afterSetPrimary := time.Now()

		// Test that set primary reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.IsPrimaryCondition(true),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					),
				),
			)
			require.NoError(ttt, err)
			// event instance.domain.primary.set
			assert.Equal(ttt, domainName, domain.Domain)
			assert.True(ttt, *domain.IsPrimary)
			assert.WithinRange(ttt, domain.UpdatedAt, beforeSetPrimary, afterSetPrimary)
		}, retryDuration, tick)
	})

	t.Run("test instance custom domain remove reduces", func(t *testing.T) {
		// Add a domain to the instance
		domainName := gofakeit.DomainName()
		_, err := instance.Client.InstanceV2Beta.AddCustomDomain(CTX, &v2beta.AddCustomDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)

		// Wait for domain to be created and verify it exists
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			_, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					),
				),
			)
			require.NoError(ttt, err)
		}, retryDuration, tick)

		// Remove the domain
		_, err = instance.Client.InstanceV2Beta.RemoveCustomDomain(CTX, &v2beta.RemoveCustomDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)

		// Test that domain remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeCustom),
					),
				),
			)
			// event instance.domain.removed
			assert.Nil(ttt, domain)
			require.ErrorIs(ttt, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})

	t.Run("test instance trusted domain add reduces", func(t *testing.T) {
		// Add a domain to the instance
		domainName := gofakeit.DomainName()
		beforeAdd := time.Now()
		_, err := instance.Client.InstanceV2Beta.AddTrustedDomain(CTX, &v2beta.AddTrustedDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		t.Cleanup(func() {
			_, err := instance.Client.InstanceV2Beta.RemoveTrustedDomain(CTX, &v2beta.RemoveTrustedDomainRequest{
				InstanceId: instance.Instance.Id,
				Domain:     domainName,
			})
			if err != nil {
				t.Logf("Failed to delete instance domain on cleanup: %v", err)
			}
		})

		// Test that domain add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeTrusted),
					),
				),
			)
			require.NoError(ttt, err)
			// event instance.domain.added
			assert.Equal(ttt, domainName, domain.Domain)
			assert.Equal(ttt, instance.Instance.Id, domain.InstanceID)
			assert.WithinRange(ttt, domain.CreatedAt, beforeAdd, afterAdd)
			assert.WithinRange(ttt, domain.UpdatedAt, beforeAdd, afterAdd)
		}, retryDuration, tick)
	})

	t.Run("test instance trusted domain remove reduces", func(t *testing.T) {
		// Add a domain to the instance
		domainName := gofakeit.DomainName()
		_, err := instance.Client.InstanceV2Beta.AddTrustedDomain(CTX, &v2beta.AddTrustedDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)

		// Wait for domain to be created and verify it exists
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			_, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeTrusted),
					),
				),
			)
			require.NoError(ttt, err)
		}, retryDuration, tick)

		// Remove the domain
		_, err = instance.Client.InstanceV2Beta.RemoveTrustedDomain(CTX, &v2beta.RemoveTrustedDomainRequest{
			InstanceId: instance.Instance.Id,
			Domain:     domainName,
		})
		require.NoError(t, err)

		// Test that domain remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := instanceDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						instanceDomainRepo.InstanceIDCondition(instance.Instance.Id),
						instanceDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
						instanceDomainRepo.TypeCondition(domain.DomainTypeTrusted),
					),
				),
			)
			// event instance.domain.removed
			assert.Nil(ttt, domain)
			require.ErrorIs(ttt, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
