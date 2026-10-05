package user

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
)

// TestDeliveryInfo ensures that the delivery information of the email sent events survives the round trip
// and that events stored before it existed are still mapped.
func TestDeliveryInfo(t *testing.T) {
	ctx := t.Context()
	agg := &eventstore.Aggregate{ID: "user1", Type: AggregateType, ResourceOwner: "org1", InstanceID: "instance1", Version: AggregateVersion}

	tests := []struct {
		name   string
		event  func(info senders.DeliveryInfo) eventstore.Command
		mapper func(eventstore.Event) (eventstore.Event, error)
		info   func(eventstore.Event) senders.DeliveryInfo
		// payloadLessBefore marks events which had no payload before the delivery information was introduced
		payloadLessBefore bool
	}{
		{
			name:              "initial code",
			payloadLessBefore: true,
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanInitialCodeSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[HumanInitialCodeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanInitialCodeSentEvent).DeliveryInfo },
		},
		{
			name:              "email code",
			payloadLessBefore: true,
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanEmailCodeSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[HumanEmailCodeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanEmailCodeSentEvent).DeliveryInfo },
		},
		{
			name:              "invite code",
			payloadLessBefore: true,
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanInviteCodeSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[HumanInviteCodeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanInviteCodeSentEvent).DeliveryInfo },
		},
		{
			name: "password code",
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanPasswordCodeSentEvent(ctx, agg, nil, info)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordCodeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanPasswordCodeSentEvent).DeliveryInfo },
		},
		{
			name: "otp email code",
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanOTPEmailCodeSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[HumanOTPEmailCodeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanOTPEmailCodeSentEvent).DeliveryInfo },
		},
		{
			name: "passwordless init code",
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanPasswordlessInitCodeSentEvent(ctx, agg, "code1", info)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordlessInitCodeSentEvent],
			info: func(e eventstore.Event) senders.DeliveryInfo {
				return e.(*HumanPasswordlessInitCodeSentEvent).DeliveryInfo
			},
		},
		{
			name:              "domain claimed",
			payloadLessBefore: true,
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewDomainClaimedSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[DomainClaimedSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*DomainClaimedSentEvent).DeliveryInfo },
		},
		{
			name:              "password change",
			payloadLessBefore: true,
			event: func(info senders.DeliveryInfo) eventstore.Command {
				return NewHumanPasswordChangeSentEvent(ctx, agg, info)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordChangeSentEvent],
			info:   func(e eventstore.Event) senders.DeliveryInfo { return e.(*HumanPasswordChangeSentEvent).DeliveryInfo },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, info := range []senders.DeliveryInfo{
				{},
				{ProviderID: "provider1"},
				{DeliverySuppressed: true},
			} {
				payload, err := json.Marshal(tt.event(info).Payload())
				require.NoError(t, err)
				switch {
				case info.DeliverySuppressed:
					assert.Contains(t, string(payload), `"deliverySuppressed":true`)
					assert.NotContains(t, string(payload), "providerId")
				case info.ProviderID != "":
					assert.Contains(t, string(payload), `"providerId":"provider1"`)
					assert.NotContains(t, string(payload), "deliverySuppressed")
				default:
					assert.NotContains(t, string(payload), "deliverySuppressed")
					assert.NotContains(t, string(payload), "providerId")
					if tt.payloadLessBefore {
						// the stored payload must not change for events which had none before
						assert.Nil(t, tt.event(info).Payload())
					}
				}
				mapped, err := tt.mapper(&eventstore.BaseEvent{Agg: agg, Data: payload})
				require.NoError(t, err)
				assert.Equal(t, info, tt.info(mapped))
			}

			// events stored before the delivery information existed carry no payload at all
			mapped, err := tt.mapper(&eventstore.BaseEvent{Agg: agg})
			require.NoError(t, err)
			assert.Zero(t, tt.info(mapped))
		})
	}
}
