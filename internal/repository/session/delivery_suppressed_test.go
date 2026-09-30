package session

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestOTPEmailSentEvent_DeliverySuppressed(t *testing.T) {
	agg := &eventstore.Aggregate{ID: "session1", Type: AggregateType, ResourceOwner: "instance1", InstanceID: "instance1", Version: AggregateVersion}
	for _, suppressed := range []bool{false, true} {
		event := NewOTPEmailSentEvent(t.Context(), agg, suppressed)
		payload, err := json.Marshal(event.Payload())
		require.NoError(t, err)
		mapped, err := eventstore.GenericEventMapper[OTPEmailSentEvent](&eventstore.BaseEvent{Agg: agg, Data: payload})
		require.NoError(t, err)
		assert.Equal(t, suppressed, mapped.(*OTPEmailSentEvent).DeliverySuppressed)
	}
	mapped, err := eventstore.GenericEventMapper[OTPEmailSentEvent](&eventstore.BaseEvent{Agg: agg})
	require.NoError(t, err)
	assert.False(t, mapped.(*OTPEmailSentEvent).DeliverySuppressed)
}
