package user

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// TestDeliverySuppressed ensures that the flag of the email sent events survives the round trip
// and that events stored before the flag existed are still mapped.
func TestDeliverySuppressed(t *testing.T) {
	ctx := t.Context()
	agg := &eventstore.Aggregate{ID: "user1", Type: AggregateType, ResourceOwner: "org1", InstanceID: "instance1", Version: AggregateVersion}

	tests := []struct {
		name   string
		event  func(suppressed bool) eventstore.Command
		mapper func(eventstore.Event) (eventstore.Event, error)
		flag   func(eventstore.Event) bool
	}{
		{
			name: "initial code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanInitialCodeSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanInitialCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanInitialCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "email code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanEmailCodeSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanEmailCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanEmailCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "invite code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanInviteCodeSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanInviteCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanInviteCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "password code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanPasswordCodeSentEvent(ctx, agg, nil, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanPasswordCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "otp email code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanOTPEmailCodeSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanOTPEmailCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanOTPEmailCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "passwordless init code",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanPasswordlessInitCodeSentEvent(ctx, agg, "code1", suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordlessInitCodeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanPasswordlessInitCodeSentEvent).DeliverySuppressed },
		},
		{
			name: "domain claimed",
			event: func(suppressed bool) eventstore.Command {
				return NewDomainClaimedSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[DomainClaimedSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*DomainClaimedSentEvent).DeliverySuppressed },
		},
		{
			name: "password change",
			event: func(suppressed bool) eventstore.Command {
				return NewHumanPasswordChangeSentEvent(ctx, agg, suppressed)
			},
			mapper: eventstore.GenericEventMapper[HumanPasswordChangeSentEvent],
			flag:   func(e eventstore.Event) bool { return e.(*HumanPasswordChangeSentEvent).DeliverySuppressed },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, suppressed := range []bool{false, true} {
				payload, err := json.Marshal(tt.event(suppressed).Payload())
				require.NoError(t, err)
				if suppressed {
					assert.Contains(t, string(payload), `"deliverySuppressed":true`)
				} else {
					assert.NotContains(t, string(payload), "deliverySuppressed")
				}
				mapped, err := tt.mapper(&eventstore.BaseEvent{Agg: agg, Data: payload})
				require.NoError(t, err)
				assert.Equal(t, suppressed, tt.flag(mapped))
			}

			// events stored before the flag existed carry no payload at all
			mapped, err := tt.mapper(&eventstore.BaseEvent{Agg: agg})
			require.NoError(t, err)
			assert.False(t, tt.flag(mapped))
		})
	}
}
