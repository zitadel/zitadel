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
	v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestServer_TestOrgMetadataReduces(t *testing.T) {
	org, err := OrgClient.CreateOrganization(CTX, &v2beta.CreateOrganizationRequest{
		Name: gofakeit.Name(),
	})
	require.NoError(t, err)

	orgRepo := repository.OrganizationRepository(pool)
	orgMetadataRepo := orgRepo.Metadata(false)

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
	assert.EventuallyWithT(t, func(t *assert.CollectT) {
		_, err := orgRepo.Get(CTX,
			database.WithCondition(orgRepo.IDCondition(org.GetId())),
		)
		assert.NoError(t, err)
	}, retryDuration, tick)

	// The API call also sets the metadata as primary, so we don't do a separate test for that.
	t.Run("test organization metadata add reduces", func(t *testing.T) {
		// Add a metadata to the organization
		beforeAdd := time.Now()
		_, err := OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-bool",
					Value: []byte("false"),
				},
				{
					Key:   "test-number",
					Value: []byte("123"),
				},
				{
					Key:   "test-object",
					Value: []byte(`{"text":"value", "number":123, "bool": false}`),
				},
				{
					Key:   "test-text",
					Value: []byte(`"test-value"`),
				},
				{
					Key:   "test-bytes",
					Value: []byte(`test-value`),
				},
			},
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		t.Cleanup(func() {
			_, err := OrgClient.DeleteOrganizationMetadata(CTX, &v2beta.DeleteOrganizationMetadataRequest{
				OrganizationId: org.GetId(),
				Keys:           []string{"test-text", "test-number", "test-bool", "test-object"},
			})
			if err != nil {
				t.Logf("Failed to delete metadata on cleanup: %v", err)
			}
		})

		// Test that metadata add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

		var gottenMetadata []*domain.OrganizationMetadata

		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gottenMetadata, err = orgMetadataRepo.List(CTX,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrgIDCondition(org.Id),
						orgMetadataRepo.KeyCondition(database.TextOperationStartsWith, "test-"),
					),
				),
				database.WithOrderByAscending(orgMetadataRepo.KeyColumn()),
			)
			require.NoError(t, err)
			require.Len(t, gottenMetadata, 5)
		}, retryDuration, tick)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[0].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[0].OrgID)
		assert.Equal(t, "test-bool", gottenMetadata[0].Key)
		assert.Equal(t, []byte(`false`), gottenMetadata[0].Value)
		assert.WithinRange(t, gottenMetadata[0].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[0].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[1].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[1].OrgID)
		assert.Equal(t, "test-number", gottenMetadata[1].Key)
		assert.Equal(t, []byte(`123`), gottenMetadata[1].Value)
		assert.WithinRange(t, gottenMetadata[1].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[1].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[2].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[2].OrgID)
		assert.Equal(t, "test-object", gottenMetadata[2].Key)
		assert.JSONEq(t, `{"text":"value", "number":123, "bool": false}`, string(gottenMetadata[2].Value))
		assert.WithinRange(t, gottenMetadata[2].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[2].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[3].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[3].OrgID)
		assert.Equal(t, "test-text", gottenMetadata[3].Key)
		assert.Equal(t, []byte(`"test-value"`), gottenMetadata[3].Value)
		assert.WithinRange(t, gottenMetadata[3].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[3].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[4].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[4].OrgID)
		assert.Equal(t, "test-text", gottenMetadata[4].Key)
		assert.Equal(t, []byte(`test-value`), gottenMetadata[4].Value)
		assert.WithinRange(t, gottenMetadata[4].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[4].UpdatedAt, beforeAdd, afterAdd)

	})

	t.Run("test org metadata remove reduces", func(t *testing.T) {
		// Add a metadata to the organization
		_, err := OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-bool",
					Value: []byte("false"),
				},
				{
					Key:   "test-number",
					Value: []byte("123"),
				},
				{
					Key:   "test-text",
					Value: []byte(`"test-value"`),
				},
			},
		})
		require.NoError(t, err)

		t.Cleanup(func() {
			_, err := OrgClient.DeleteOrganizationMetadata(CTX, &v2beta.DeleteOrganizationMetadataRequest{
				OrganizationId: org.GetId(),
				Keys:           []string{"test-number"},
			})
			if err != nil {
				t.Logf("Failed to delete metadata on cleanup: %v", err)
			}
		})

		// Remove the metadata
		_, err = OrgClient.DeleteOrganizationMetadata(CTX, &v2beta.DeleteOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Keys:           []string{"test-bool", "test-text"},
		})
		require.NoError(t, err)

		// Test that metadata remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			metadata, err := orgMetadataRepo.Get(CTX,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrgIDCondition(org.Id),
					),
				),
			)
			// event instance.metadata.removed
			require.NoError(t, err)
			assert.Equal(t, "test-number", metadata.Key)
		}, retryDuration, tick)
	})

	t.Run("test org metadata removed on org remove", func(t *testing.T) {
		// Add a metadata to the organization
		org, err := OrgClient.CreateOrganization(CTX, &v2beta.CreateOrganizationRequest{
			Name: "some funny name",
		})
		require.NoError(t, err)

		_, err = OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-bool",
					Value: []byte("false"),
				},
				{
					Key:   "test-number",
					Value: []byte("123"),
				},
				{
					Key:   "test-text",
					Value: []byte(`"test-value"`),
				},
			},
		})
		require.NoError(t, err)

		// await metadata creation
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			metadata, err := orgMetadataRepo.List(CTX,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrgIDCondition(org.Id),
					),
				),
			)
			require.NoError(t, err)
			assert.Len(t, metadata, 3)
		}, retryDuration, tick)

		_, err = OrgClient.DeleteOrganization(CTX, &v2beta.DeleteOrganizationRequest{
			Id: org.Id,
		})
		require.NoError(t, err)

		// Test that metadata remove reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			metadata, err := orgMetadataRepo.List(CTX,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrgIDCondition(org.Id),
					),
				),
			)
			// event instance.metadata.removed
			require.NoError(t, err)
			assert.Len(t, metadata, 0)
		}, retryDuration, tick)
	})
}
