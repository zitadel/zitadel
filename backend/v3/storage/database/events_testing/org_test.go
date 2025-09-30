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
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func createOrg(t *testing.T) *org.CreateOrganizationResponse {
	t.Helper()
	org, err := OrgClient.CreateOrganization(CTX, &org.CreateOrganizationRequest{
		Name: gofakeit.Name(),
	})
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(t.Context(), time.Minute)
	assert.EventuallyWithT(t, func(tc *assert.CollectT) {
		_, err := orgRepo.Get(t.Context(), pool,
			database.WithCondition(database.And(
				orgRepo.InstanceIDCondition(Instance.Instance.Id),
				orgRepo.IDCondition(org.GetId())),
			),
		)
		assert.NoError(tc, err)
	}, retryDuration, tick)

	return org
}

func createTestScopedOrg(t *testing.T) *org.CreateOrganizationResponse {
	org := createOrg(t)

	t.Cleanup(func() {
		_, err := OrgClient.DeleteOrganization(CTX, &v2beta.DeleteOrganizationRequest{
			Id: org.GetId(),
		})
		if err != nil {
			t.Logf("Failed to delete organization on cleanup: %v", err)
		}
	})

	return org
}
