package instance

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// TestClientIDMetadataDocumentAllowedURLField pins that the unique field survives the
// eventstore lowercasing unique fields: client_id URLs are compared case-sensitively, so two
// URLs that differ only in case must get different fields.
func TestClientIDMetadataDocumentAllowedURLField(t *testing.T) {
	lower := clientIDMetadataDocumentAllowedURLField("https://app.example.com/client")
	upper := clientIDMetadataDocumentAllowedURLField("https://app.example.com/Client")
	assert.NotEqual(t, strings.ToLower(lower), strings.ToLower(upper))
	assert.Equal(t, lower, strings.ToLower(lower), "the field must already be lowercase")
	assert.Len(t, clientIDMetadataDocumentAllowedURLField("https://app.example.com/"+strings.Repeat("a", 2000)), 64)
}

func TestSecurityPolicyClientIDMetadataDocumentAllowedURLEvents_UniqueConstraints(t *testing.T) {
	agg := &NewAggregate("instance").Aggregate
	const allowedURL = "https://app.example.com/client"

	added := NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(), agg, allowedURL).UniqueConstraints()
	require.Len(t, added, 1)
	assert.Equal(t, eventstore.UniqueConstraintAdd, added[0].Action)
	assert.Equal(t, UniqueClientIDMetadataDocumentAllowedURL, added[0].UniqueType)
	assert.Equal(t, clientIDMetadataDocumentAllowedURLField(allowedURL), added[0].UniqueField)
	assert.Equal(t, []string{eventstore.OwnerTag(clientIDMetadataDocumentAllowedURLOwner, "instance")}, added[0].Owners,
		"an added url must be tagged so that replacing the list releases it")

	removed := NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(), agg, allowedURL).UniqueConstraints()
	require.Len(t, removed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, removed[0].Action)
	assert.Equal(t, added[0].UniqueField, removed[0].UniqueField)
}

// TestSecurityPolicySetEvent_UniqueConstraints pins that replacing the list releases every
// allowed URL constraint of the instance and takes one per URL in the new list. The constraints
// depend only on the stored list, so a push retried after a concurrent change still leaves
// exactly the constraints of the list it stores.
func TestSecurityPolicySetEvent_UniqueConstraints(t *testing.T) {
	event, err := NewSecurityPolicySetEvent(context.Background(), &NewAggregate("instance").Aggregate, []SecurityPolicyChanges{
		ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://a.example.com/", "https://b.example.com/"}),
	})
	require.NoError(t, err)

	constraints := event.UniqueConstraints()
	require.Len(t, constraints, 3)
	assert.Equal(t, eventstore.UniqueConstraintRemoveByOwner, constraints[0].Action)
	assert.Equal(t, []string{eventstore.OwnerTag(clientIDMetadataDocumentAllowedURLOwner, "instance")}, constraints[0].Owners)
	assert.Equal(t, eventstore.UniqueConstraintAdd, constraints[1].Action)
	assert.Equal(t, clientIDMetadataDocumentAllowedURLField("https://a.example.com/"), constraints[1].UniqueField)
	assert.Equal(t, clientIDMetadataDocumentAllowedURLField("https://b.example.com/"), constraints[2].UniqueField)

	cleared, err := NewSecurityPolicySetEvent(context.Background(), &NewAggregate("instance").Aggregate, []SecurityPolicyChanges{
		ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs(nil),
	})
	require.NoError(t, err)
	require.Len(t, cleared.UniqueConstraints(), 1, "clearing the list releases every constraint")

	unchanged, err := NewSecurityPolicySetEvent(context.Background(), &NewAggregate("instance").Aggregate, []SecurityPolicyChanges{
		ChangeSecurityPolicyEnableClientIDMetadataDocument(true),
	})
	require.NoError(t, err)
	assert.Empty(t, unchanged.UniqueConstraints(), "a set that does not replace the list touches no constraint")
}
