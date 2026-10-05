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

	removed := NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(), agg, allowedURL).UniqueConstraints()
	require.Len(t, removed, 1)
	assert.Equal(t, eventstore.UniqueConstraintRemove, removed[0].Action)
	assert.Equal(t, added[0].UniqueField, removed[0].UniqueField)
}

// TestChangeSecurityPolicyClientIDMetadataDocumentAllowedURLConstraints pins that replacing the
// list takes the constraint of every new URL and releases the one of every dropped URL, and
// leaves URLs that stay untouched.
func TestChangeSecurityPolicyClientIDMetadataDocumentAllowedURLConstraints(t *testing.T) {
	event, err := NewSecurityPolicySetEvent(context.Background(), &NewAggregate("instance").Aggregate, []SecurityPolicyChanges{
		ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://kept.example.com/", "https://new.example.com/"}),
		ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLConstraints(
			[]string{"https://kept.example.com/", "https://dropped.example.com/"},
			[]string{"https://kept.example.com/", "https://new.example.com/"},
		),
	})
	require.NoError(t, err)

	constraints := event.UniqueConstraints()
	require.Len(t, constraints, 2)
	assert.Equal(t, eventstore.UniqueConstraintRemove, constraints[0].Action)
	assert.Equal(t, clientIDMetadataDocumentAllowedURLField("https://dropped.example.com/"), constraints[0].UniqueField)
	assert.Equal(t, eventstore.UniqueConstraintAdd, constraints[1].Action)
	assert.Equal(t, clientIDMetadataDocumentAllowedURLField("https://new.example.com/"), constraints[1].UniqueField)

	unchanged, err := NewSecurityPolicySetEvent(context.Background(), &NewAggregate("instance").Aggregate, []SecurityPolicyChanges{
		ChangeSecurityPolicyEnableClientIDMetadataDocument(true),
	})
	require.NoError(t, err)
	assert.Empty(t, unchanged.UniqueConstraints(), "a set that does not replace the list touches no constraint")
}
