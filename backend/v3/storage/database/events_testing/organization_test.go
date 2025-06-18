//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestServer_TestOrganizationReduces(t *testing.T) {
	t.Run("test org add reduces", func(t *testing.T) {
		beforeCreate := time.Now()
		orgName := gofakeit.Name()

		_, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		assert.NoError(t, err)
		afterCreate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(tt *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(tt, err)

			// event org.added
			assert.NotNil(t, organization.ID)
			assert.Equal(t, orgName, organization.Name)
			assert.NotNil(t, organization.InstanceID)
			assert.Equal(t, domain.OrgStateActive.String(), organization.State)
			assert.WithinRange(t, organization.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, organization.UpdatedAt, beforeCreate, afterCreate)
			assert.Nil(t, organization.DeletedAt)
		}, retryDuration, tick)
	})

	t.Run("test org change reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		assert.NoError(t, err)

		// 2. update org name
		beforeUpdate := time.Now()
		orgName = orgName + "_new"
		_, err = OrgClient.UpdateOrganization(CTX, &v2beta_org.UpdateOrganizationRequest{
			Id:   organization.Id,
			Name: orgName,
		})
		assert.NoError(t, err)
		afterUpdate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			// event org.changed
			assert.Equal(t, orgName, organization.Name)
			assert.WithinRange(t, organization.UpdatedAt, beforeUpdate, afterUpdate)
		}, retryDuration, tick)
	})

	t.Run("test org deactivate reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		assert.NoError(t, err)

		// 2. deactivate org name
		beforeDeactivate := time.Now()
		_, err = OrgClient.DeactivateOrganization(CTX, &v2beta_org.DeactivateOrganizationRequest{
			Id: organization.Id,
		})

		assert.NoError(t, err)
		afterDeactivate := time.Now()

		orgRepo := repository.OrganizationRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			// event org.deactivate
			assert.Equal(t, orgName, organization.Name)
			assert.Equal(t, domain.OrgStateInactive.String(), organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeDeactivate, afterDeactivate)
		}, retryDuration, tick)
	})

	t.Run("test org activate reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		assert.NoError(t, err)

		// 2. deactivate org name
		_, err = OrgClient.DeactivateOrganization(CTX, &v2beta_org.DeactivateOrganizationRequest{
			Id: organization.Id,
		})
		assert.NoError(t, err)

		orgRepo := repository.OrganizationRepository(pool)
		// 3. check org deactivated
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			assert.Equal(t, orgName, organization.Name)
			assert.Equal(t, domain.OrgStateInactive.String(), organization.State)
		}, retryDuration, tick)

		// 4. activate org name
		beforeActivate := time.Now()
		_, err = OrgClient.ActivateOrganization(CTX, &v2beta_org.ActivateOrganizationRequest{
			Id: organization.Id,
		})
		assert.NoError(t, err)
		afterActivate := time.Now()

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			// event org.reactivate
			assert.Equal(t, orgName, organization.Name)
			assert.Equal(t, domain.OrgStateActive.String(), organization.State)
			assert.WithinRange(t, organization.UpdatedAt, beforeActivate, afterActivate)
		}, retryDuration, tick)
	})

	t.Run("test org remove reduces", func(t *testing.T) {
		orgName := gofakeit.Name()

		// 1. create org
		organization, err := OrgClient.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
			Name: orgName,
		})
		assert.NoError(t, err)

		// 2. check org retrivable
		orgRepo := repository.OrganizationRepository(pool)
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			if organization == nil {
				assert.Fail(t, "this error is here because of a race condition")
			}
			assert.Equal(t, orgName, organization.Name)
		}, retryDuration, tick)

		// 3. delete org
		_, err = OrgClient.DeleteOrganization(CTX, &v2beta_org.DeleteOrganizationRequest{
			Id: organization.Id,
		})
		assert.NoError(t, err)

		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			organization, err := orgRepo.Get(CTX,
				orgRepo.NameCondition(database.TextOperationEqual, orgName),
			)
			assert.NoError(t, err)

			// event org.remove
			assert.Nil(t, organization)
		}, retryDuration, tick)
	})
}
