package schemauser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestEmailCodeSentEvent_DeliverySuppressed(t *testing.T) {
	agg := &eventstore.Aggregate{ID: "user1", Type: AggregateType, ResourceOwner: "org1", InstanceID: "instance1", Version: AggregateVersion}
	for _, suppressed := range []bool{false, true} {
		event := NewEmailCodeSentEvent(t.Context(), agg, suppressed)
		payload, err := json.Marshal(event.Payload())
		require.NoError(t, err)
		mapped, err := eventstore.GenericEventMapper[EmailCodeSentEvent](&eventstore.BaseEvent{Agg: agg, Data: payload})
		require.NoError(t, err)
		assert.Equal(t, suppressed, mapped.(*EmailCodeSentEvent).DeliverySuppressed)
	}
	mapped, err := eventstore.GenericEventMapper[EmailCodeSentEvent](&eventstore.BaseEvent{Agg: agg})
	require.NoError(t, err)
	assert.False(t, mapped.(*EmailCodeSentEvent).DeliverySuppressed)
}
