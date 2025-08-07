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
	v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestServer_TestOrgDomainReduces(t *testing.T) {
	org, err := OrgClient.CreateOrganization(CTX, &v2beta.CreateOrganizationRequest{
		Name: gofakeit.Name(),
	})
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(pool)
	orgDomainRepo := orgRepo.Domains(false)

	t.Cleanup(func() {
		_, err := OrgClient.DeleteOrganization(CTX, &v2beta.DeleteOrganizationRequest{
			Id: org.GetId(),
		})
		if err != nil {
			t.Logf("Failed to delete organization on cleanup: %v", err)
		}
	})

	// Wait for org to be created
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := orgRepo.Get(CTX,
			database.WithCondition(orgRepo.IDCondition(org.GetId())),
		)
		assert.NoError(ttt, err)
	}, retryDuration, tick)

	// The API call also sets the domain as primary, so we don't do a separate test for that.
	t.Run("test organization domain add reduces", func(t *testing.T) {
		// Add a domain to the organization
		domainName := gofakeit.DomainName()
		beforeAdd := time.Now()
		_, err := OrgClient.AddOrganizationDomain(CTX, &v2beta.AddOrganizationDomainRequest{
			OrganizationId: org.GetId(),
			Domain:         domainName,
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		t.Cleanup(func() {
			_, err := OrgClient.DeleteOrganizationDomain(CTX, &v2beta.DeleteOrganizationDomainRequest{
				OrganizationId: org.GetId(),
				Domain:         domainName,
			})
			if err != nil {
				t.Logf("Failed to delete domain on cleanup: %v", err)
			}
		})

		// Test that domain add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			gottenDomain, err := orgDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						orgDomainRepo.InstanceIDCondition(Instance.Instance.Id),
						orgDomainRepo.OrgIDCondition(org.Id),
						orgDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					),
				),
			)
			require.NoError(ttt, err)
			// event org.domain.added
			assert.Equal(t, domainName, gottenDomain.Domain)
			assert.Equal(t, Instance.Instance.Id, gottenDomain.InstanceID)
			assert.Equal(t, org.Id, gottenDomain.OrgID)

			assert.WithinRange(t, gottenDomain.CreatedAt, beforeAdd, afterAdd)
			assert.WithinRange(t, gottenDomain.UpdatedAt, beforeAdd, afterAdd)
		}, retryDuration, tick)
	})

	t.Run("test org domain remove reduces", func(t *testing.T) {
		// Add a domain to the organization
		domainName := gofakeit.DomainName()
		_, err := OrgClient.AddOrganizationDomain(CTX, &v2beta.AddOrganizationDomainRequest{
			OrganizationId: org.GetId(),
			Domain:         domainName,
		})
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := OrgClient.DeleteOrganizationDomain(CTX, &v2beta.DeleteOrganizationDomainRequest{
				OrganizationId: org.GetId(),
				Domain:         domainName,
			})
			if err != nil {
				t.Logf("Failed to delete domain on cleanup: %v", err)
			}
		})

		// Remove the domain
		_, err = OrgClient.DeleteOrganizationDomain(CTX, &v2beta.DeleteOrganizationDomainRequest{
			OrganizationId: org.GetId(),
			Domain:         domainName,
		})
		require.NoError(t, err)

		// Test that domain remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
			domain, err := orgDomainRepo.Get(CTX,
				database.WithCondition(
					database.And(
						orgDomainRepo.InstanceIDCondition(Instance.Instance.Id),
						orgDomainRepo.DomainCondition(database.TextOperationEqual, domainName),
					),
				),
			)
			// event instance.domain.removed
			assert.Nil(ttt, domain)
			require.ErrorIs(ttt, err, new(database.NoRowFoundError))
		}, retryDuration, tick)
	})
}
