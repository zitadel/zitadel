package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/session"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestIDPIntentCheckCommand_Events(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		command        *IDPIntentCheckCommand
		expectedEvents []eventstore.Command
	}{
		{
			name: "nil CheckIntent returns nil events",
			command: &IDPIntentCheckCommand{
				CheckIntent:     nil,
				isCheckComplete: true,
			},
		},
		{
			name: "isCheckComplete false returns nil events",
			command: &IDPIntentCheckCommand{
				CheckIntent:     &session_grpc.CheckIDPIntent{},
				isCheckComplete: false,
			},
		},
		{
			name: "valid command returns events",
			command: &IDPIntentCheckCommand{
				CheckIntent:     &session_grpc.CheckIDPIntent{IdpIntentId: "intent-123"},
				sessionID:       "session-456",
				instanceID:      "instance-789",
				isCheckComplete: true,
			},

			expectedEvents: []eventstore.Command{
				session.NewIntentCheckedEvent(t.Context(), &session.NewAggregate("session-456", "instance-789").Aggregate, time.Now()),
				idpintent.NewConsumedEvent(t.Context(), &idpintent.NewAggregate("intent-123", "").Aggregate),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Test
			events, err := tc.command.Events(context.Background(), &InvokeOpts{})

			// Verify
			assert.NoError(t, err)

			require.Len(t, events, len(tc.expectedEvents))

			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *session.IntentCheckedEvent:
					actualAssertedType, ok := events[i].(*session.IntentCheckedEvent)
					require.True(t, ok)
					assert.InDelta(t, expectedAssertedType.CheckedAt.UnixMicro(), actualAssertedType.CheckedAt.UnixMicro(), 1.5)
				case *idpintent.ConsumedEvent:
					_, ok := events[i].(*idpintent.ConsumedEvent)
					require.True(t, ok)
				}
			}
		})
	}
}
