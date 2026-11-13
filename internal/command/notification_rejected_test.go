package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_NotificationRejected(t *testing.T) {
	rejection := channels.Rejection{
		TriggeringEventType: user.HumanInviteCodeAddedType,
		Reason:              channels.RejectionReasonLimitExceeded,
		ProviderID:          "provider1",
	}
	userAdded := user.NewHumanAddedEvent(context.Background(),
		&user.NewAggregate("user1", "org1").Aggregate,
		"username", "firstName", "lastName", "nickName", "displayName",
		language.English, domain.GenderUnspecified, "email@example.com", false,
	)
	sessionAdded := session.NewAddedEvent(context.Background(),
		&session.NewAggregate("session1", "instance1").Aggregate,
		&domain.UserAgent{},
	)
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		aggregate  *eventstore.Aggregate
		wantErr    error
	}{
		{
			name:       "aggregate missing",
			eventstore: expectEventstore(),
			aggregate:  nil,
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Rj3k1", "Errors.IDMissing"),
		},
		{
			name:       "aggregate id missing",
			eventstore: expectEventstore(),
			aggregate:  &eventstore.Aggregate{Type: user.AggregateType, ResourceOwner: "org1"},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Rj3k1", "Errors.IDMissing"),
		},
		{
			name:       "aggregate type without notifications",
			eventstore: expectEventstore(),
			aggregate:  &eventstore.Aggregate{ID: "org1", Type: org.AggregateType, ResourceOwner: "org1"},
			wantErr:    zerrors.ThrowInvalidArgument(nil, "COMMAND-Rj3k2", "notifications of aggregate type org cannot be rejected"),
		},
		{
			name: "user not found",
			eventstore: expectEventstore(
				expectFilter(),
			),
			aggregate: &eventstore.Aggregate{ID: "user1", Type: user.AggregateType, ResourceOwner: "org1"},
			wantErr:   zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rj3k3", "Errors.User.NotFound"),
		},
		{
			name: "user removed",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(userAdded),
					eventFromEventPusher(
						user.NewUserRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username", nil, true,
						),
					),
				),
			),
			aggregate: &eventstore.Aggregate{ID: "user1", Type: user.AggregateType, ResourceOwner: "org1"},
			wantErr:   zerrors.ThrowPreconditionFailed(nil, "COMMAND-Rj3k3", "Errors.User.NotFound"),
		},
		{
			name: "user notification rejected",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(userAdded),
				),
				expectPush(
					user.NewNotificationRejectedEvent(context.Background(),
						&user.NewAggregate("user1", "org1").Aggregate,
						rejection,
					),
				),
			),
			aggregate: &eventstore.Aggregate{ID: "user1", Type: user.AggregateType, ResourceOwner: "org1"},
		},
		{
			name: "session not found",
			eventstore: expectEventstore(
				expectFilter(),
			),
			aggregate: &eventstore.Aggregate{ID: "session1", Type: session.AggregateType, ResourceOwner: "instance1"},
			wantErr:   zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
		},
		{
			name: "session terminated",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(sessionAdded),
					eventFromEventPusher(
						session.NewTerminateEvent(context.Background(), &session.NewAggregate("session1", "instance1").Aggregate),
					),
				),
			),
			aggregate: &eventstore.Aggregate{ID: "session1", Type: session.AggregateType, ResourceOwner: "instance1"},
			wantErr:   zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hewfq", "Errors.Session.Terminated"),
		},
		{
			name: "session notification rejected",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(sessionAdded),
				),
				expectPush(
					session.NewNotificationRejectedEvent(context.Background(),
						&session.NewAggregate("session1", "instance1").Aggregate,
						rejection,
					),
				),
			),
			aggregate: &eventstore.Aggregate{ID: "session1", Type: session.AggregateType, ResourceOwner: "instance1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{eventstore: tt.eventstore(t)}
			err := c.NotificationRejected(context.Background(), tt.aggregate, rejection)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
