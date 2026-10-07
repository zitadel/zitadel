package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_stateRejection(t *testing.T) {
	aggregate := &eventstore.Aggregate{ID: userID, Type: user.AggregateType, ResourceOwner: orgID, InstanceID: instanceID}
	rejection := channels.Rejection{
		TriggeringEventType: user.HumanInviteCodeAddedType,
		Reason:              channels.RejectionReasonLimitExceeded,
		ProviderID:          "provider1",
	}
	sendErr := errors.New("send error")
	commandErr := errors.New("command error")

	tests := []struct {
		name    string
		err     error
		expect  func(commands *mock.MockCommands)
		wantErr func(t *testing.T, err error)
	}{
		{
			name: "no error",
			err:  nil,
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "other error is returned unchanged",
			err:  sendErr,
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, sendErr)
				assert.NotErrorIs(t, err, new(channels.CancelError))
			},
		},
		{
			name: "cancel error is returned unchanged",
			err:  channels.NewCancelError(sendErr),
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, new(channels.CancelError))
			},
		},
		{
			name: "rejection is stated and canceled",
			err:  channels.NewRejectedError(channels.RejectionReasonLimitExceeded, "provider1"),
			expect: func(commands *mock.MockCommands) {
				commands.EXPECT().NotificationRejected(gomock.Any(), aggregate, rejection).Return(nil)
			},
			wantErr: func(t *testing.T, err error) {
				// canceled, so it is not retried
				assert.ErrorIs(t, err, new(channels.CancelError))
				assert.ErrorAs(t, err, new(*channels.RejectedError))
			},
		},
		{
			name: "aggregate does not exist anymore, canceled",
			err:  channels.NewRejectedError(channels.RejectionReasonLimitExceeded, "provider1"),
			expect: func(commands *mock.MockCommands) {
				commands.EXPECT().NotificationRejected(gomock.Any(), aggregate, rejection).
					Return(zerrors.ThrowPreconditionFailed(nil, "COMMAND-uXHNj", "Errors.User.NotFound"))
			},
			wantErr: func(t *testing.T, err error) {
				// there is nothing to state the rejection on, but the notification must not be retried
				assert.ErrorIs(t, err, new(channels.CancelError))
			},
		},
		{
			name: "rejection cannot be stated, retried",
			err:  channels.NewRejectedError(channels.RejectionReasonLimitExceeded, "provider1"),
			expect: func(commands *mock.MockCommands) {
				commands.EXPECT().NotificationRejected(gomock.Any(), aggregate, rejection).Return(commandErr)
			},
			wantErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, commandErr)
				assert.NotErrorIs(t, err, new(channels.CancelError))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := mock.NewMockCommands(gomock.NewController(t))
			if tt.expect != nil {
				tt.expect(commands)
			}
			err := stateRejection(t.Context(), commands, aggregate, user.HumanInviteCodeAddedType, tt.err)
			require.NotNil(t, tt.wantErr)
			tt.wantErr(t, err)
		})
	}
}

func Test_stateRejection_session(t *testing.T) {
	aggregate := &eventstore.Aggregate{ID: "session1", Type: session.AggregateType, ResourceOwner: instanceID, InstanceID: instanceID}
	commands := mock.NewMockCommands(gomock.NewController(t))
	commands.EXPECT().NotificationRejected(gomock.Any(), aggregate, channels.Rejection{
		TriggeringEventType: session.OTPEmailChallengedType,
		Reason:              channels.RejectionReasonLimitExceeded,
		ProviderID:          "provider1",
	}).Return(nil)

	err := stateRejection(t.Context(), commands, aggregate, session.OTPEmailChallengedType,
		channels.NewRejectedError(channels.RejectionReasonLimitExceeded, "provider1"))
	assert.ErrorIs(t, err, new(channels.CancelError))
}

func Test_notSent(t *testing.T) {
	event := &user.HumanInviteCodeAddedEvent{
		BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
			InstanceID:    instanceID,
			AggregateID:   userID,
			AggregateType: user.AggregateType,
			ResourceOwner: sql.NullString{String: orgID},
			Typ:           user.HumanInviteCodeAddedType,
		}),
	}
	sendErr := errors.New("send error")

	tests := []struct {
		name     string
		err      error
		expect   func(commands *mock.MockCommands)
		wantDone bool
		wantErr  error
	}{
		{
			name:     "sent",
			err:      nil,
			wantDone: false,
		},
		{
			name:     "send failed, retried",
			err:      sendErr,
			wantDone: true,
			wantErr:  sendErr,
		},
		{
			name:     "canceled, not retried",
			err:      channels.NewCancelError(sendErr),
			wantDone: true,
		},
		{
			name: "rejected, stated and not retried",
			err:  channels.NewRejectedError(channels.RejectionReasonLimitExceeded, "provider1"),
			expect: func(commands *mock.MockCommands) {
				commands.EXPECT().NotificationRejected(gomock.Any(), event.Aggregate(), channels.Rejection{
					TriggeringEventType: user.HumanInviteCodeAddedType,
					Reason:              channels.RejectionReasonLimitExceeded,
					ProviderID:          "provider1",
				}).Return(nil)
			},
			wantDone: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := mock.NewMockCommands(gomock.NewController(t))
			if tt.expect != nil {
				tt.expect(commands)
			}
			done, err := notSent(t.Context(), commands, event, tt.err)
			assert.Equal(t, tt.wantDone, done)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// Test_alreadyHandled_rejected ensures that a rejected notification counts as handled for the aggregates
// that trigger notifications, and that the filter matches the payload of the rejected events.
func Test_alreadyHandled_rejected(t *testing.T) {
	trigger := func(aggregateType eventstore.AggregateType, eventType eventstore.EventType) eventstore.Event {
		return eventstore.BaseEventFromRepo(&repository.Event{
			InstanceID:    instanceID,
			AggregateID:   "aggregate1",
			AggregateType: aggregateType,
			ResourceOwner: sql.NullString{String: orgID},
			Typ:           eventType,
		})
	}
	rejection := channels.Rejection{TriggeringEventType: user.HumanInviteCodeAddedType, Reason: channels.RejectionReasonLimitExceeded}
	userRejected := user.NewNotificationRejectedEvent(t.Context(), &user.NewAggregate("aggregate1", orgID).Aggregate, rejection)
	sessionRejected := session.NewNotificationRejectedEvent(t.Context(), &session.NewAggregate("aggregate1", instanceID).Aggregate, rejection)

	userQuery := (&alreadyHandled{
		event:      trigger(user.AggregateType, user.HumanInviteCodeAddedType),
		eventTypes: []eventstore.EventType{user.HumanInviteCodeSentType},
	}).Query()
	assert.Len(t, userQuery.Matches(userRejected), 1, "rejected user notification is handled")
	assert.Empty(t, userQuery.Matches(sessionRejected), "rejections of other aggregates are ignored")

	sessionQuery := (&alreadyHandled{
		event:      trigger(session.AggregateType, session.OTPEmailChallengedType),
		eventTypes: []eventstore.EventType{session.OTPEmailSentType},
	}).Query()
	assert.Len(t, sessionQuery.Matches(sessionRejected), 1, "rejected session notification is handled")

	// the filter on the payload must match what the rejected events store
	payload, err := json.Marshal(userRejected.Payload())
	require.NoError(t, err)
	stored := make(map[string]any)
	require.NoError(t, json.Unmarshal(payload, &stored))
	for key, value := range channels.RejectionsOf(user.HumanInviteCodeAddedType) {
		assert.EqualValues(t, value, stored[key], key)
	}
}
