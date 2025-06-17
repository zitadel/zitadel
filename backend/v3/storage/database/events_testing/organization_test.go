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
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestServer_TestOrganizationReduces(t *testing.T) {
	t.Run("test org add reduces", func(t *testing.T) {
		beforeCreate := time.Now()
		orgName := gofakeit.Name()

		_, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)
		afterCreate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			// event org.added
			require.NotNil(t, organization.ID)
			require.Equal(t, orgName, organization.Name)
			require.NotNil(t, organization.InstanceID)
			require.Equal(t, domain.Active, organization.State)
			assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
			require.Nil(t, organization.DeletedAt)
		}, retryDuration, tick)
	})

	t.Run("test org change reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		_, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)

		// 2. update org name
		beforeUpdate := time.Now()
		orgName = orgName + "_new"
		_, err = MgmtClient.UpdateOrg(CTX, &management.UpdateOrgRequest{
			Name: orgName,
		})
		require.NoError(t, err)
		afterUpdate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			// event org.changed
			require.Equal(t, orgName, organization.Name)
			assert.WithinRange(t, organization.UpdatedAt, beforeUpdate, afterUpdate)
		}, retryDuration, tick)
	})

	t.Run("test org deactivate reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)

		// 2. deactivate org name
		beforeDeactivate := time.Now()
		_ = Instance.DeactivateOrganization(CTX, organization.OrganizationId)
		require.NoError(t, err)
		afterDeactivate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			// event org.deactivate
			require.Equal(t, orgName, organization.Name)
			require.Equal(t, domain.Inactive, organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeDeactivate, afterDeactivate)
		}, retryDuration, tick)
	})

	t.Run("test org activate reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)

		// 2. deactivate org name
		_ = Instance.DeactivateOrganization(CTX, organization.OrganizationId)
		require.NoError(t, err)

		orgRepo := repository.OrganizationRepository(pool)
		// 3. check org deactivated
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			require.Equal(t, orgName, organization.Name)
			require.Equal(t, domain.Inactive, organization.State)
		}, retryDuration, tick)

		// 4. activate org name
		beforeActivate := time.Now()
		_ = Instance.ReactivateOrganization(CTX, organization.OrganizationId)
		require.NoError(t, err)
		afterActivate := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			// event org.reactivate
			require.Equal(t, orgName, organization.Name)
			require.Equal(t, domain.Active, organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeActivate, afterActivate)
		}, retryDuration, tick)
	})

	t.Run("test org remove reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
			Name: orgName,
		})
		require.NoError(t, err)

		// 2. check org retrivable
		orgRepo := repository.OrganizationRepository(pool)
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			if organization == nil {
				require.Fail(t, "this error is here because of a race condition")
			}
			require.Equal(t, orgName, organization.Name)
		}, retryDuration, tick)

		// 3. delete org
		_ = Instance.RemoveOrganization(CTX, organization.OrganizationId)
		require.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			require.NoError(t, err)

			// event org.remove
			require.Nil(t, organization)
		}, retryDuration, tick)
	})
}
