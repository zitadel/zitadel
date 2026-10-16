package session

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
)

func TestOTPEmailSentEvent_DeliveryInfo(t *testing.T) {
	agg := &eventstore.Aggregate{ID: "session1", Type: AggregateType, ResourceOwner: "instance1", InstanceID: "instance1", Version: AggregateVersion}
	for _, info := range []senders.DeliveryInfo{
		{},
		{ProviderID: "provider1"},
		{DeliverySuppressed: true},
	} {
		event := NewOTPEmailSentEvent(t.Context(), agg, info)
		payload, err := json.Marshal(event.Payload())
		require.NoError(t, err)
		mapped, err := eventstore.GenericEventMapper[OTPEmailSentEvent](&eventstore.BaseEvent{Agg: agg, Data: payload})
		require.NoError(t, err)
		assert.Equal(t, info, mapped.(*OTPEmailSentEvent).DeliveryInfo)
	}
	// events stored before the delivery information existed carry no payload at all
	mapped, err := eventstore.GenericEventMapper[OTPEmailSentEvent](&eventstore.BaseEvent{Agg: agg})
	require.NoError(t, err)
	assert.Zero(t, mapped.(*OTPEmailSentEvent).DeliveryInfo)
}
