package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_userStateForAuthentication(t *testing.T) {
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		userErrorID   string
		orgErrorID    string
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		args       args
		wantErr    func(error) bool
	}{
		{
			name: "user not active",
			eventstore: expectEventstore(
				expectFilter(
					user.NewHumanAddedEvent(
						context.Background(),
						&user.NewAggregate("userID", "org1").Aggregate,
						"username",
						"firstname",
						"lastname",
						"nickname",
						"displayname",
						language.English,
						domain.GenderUnspecified,
						"email",
						false,
					),
					user.NewUserDeactivatedEvent(
						context.Background(),
						&user.NewAggregate("userID", "org1").Aggregate,
					),
				),
			),
			args: args{
				ctx:           context.Background(),
				userID:        "userID",
				resourceOwner: "org1",
				userErrorID:   "TEST-1",
				orgErrorID:    "TEST-1-ORG",
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "org not active",
			eventstore: expectEventstore(
				expectFilter(
					user.NewHumanAddedEvent(
						context.Background(),
						&user.NewAggregate("userID", "org1").Aggregate,
						"username",
						"firstname",
						"lastname",
						"nickname",
						"displayname",
						language.English,
						domain.GenderUnspecified,
						"email",
						false,
					),
				),
				expectFilter(
					org.NewOrgAddedEvent(context.Background(),
						&org.NewAggregate("org1").Aggregate,
						"org"),
					org.NewOrgDeactivatedEvent(context.Background(),
						&org.NewAggregate("org1").Aggregate),
				),
			),
			args: args{
				ctx:           context.Background(),
				userID:        "userID",
				resourceOwner: "org1",
				userErrorID:   "TEST-2",
				orgErrorID:    "TEST-2-ORG",
			},
			wantErr: zerrors.IsPreconditionFailed,
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectFilter(
					user.NewHumanAddedEvent(
						context.Background(),
						&user.NewAggregate("userID", "org1").Aggregate,
						"username",
						"firstname",
						"lastname",
						"nickname",
						"displayname",
						language.English,
						domain.GenderUnspecified,
						"email",
						false,
					),
				),
				expectFilter(
					org.NewOrgAddedEvent(context.Background(),
						&org.NewAggregate("org1").Aggregate,
						"org"),
				),
			),
			args: args{
				ctx:           context.Background(),
				userID:        "userID",
				resourceOwner: "org1",
				userErrorID:   "TEST-3",
				orgErrorID:    "TEST-3-ORG",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.userStateForAuthentication(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.userErrorID, tt.args.orgErrorID)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, tt.wantErr(err))
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, "userID", got.AggregateID)
		})
	}
}
