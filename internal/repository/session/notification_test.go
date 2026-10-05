package session

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

func TestNotificationRejectedEvent(t *testing.T) {
	agg := &eventstore.Aggregate{ID: "session1", Type: AggregateType, ResourceOwner: "instance1", InstanceID: "instance1", Version: AggregateVersion}
	rejection := channels.Rejection{
		TriggeringEventType: OTPEmailChallengedType,
		Reason:              channels.RejectionReasonLimitExceeded,
		ProviderID:          "provider1",
	}
	event := NewNotificationRejectedEvent(t.Context(), agg, rejection)
	assert.Equal(t, eventstore.EventType(NotificationRejectedType), event.Type())

	payload, err := json.Marshal(event.Payload())
	require.NoError(t, err)
	assert.JSONEq(t, `{"triggeringEventType":"session.otp.email.challenged","reason":"limit_exceeded","providerId":"provider1"}`, string(payload))

	mapped, err := eventstore.GenericEventMapper[NotificationRejectedEvent](&eventstore.BaseEvent{Agg: agg, Data: payload})
	require.NoError(t, err)
	assert.Equal(t, rejection, mapped.(*NotificationRejectedEvent).Rejection)
}
