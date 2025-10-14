//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestServer_TestOrgMetadataReduces(t *testing.T) {
	orgMetadataRepo := repository.OrganizationMetadataRepository()

	// The API call also sets the metadata as primary, so we don't do a separate test for that.
	t.Run("test organization metadata add reduces", func(t *testing.T) {
		// Add a metadata to the organization
		org := createTestScopedOrg(t)
		beforeAdd := time.Now()
		_, err := OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-1-bool",
					Value: []byte("false"),
				},
				{
					Key:   "test-2-number",
					Value: []byte("123"),
				},
				{
					Key:   "test-3-object",
					Value: []byte(`{"text":"value", "number":123, "bool": false}`),
				},
				{
					Key:   "test-4-text",
					Value: []byte(`"test-value"`),
				},
				{
					Key:   "test-5-bytes",
					Value: []byte(`test-value`),
				},
			},
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		// Test that metadata add reduces
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		var gottenMetadata []*domain.OrganizationMetadata

		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gottenMetadata, err = orgMetadataRepo.List(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
						orgMetadataRepo.KeyCondition(database.TextOperationStartsWith, "test-"),
					),
				),
				database.WithOrderByAscending(orgMetadataRepo.KeyColumn()),
			)
			require.NoError(t, err)
			require.Len(t, gottenMetadata, 5)
		}, retryDuration, tick)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[0].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[0].OrganizationID)
		assert.Equal(t, "test-1-bool", gottenMetadata[0].Key)
		assert.Equal(t, []byte(`false`), gottenMetadata[0].Value)
		assert.WithinRange(t, gottenMetadata[0].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[0].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[1].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[1].OrganizationID)
		assert.Equal(t, "test-2-number", gottenMetadata[1].Key)
		assert.Equal(t, []byte(`123`), gottenMetadata[1].Value)
		assert.WithinRange(t, gottenMetadata[1].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[1].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[2].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[2].OrganizationID)
		assert.Equal(t, "test-3-object", gottenMetadata[2].Key)
		assert.JSONEq(t, `{"text":"value", "number":123, "bool": false}`, string(gottenMetadata[2].Value))
		assert.WithinRange(t, gottenMetadata[2].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[2].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[3].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[3].OrganizationID)
		assert.Equal(t, "test-4-text", gottenMetadata[3].Key)
		assert.Equal(t, []byte(`"test-value"`), gottenMetadata[3].Value)
		assert.WithinRange(t, gottenMetadata[3].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[3].UpdatedAt, beforeAdd, afterAdd)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[4].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[4].OrganizationID)
		assert.Equal(t, "test-5-bytes", gottenMetadata[4].Key)
		assert.Equal(t, []byte(`test-value`), gottenMetadata[4].Value)
		assert.WithinRange(t, gottenMetadata[4].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[4].UpdatedAt, beforeAdd, afterAdd)
	})

	// The API call also sets the metadata as primary, so we don't do a separate test for that.
	t.Run("ensure update works", func(t *testing.T) {
		// Add a metadata to the organization
		org := createTestScopedOrg(t)
		beforeAdd := time.Now()
		_, err := OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-1-bool",
					Value: []byte("false"),
				},
			},
		})
		require.NoError(t, err)
		afterAdd := time.Now()

		// Test that metadata add reduces
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		var gottenMetadata []*domain.OrganizationMetadata
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gottenMetadata, err = orgMetadataRepo.List(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
						orgMetadataRepo.KeyCondition(database.TextOperationStartsWith, "test-"),
					),
				),
				database.WithOrderByAscending(orgMetadataRepo.KeyColumn()),
			)
			require.NoError(t, err)
			require.Len(t, gottenMetadata, 1)
		}, retryDuration, tick)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[0].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[0].OrganizationID)
		assert.Equal(t, "test-1-bool", gottenMetadata[0].Key)
		assert.Equal(t, []byte(`false`), gottenMetadata[0].Value)
		assert.WithinRange(t, gottenMetadata[0].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[0].UpdatedAt, beforeAdd, afterAdd)

		_, err = OrgClient.SetOrganizationMetadata(CTX, &v2beta.SetOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Metadata: []*v2beta.Metadata{
				{
					Key:   "test-1-bool",
					Value: []byte("true"),
				},
			},
		})
		require.NoError(t, err)
		afterUpdate := time.Now()

		// Test that metadata add reduces
		retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			gottenMetadata, err = orgMetadataRepo.List(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
						orgMetadataRepo.KeyCondition(database.TextOperationStartsWith, "test-"),
					),
				),
				database.WithOrderByAscending(orgMetadataRepo.KeyColumn()),
			)
			require.NoError(t, err)
			require.Len(t, gottenMetadata, 1)
			assert.Greater(t, gottenMetadata[0].UpdatedAt, gottenMetadata[0].CreatedAt)
		}, retryDuration, tick)

		assert.Equal(t, Instance.Instance.Id, gottenMetadata[0].InstanceID)
		assert.Equal(t, org.Id, gottenMetadata[0].OrganizationID)
		assert.Equal(t, "test-1-bool", gottenMetadata[0].Key)
		assert.Equal(t, []byte(`true`), gottenMetadata[0].Value)
		assert.WithinRange(t, gottenMetadata[0].CreatedAt, beforeAdd, afterAdd)
		assert.WithinRange(t, gottenMetadata[0].UpdatedAt, afterAdd, afterUpdate)
	})

	t.Run("test org metadata remove reduces", func(t *testing.T) {
		// Add a metadata to the organization
		org := createTestScopedOrg(t)

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

		// Remove the metadata
		_, err = OrgClient.DeleteOrganizationMetadata(CTX, &v2beta.DeleteOrganizationMetadataRequest{
			OrganizationId: org.GetId(),
			Keys:           []string{"test-bool", "test-text"},
		})
		require.NoError(t, err)

		// Test that metadata remove reduces
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			metadata, err := orgMetadataRepo.Get(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
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
		org := createOrg(t)

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

		// await metadata creation
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			metadata, err := orgMetadataRepo.List(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
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
			metadata, err := orgMetadataRepo.List(CTX, pool,
				database.WithCondition(
					database.And(
						orgMetadataRepo.InstanceIDCondition(Instance.Instance.Id),
						orgMetadataRepo.OrganizationIDCondition(org.Id),
					),
				),
			)
			// event instance.metadata.removed
			require.NoError(t, err)
			assert.Len(t, metadata, 0)
		}, retryDuration, tick)
	})
}
