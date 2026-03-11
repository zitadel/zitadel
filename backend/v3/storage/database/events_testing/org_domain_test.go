//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	v2 "github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestServer_TestOrgDomainReduces(t *testing.T) {
	org, err := OrgClient.AddOrganization(CTX, &v2.AddOrganizationRequest{
		Name: gofakeit.Name(),
	})
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository()
	orgDomainRepo := repository.OrganizationDomainRepository()

	t.Cleanup(func() {
		_, err := OrgClient.DeleteOrganization(CTX, &v2.DeleteOrganizationRequest{
			OrganizationId: org.GetOrganizationId(),
		})
		if err != nil {
			t.Logf("Failed to delete organization on cleanup: %v", err)
		}
	})

	// Wait for org to be created
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		_, err := orgRepo.Get(CTX, pool,
			database.WithCondition(
				orgRepo.PrimaryKeyCondition(Instance.Instance.Id, org.GetOrganizationId()),
			),
		)
		assert.NoError(t, err)
	}, retryDuration, tick)

	// The API call also sets the domain as primary, so we don't do a separate test for that.
	t.Run("test organization domain add reduces", func(t *testing.T) {
		// Add a domain to the organization
		domainName := gofakeit.DomainName()
		beforeAdd := time.Now()
		_, err := OrgClient.AddOrganizationDomain(CTX, &v2.AddOrganizationDomainRequest{
			OrganizationId: org.GetOrganizationId(),
			Domain:         domainName,
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		t.Cleanup(func() {
			_, err := OrgClient.DeleteOrganizationDomain(CTX, &v2.DeleteOrganizationDomainRequest{
				OrganizationId: org.GetOrganizationId(),
				Domain:         domainName,
			})
			if err != nil {
				t.Logf("Failed to delete domain on cleanup: %v", err)
			}
		})

		// Test that domain add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gottenDomain, err := orgDomainRepo.Get(CTX, pool,
				database.WithCondition(
					database.And(
						orgDomainRepo.InstanceIDCondition(Instance.Instance.Id),
						orgDomainRepo.OrgIDCondition(org.OrganizationId),
						orgDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					),
				),
			)
			require.NoError(t, err)
			// event org.domain.added
			assert.Equal(t, domainName, gottenDomain.Domain)
			assert.Equal(t, Instance.Instance.Id, gottenDomain.InstanceID)
			assert.Equal(t, org.OrganizationId, gottenDomain.OrgID)

			assert.WithinRange(t, gottenDomain.CreatedAt, beforeAdd, afterAdd)
			assert.WithinRange(t, gottenDomain.UpdatedAt, beforeAdd, afterAdd)
		}, retryDuration, tick)
	})

	t.Run("test org domain remove reduces", func(t *testing.T) {
		// Add a domain to the organization
		domainName := gofakeit.DomainName()
		_, err := OrgClient.AddOrganizationDomain(CTX, &v2.AddOrganizationDomainRequest{
			OrganizationId: org.GetOrganizationId(),
			Domain:         domainName,
		})
		require.NoError(t, err)

		// Remove the domain
		_, err = OrgClient.DeleteOrganizationDomain(CTX, &v2.DeleteOrganizationDomainRequest{
			OrganizationId: org.GetOrganizationId(),
			Domain:         domainName,
		})
		require.NoError(t, err)

		// Test that domain remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			domain, err := orgDomainRepo.Get(CTX, pool,
				database.WithCondition(
					database.And(
						orgDomainRepo.InstanceIDCondition(Instance.Instance.Id),
						orgDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					),
				),
			)
			// event instance.domain.removed
			assert.Nil(t, domain)
			require.ErrorIs(t, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
