package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestHumanRecoveryCodeWriteModel_Reduce(t *testing.T) {
	userAgg := user.NewAggregate("user1", "org1")
	ctx := context.Background()

	tests := []struct {
		name   string
		events []eventstore.Event
		want   *HumanRecoveryCodeWriteModel
	}{
		{
			name:   "no events",
			events: []eventstore.Event{},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "recovery codes added",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "recovery codes added multiple times",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2"},
					nil,
				),
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code3", "code4"},
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code2", "code3", "code4"},
				userLocked:     false,
			},
		},
		{
			name: "recovery code check succeeded",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckSucceededEvent(
					ctx,
					&userAgg.Aggregate,
					"code2",
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "recovery code check succeeded with code that no longer exists",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckSucceededEvent(
					ctx,
					&userAgg.Aggregate,
					"nonexistent",
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "recovery code check failed",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 1,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "multiple recovery code check failures",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 3,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "recovery codes removed",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeRemovedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateRemoved,
				FailedAttempts: 0,
				codes:          []string{},
				userLocked:     false,
			},
		},
		{
			name: "user locked",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewUserLockedEvent(
					ctx,
					&userAgg.Aggregate,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 1,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     true,
			},
		},
		{
			name: "user unlocked",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewUserLockedEvent(
					ctx,
					&userAgg.Aggregate,
				),
				user.NewUserUnlockedEvent(
					ctx,
					&userAgg.Aggregate,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "user removed",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3"},
					nil,
				),
				user.NewUserRemovedEvent(
					ctx,
					&userAgg.Aggregate,
					"username",
					nil,
					false,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateRemoved,
				FailedAttempts: 0,
				codes:          []string{"code1", "code2", "code3"},
				userLocked:     false,
			},
		},
		{
			name: "complex flow with multiple events",
			events: []eventstore.Event{
				user.NewHumanRecoveryCodesAddedEvent(
					ctx,
					&userAgg.Aggregate,
					[]string{"code1", "code2", "code3", "code4"},
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewHumanRecoveryCodeCheckSucceededEvent(
					ctx,
					&userAgg.Aggregate,
					"code2",
					nil,
				),
				user.NewHumanRecoveryCodeCheckFailedEvent(
					ctx,
					&userAgg.Aggregate,
					nil,
				),
				user.NewUserLockedEvent(
					ctx,
					&userAgg.Aggregate,
				),
				user.NewUserUnlockedEvent(
					ctx,
					&userAgg.Aggregate,
				),
				user.NewHumanRecoveryCodeCheckSucceededEvent(
					ctx,
					&userAgg.Aggregate,
					"code4",
					nil,
				),
			},
			want: &HumanRecoveryCodeWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "user1",
					ResourceOwner: "org1",
				},
				State:          domain.MFAStateReady,
				FailedAttempts: 0,
				codes:          []string{"code1", "code3"},
				userLocked:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewHumanRecoveryCodeWriteModel("user1", "org1")
			wm.Events = tt.events

			err := wm.Reduce()
			require.NoError(t, err)

			assert.Equal(t, tt.want.State, wm.State)
			assert.Equal(t, tt.want.FailedAttempts, wm.FailedAttempts)
			assert.Equal(t, tt.want.codes, wm.codes)
			assert.Equal(t, tt.want.userLocked, wm.userLocked)
		})
	}
}
