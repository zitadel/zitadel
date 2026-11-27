package samlidpmetadata

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

type queriesFunc func(context.Context, string, string, uint32) ([]query.SAMLIDPMetadataURL, error)

func (f queriesFunc) ListSAMLIDPMetadataURLs(ctx context.Context, afterInstanceID, afterID string, limit uint32) ([]query.SAMLIDPMetadataURL, error) {
	return f(ctx, afterInstanceID, afterID, limit)
}

type commandCall struct {
	ownerType     domain.IdentityProviderType
	resourceOwner string
	id            string
}

type commandsRecorder struct {
	calls []commandCall
	err   error
}

func (c *commandsRecorder) RefreshInstanceSAMLProviderMetadata(_ context.Context, id string) error {
	c.calls = append(c.calls, commandCall{ownerType: domain.IdentityProviderTypeSystem, id: id})
	return c.err
}

func (c *commandsRecorder) RefreshOrgSAMLProviderMetadata(_ context.Context, resourceOwner, id string) error {
	c.calls = append(c.calls, commandCall{ownerType: domain.IdentityProviderTypeOrg, resourceOwner: resourceOwner, id: id})
	return c.err
}

func TestWorkerWork(t *testing.T) {
	t.Run("query error", func(t *testing.T) {
		queryErr := errors.New("query failed")
		worker := NewWorker(queriesFunc(func(context.Context, string, string, uint32) ([]query.SAMLIDPMetadataURL, error) {
			return nil, queryErr
		}), new(commandsRecorder))

		err := worker.Work(context.Background(), new(river.Job[*SAMLIDPMetadataRefresh]))
		require.ErrorIs(t, err, queryErr)
	})

	t.Run("processes every page and continues after provider error", func(t *testing.T) {
		firstPage := make([]query.SAMLIDPMetadataURL, batchSize)
		for i := range firstPage {
			firstPage[i] = query.SAMLIDPMetadataURL{
				ID:         fmt.Sprintf("idp-%03d", i),
				InstanceID: "instance-1",
				OwnerType:  domain.IdentityProviderTypeSystem,
			}
		}
		secondPage := []query.SAMLIDPMetadataURL{{
			ID:            "idp-100",
			InstanceID:    "instance-2",
			ResourceOwner: "org-1",
			OwnerType:     domain.IdentityProviderTypeOrg,
		}}
		queries := queriesFunc(func(_ context.Context, afterInstanceID, afterID string, limit uint32) ([]query.SAMLIDPMetadataURL, error) {
			assert.Equal(t, uint32(batchSize), limit)
			switch {
			case afterInstanceID == "" && afterID == "":
				return firstPage, nil
			case afterInstanceID == "instance-1" && afterID == "idp-099":
				return secondPage, nil
			default:
				t.Fatalf("unexpected cursor (%q, %q)", afterInstanceID, afterID)
				return nil, nil
			}
		})
		commands := &commandsRecorder{err: errors.New("refresh failed")}
		worker := NewWorker(queries, commands)

		require.NoError(t, worker.Work(context.Background(), new(river.Job[*SAMLIDPMetadataRefresh])))
		require.Len(t, commands.calls, batchSize+1)
		assert.Equal(t, commandCall{ownerType: domain.IdentityProviderTypeOrg, resourceOwner: "org-1", id: "idp-100"}, commands.calls[batchSize])
	})
}

func TestWorkerRefresh(t *testing.T) {
	commands := new(commandsRecorder)
	worker := NewWorker(nil, commands)

	require.NoError(t, worker.refresh(context.Background(), query.SAMLIDPMetadataURL{
		ID:         "instance-idp",
		InstanceID: "instance-1",
		OwnerType:  domain.IdentityProviderTypeSystem,
	}))
	require.NoError(t, worker.refresh(context.Background(), query.SAMLIDPMetadataURL{
		ID:            "org-idp",
		InstanceID:    "instance-1",
		ResourceOwner: "org-1",
		OwnerType:     domain.IdentityProviderTypeOrg,
	}))
	err := worker.refresh(context.Background(), query.SAMLIDPMetadataURL{
		ID:         "invalid-idp",
		InstanceID: "instance-1",
		OwnerType:  domain.IdentityProviderType(99),
	})
	require.ErrorContains(t, err, "unsupported identity provider owner type")
	assert.Equal(t, []commandCall{
		{ownerType: domain.IdentityProviderTypeSystem, id: "instance-idp"},
		{ownerType: domain.IdentityProviderTypeOrg, resourceOwner: "org-1", id: "org-idp"},
	}, commands.calls)
}
